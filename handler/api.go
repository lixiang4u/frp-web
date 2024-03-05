package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatedier/frp/client"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/lixiang4u/frp-web/model"
	"github.com/lixiang4u/frp-web/utils"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func ApiServerConfig(ctx *gin.Context) {
	code, buf, _ := utils.HttpGet(fmt.Sprintf("%s/api/config", model.ApiServerHost))

	var resp gin.H
	_ = json.Unmarshal(buf, &resp)

	ctx.JSON(code, resp)
}

func ApiServerVhostList(ctx *gin.Context) {
	var params = url.Values{}
	params.Add("machine_id", model.AppMachineId)
	code, buf, _ := utils.HttpGet(fmt.Sprintf("%s/api/vhosts", model.ApiServerHost), params)

	var resp gin.H
	_ = json.Unmarshal(buf, &resp)

	ctx.JSON(code, resp)
}

func ApiServerCreateVhost(ctx *gin.Context) {
	type Req struct {
		Type      string `json:"type" form:"type"`
		LocalAddr string `json:"local_addr" form:"local_addr"`
		Name      string `json:"name" form:"name"` // 代码名称
	}
	var req Req
	_ = ctx.ShouldBind(&req)

	var body = utils.ToJsonString(gin.H{
		"type":       req.Type,
		"machine_id": model.AppMachineId,
		"local_addr": req.LocalAddr,
		"name":       req.Name,
	})
	code, buf, _ := utils.HttpPost(fmt.Sprintf("%s/api/vhost", model.ApiServerHost), []byte(body))

	var resp gin.H
	_ = json.Unmarshal(buf, &resp)

	ctx.JSON(code, resp)
}

func ApiServerRemoveVhost(ctx *gin.Context) {
	var vhostId = ctx.Param("vhost_id")

	code, buf, _ := utils.HttpDelete(fmt.Sprintf("%s/api/vhost/%s/%s", model.ApiServerHost, model.AppMachineId, vhostId), nil)

	var resp gin.H
	_ = json.Unmarshal(buf, &resp)

	ctx.JSON(code, resp)
}

func ApiFrpReload(ctx *gin.Context) {
	_, buf, _ := utils.HttpGet(fmt.Sprintf("%s/api/config", model.ApiServerHost))

	type Resp struct {
		Code   int `json:"code"`
		Config struct {
			BindPort       int    `json:"bind_port"`
			VhostHttpPort  int    `json:"vhost_http_port"`
			VhostHttpsPort int    `json:"vhost_https_port"`
			Host           string `json:"host"`
		} `json:"config"`
	}
	var resp Resp
	_ = json.Unmarshal(buf, &resp)
	if resp.Code != http.StatusOK {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "frp服务器信息获取失败",
		})
		return
	}

	vhosts, err := apiGetVhosts()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	go func() {
		log.Println(fmt.Sprintf("[frpServer] %s:%d", resp.Config.Host, resp.Config.BindPort))

		if err := runFrpClient(resp.Config.Host, resp.Config.BindPort, vhosts); err != nil {
			log.Println("[frpClientError]", err.Error())
		}
	}()

	ctx.JSON(http.StatusOK, resp)
}

func ApiNotRoute(ctx *gin.Context) {
	root, _ := filepath.Abs(filepath.Join("."))
	tmpFile, _ := filepath.Abs(filepath.Join(".", ctx.Request.RequestURI))
	_, err := os.Stat(tmpFile)
	if err == nil && strings.HasPrefix(tmpFile, root) {
		ctx.File(tmpFile)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 404,
		"msg":  "请求地址错误",
		"uri":  ctx.Request.RequestURI,
	})
}

func handlerVhostConfig(vhosts []model.Vhost) []v1.ProxyConfigurer {
	var proxyCfgs = make([]v1.ProxyConfigurer, 0)
	for _, vhost := range vhosts {
		pc := v1.NewProxyConfigurerByType(v1.ProxyType(vhost.Type))
		proxyCfgs = append(proxyCfgs, handlerVhostConfigTyped(pc, vhost))
	}
	return proxyCfgs
}

