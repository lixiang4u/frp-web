package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/frp-web/handler"
	"github.com/lixiang4u/frp-web/model"
	"github.com/lixiang4u/frp-web/utils"
	"io/fs"
	"log"
	"mime"
	"net"
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

//go:embed frp-web-h5/dist/*
var frpWebContent embed.FS

var (
	frpWebRoot      = "frp-web" // 同 frp-web-h5/vite.config.js 中 base 值
	port            = utils.IWantUseHttpPort()
	localServerFile = utils.AppTempFile("local-web-server.json")
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	appOneInstanceCheck(port)
	go getVhostListOrCreate(port)
	go httpServer(port)

	select {
	case _sig := <-sig:
		_ = os.Remove(localServerFile)
		log.Println(fmt.Sprintf("[stop] %v\n", _sig))
	}

}

func httpServer(port int) {
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

	r.GET(fmt.Sprintf("/%s", frpWebRoot), frpWebFileServer)
	r.GET(fmt.Sprintf("/%s/*filepath", frpWebRoot), frpWebFileServer)

	r.GET("/api/config", handler.ApiRecover(handler.ApiAuth(handler.ApiServerConfig)))
	r.POST("/api/vhost", handler.ApiRecover(handler.ApiAuth(handler.ApiServerCreateVhost)))
	r.GET("/api/vhosts", handler.ApiRecover(handler.ApiAuth(handler.ApiServerVhostList)))
	r.DELETE("/api/vhost/:vhost_id", handler.ApiRecover(handler.ApiAuth(handler.ApiServerRemoveVhost)))
	r.POST("/api/frp/reload", handler.ApiRecover(handler.ApiAuth(handler.ApiFrpReload)))
	r.POST("/api/use-port-check", handler.ApiRecover(handler.ApiAuth(handler.ApiUsePortCheck)))

	r.NoRoute(handler.ApiRecover(handler.ApiNotRoute))

	go func() { _ = r.Run(fmt.Sprintf(":%d", port)) }()
	go openBrowser(port)
}

func openBrowser(port int) {
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
		utils.WaitInputExit()
	}
	if err := handler.ClientVhostList(); err != nil {
		log.Println("[ClientVhostListError]", err.Error())
		utils.WaitInputExit()
	}
}

func appOneInstanceCheck(port int) {
	var isRun = make(chan bool, 1)

	go func() {
		if len(model.AppInstance1) == 0 {
			return
		}
		l, err := net.Listen("tcp", model.AppInstance1)
		if err != nil {
			log.Println("[程序使用(tcp://127.0.0.98:61234)检测多开问题]", err.Error())

			// 如果启动了则打开网页
			var ls = model.LocalServer{}
			buf, _ := os.ReadFile(localServerFile)
			if err = json.Unmarshal(buf, &ls); err == nil && ls.Port > 0 {
				openBrowser(ls.Port)
			}

			isRun <- false
		} else {
			go func() { _, _ = l.Accept() }()

			// 如果没启动，则写入启动文件
			var ls = model.LocalServer{
				Url:  fmt.Sprintf("http://127.0.0.1:%d/%s", port, frpWebRoot),
				Port: port,
			}
			_ = os.WriteFile(localServerFile, []byte(utils.ToJsonString(ls)), fs.ModePerm)

			isRun <- true
		}
	}()
	if !<-isRun {
		utils.WaitInputExit()
		os.Exit(1)
	}

}

func frpWebFileServer(ctx *gin.Context) {
	var fp = ctx.Param("filepath")
	if len(fp) == 0 || fp == "/" {
		fp = "index.html"
	}
	var cf = strings.ReplaceAll(filepath.Join("frp-web-h5", "dist", fp), "\\", "/")
	buf, err := fs.ReadFile(frpWebContent, cf)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "请求地址错误(embed.FS)",
			"path": ctx.Param("filepath"),
		})
		return
	}
	ctx.Data(http.StatusOK, mime.TypeByExtension(filepath.Ext(cf)), buf)
}
