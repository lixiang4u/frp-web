package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/frp-web/handler"
	"github.com/lixiang4u/frp-web/utils"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var (
	frpWebRoot  = "frp-web" // 同 frp-web-h5/vite.config.js 中 base 值
	port        = utils.IWantUseHttpPort()
	appLockFile = "run.lock"
	appRunFile  *os.File
)

func main() {
	defer func() {
		if appRunFile != nil {
			_ = appRunFile.Close()
		}
		_ = os.Remove(appLockFile)
	}()

	appOneInstanceCheck()

	go getVhostListOrCreate(port)

	httpServer(port)

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

	r.Static(fmt.Sprintf("/%s/", frpWebRoot), filepath.Join("frp-web-h5", "dist"))
	r.GET("/api/config", handler.ApiServerConfig)
	r.POST("/api/vhost", handler.ApiServerCreateVhost)
	r.GET("/api/vhosts", handler.ApiServerVhostList)
	r.DELETE("/api/vhost/:vhost_id", handler.ApiServerRemoveVhost)
	r.POST("/api/frp/reload", handler.ApiFrpReload)

	r.NoRoute(handler.ApiNotRoute)

	go func() { _ = r.Run(fmt.Sprintf(":%d", port)) }()
	go openBrowser()

	select {
	case _sig := <-sig:
		log.Println(fmt.Sprintf("[stop] %v\n", _sig))
	}
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

func appOneInstanceCheck() {
	var err error
	appRunFile, err = os.OpenFile(appLockFile, os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		log.Println("程序已运行：" + err.Error())
		os.Exit(1)
	}
}