func handlerVhostConfigTyped(pc v1.ProxyConfigurer, vhost model.Vhost) (proxyCfg v1.ProxyConfigurer) {
	host, port, _ := net.SplitHostPort(vhost.LocalAddr)
	switch tmpC := pc.(type) {
	case *v1.HTTPProxyConfig:
		//tmpC.Type = ""
		tmpC.Name = vhost.Name
		tmpC.Transport.BandwidthLimitMode = "client"

		tmpC.LocalIP = host
		tmpC.LocalPort = utils.StringToInt(port)

		tmpC.CustomDomains = make([]string, 0)
		tmpC.CustomDomains = append(tmpC.CustomDomains, vhost.CustomDomain)

		proxyCfg = tmpC
	case *v1.HTTPSProxyConfig:
		proxyCfg = tmpC
	default:

	}
	return proxyCfg
}

func apiGetVhosts() ([]model.Vhost, error) {
	var params = url.Values{}
	params.Add("machine_id", model.AppMachineId)

	code, buf, _ := utils.HttpGet(fmt.Sprintf("%s/api/vhosts", model.ApiServerHost), params)

	type Resp struct {
		Code   int           `json:"code"`
		Msg    string        `json:"msg"`
		Vhosts []model.Vhost `json:"vhosts"`
	}
	var resp Resp
	_ = json.Unmarshal(buf, &resp)
	if code != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("请求失败：status_code=%d", code))
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Msg)
	}
	return resp.Vhosts, nil
}

//
//func apiGetConfig() {
//	code, buf, _ := utils.HttpGet(fmt.Sprintf("%s/api/config", model.ApiServerHost))
//
//	type Resp struct {
//		Code   int `json:"code"`
//		Config struct {
//			BindPort       int    `json:"bind_port"`
//			VhostHttpPort  int    `json:"vhost_http_port"`
//			VhostHttpsPort int    `json:"vhost_https_port"`
//			Host           string `json:"host"`
//		} `json:"config"`
//	}
//
//	var resp Resp
//	_ = json.Unmarshal(buf, &resp)
//
//	ctx.JSON(code, resp)
//
//}

func runFrpClient(serverAddr string, serverPort int, vhosts []model.Vhost) error {
	var cfg = &v1.ClientCommonConfig{}
	cfg.Complete()

	cfg.ServerAddr = serverAddr
	cfg.ServerPort = serverPort

	var proxyCfgs = handlerVhostConfig(vhosts)
	var visitorCfgs = make([]v1.VisitorConfigurer, 0)

	svr, err := client.NewService(client.ServiceOptions{
		Common:         cfg,
		ProxyCfgs:      proxyCfgs,
		VisitorCfgs:    visitorCfgs,
		ConfigFilePath: "",
	})
	if err != nil {
		return err
	}

	if model.AppFrpRun == false {
		model.AppFrpRun = true

		shouldGracefulClose := cfg.Transport.Protocol == "kcp" || cfg.Transport.Protocol == "quic"
		// Capture the exit signal if we use kcp or quic.
		if shouldGracefulClose {
			go utils.FrpTermSignal(svr)
		}

		err = svr.Run(context.Background())
		if err != nil {
			log.Println("[frpRunError]", err.Error())
		}
	} else {
		err = svr.UpdateAllConfigurer(proxyCfgs, visitorCfgs)
		if err != nil {
			log.Println("[frpUpdateConfigError]", err.Error())
		}
	}

	return err
}

func NewClientVhost(localPort int) error {
	var body = utils.ToJsonString(gin.H{
		"type":       string(v1.ProxyTypeHTTP),
		"machine_id": model.AppMachineId,
		"local_addr": fmt.Sprintf("127.0.0.1:%d", localPort),
		"name":       fmt.Sprintf("frp-%s-%d", model.AppMachineId[:6], localPort),
	})
	code, buf, err := utils.HttpPost(fmt.Sprintf("%s/api/vhost", model.ApiServerHost), []byte(body))
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return errors.New(fmt.Sprintf("statusCode: %d", code))
	}

	var resp gin.H
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return err
	}
	if v, ok := resp["code"]; !ok || int(v.(float64)) != http.StatusOK {
		msg, _ := resp["msg"]
		return errors.New(msg.(string))
	}

	return nil
}

func ClientVhostList() error {
	var params = url.Values{}
	params.Add("machine_id", model.AppMachineId)

	code, buf, err := utils.HttpGet(fmt.Sprintf("%s/api/vhosts", model.ApiServerHost), params)
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return errors.New(fmt.Sprintf("statusCode: %d", code))
	}

	var resp gin.H
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return err
	}
	if v, ok := resp["code"]; !ok || int(v.(float64)) != http.StatusOK {
		msg, _ := resp["msg"]
		return errors.New(msg.(string))
	}

	return nil

}
