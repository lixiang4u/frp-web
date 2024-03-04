package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/frp-web/handler"
	"github.com/lixiang4u/frp-web/utils"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

//go:embed frp-web-h5/dist
var content embed.FS

var (
	cfgFilePath = fakeEmptyConfig()
	port        = utils.IWantUseHttpPort()
	frpWebRoot  = "frp-web" // 同 frp-web-h5/vite.config.js 中 base 值
)

func main() {
	go getVhostListOrCreate(port)

	httpServer(port)

}

func fakeEmptyConfig() string {
	var yamlFile = filepath.Join(os.TempDir(), "frp-web", fmt.Sprintf("tmp-frp-config-empty.toml"))
	_ = os.MkdirAll(filepath.Dir(yamlFile), fs.ModePerm)
	_ = os.Remove(yamlFile)
	var v = viper.New()
	if e := v.WriteConfigAs(yamlFile); e != nil {
		log.Println("[writeConfigError]", e.Error())
	}
	return yamlFile
}

func httpServer(port int) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.Static(fmt.Sprintf("/%s", frpWebRoot), filepath.Join("frp-web-h5", "dist"))
	r.GET("/api/config", handler.ApiServerConfig)
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/%s", frpWebRoot))
	})

	go func() { _ = r.Run(fmt.Sprintf(":%d", port)) }()
	go openBrowser()

	select {
	case _sig := <-sig:
		log.Println(fmt.Sprintf("[stop] %v\n", _sig))
	}
}

func runFrpService() error {
	cfg, proxyCfgs, visitorCfgs, isLegacyFormat, err := config.LoadClientConfig(cfgFilePath, false)
	if err != nil {
		return err
	}
	if isLegacyFormat {
		fmt.Printf("WARNING: ini format is deprecated and the support will be removed in the future, " +
			"please use yaml/json/toml format instead!\n")
	}

	var a = v1.ClientConfig{}
	a.Proxies = make([]v1.TypedProxyConfig, 1)
	a.Proxies[0].Type = string(v1.ProxyTypeHTTP)
	a.Proxies[0].ProxyConfigurer = v1.NewProxyConfigurerByType(v1.ProxyTypeHTTP)

	a.Proxies[0] = v1.TypedProxyConfig{Type: "A", ProxyConfigurer: v1.NewProxyConfigurerByType(v1.ProxyTypeHTTP)}

	proxyCfgs = append(proxyCfgs, a.Proxies[0])

	//log.InitLog(cfg.Log.To, cfg.Log.Level, cfg.Log.MaxDays, cfg.Log.DisablePrintColor)

	//if cfgFile != "" {
	//	log.Info("start frpc service for config file [%s]", cfgFile)
	//	defer log.Info("frpc service for config file [%s] stopped", cfgFile)
	//}
	cfgFilePath = ""
	svr, err := client.NewService(client.ServiceOptions{
		Common:         cfg,
		ProxyCfgs:      proxyCfgs,
		VisitorCfgs:    visitorCfgs,
		ConfigFilePath: cfgFilePath,
	})
	if err != nil {
		return err
	}

	//_ = os.Remove(cfgFilePath)

	shouldGracefulClose := cfg.Transport.Protocol == "kcp" || cfg.Transport.Protocol == "quic"
	// Capture the exit signal if we use kcp or quic.
	if shouldGracefulClose {
		go handleTermSignal(svr)
	}

	log.Println("[====================1>]", utils.ToJsonString(map[string]interface{}{
		"cfgFilePath":    cfgFilePath,
		"cfg":            cfg,
		"proxyCfgs":      proxyCfgs,
		"visitorCfgs":    visitorCfgs,
		"isLegacyFormat": isLegacyFormat,
	}))

	var cgf = v1.ClientConfig{}
	cgf.Complete()

	var ff = v1.ProxyBaseConfig{}
	ff.Complete("http")
	ff.LocalIP = ""
	var dd = v1.TypedClientPluginOptions{}
	dd.ClientPluginOptions = v1.HTTPS2HTTPPluginOptions{}
	ff.Plugin = dd

	log.Println("[====================2>]", utils.ToJsonString(map[string]interface{}{
		"cfgFilePath": cfgFilePath,
		"cfg":         cgf,
		"proxyCfgs":   ff,
	}))

	return svr.Run(context.Background())
}

func handleTermSignal(svr *client.Service) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
}

func openBrowser() {
	var osName = strings.ToLower(runtime.GOOS)
	switch osName {
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", fmt.Sprintf("http://127.0.0.1:%d/%s", port, frpWebRoot))
		_ = cmd.Run()
	case "darwin":
		cmd := exec.Command("open", fmt.Sprintf("http://127.0.0.1:%d/%s", port, frpWebRoot))
		_ = cmd.Run()
	}
}

func getVhostListOrCreate(localPort int) {
	if err := handler.NewClientVhost(localPort); err != nil {
		log.Println("[NewClientVhostError]", err.Error())
		os.Exit(1)
	}
	if err := handler.ClientVhostList(); err != nil {
		log.Println("[ClientVhostListError]", err.Error())
		os.Exit(1)
	}
}
