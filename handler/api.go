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
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func ApiRecover(h gin.HandlerFunc) gin.HandlerFunc {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[recover]", err)
		}
	}()

	return func(ctx *gin.Context) {
		h(ctx)
	}
}

func ApiAuth(h gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ok, err := utils.IsIntranet(ctx.Request.Host)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":  http.StatusInternalServerError,
				"msg":   "请求来源异常",
				"debug": err.Error(),
			})
			return
		}
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"code":  http.StatusInternalServerError,
				"msg":   "请求来源异常，该请求只能内网访问",
				"debug": ctx.Request.Host,
			})
			return
		}
		h(ctx)
	}
}

func ApiServerConfig(ctx *gin.Context) {
	code, buf, _ := utils.HttpGet(fmt.Sprintf("%s/api/config", model.ApiServerHost))

	var resp gin.H
	_ = json.Unmarshal(buf, &resp)

	resp["machine_id"] = model.AppMachineId

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
		Id         string `json:"id" form:"id"`
		Type       string `json:"type" form:"type"`
		LocalAddr  string `json:"local_addr" form:"local_addr"`
		RemotePort int    `json:"remote_port" form:"remote_port"`
		Name       string `json:"name" form:"name"`     // 代码名称
		Status     bool   `json:"status" form:"status"` //true.运行，false.未运行
	}
	var req Req
	_ = ctx.ShouldBind(&req)

	var body = utils.ToJsonString(gin.H{
		"id":          req.Id,
		"type":        req.Type,
		"machine_id":  model.AppMachineId,
		"local_addr":  req.LocalAddr,
		"remote_port": req.RemotePort,
		"name":        req.Name,
		"status":      req.Status,
	})
	code, buf, _ := utils.HttpPost(fmt.Sprintf("%s/api/vhost", model.ApiServerHost), []byte(body))

	type Resp struct {
		Code  int         `json:"code"`
		Msg   string      `json:"msg"`
		Vhost model.Vhost `json:"vhost"`
	}
	var resp Resp
	_ = json.Unmarshal(buf, &resp)
	if resp.Code != http.StatusOK {
		ctx.JSON(http.StatusOK, gin.H{
			"code": resp.Code,
			"msg":  resp.Msg,
		})
		return
	}

	if resp.Vhost.Type == string(v1.ProxyTypeHTTPS) {
		if !strings.Contains(resp.Vhost.CrtPath, "CERTIFICATE") {
			_, _, _ = utils.HttpDelete(fmt.Sprintf("%s/api/vhost/%s/%s", model.ApiServerHost, model.AppMachineId, resp.Vhost.Id), nil)
			ctx.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "证书文件错误(cert)",
			})
			return
		}
		if !strings.Contains(resp.Vhost.KeyPath, "PRIVATE KEY") {
			_, _, _ = utils.HttpDelete(fmt.Sprintf("%s/api/vhost/%s/%s", model.ApiServerHost, model.AppMachineId, resp.Vhost.Id), nil)
			ctx.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "证书文件错误(key)",
			})
			return
		}
		removeCertFile(resp.Vhost.Id)
		_, _, err := parseCertToFile(resp.Vhost.Id, []byte(resp.Vhost.CrtPath), []byte(resp.Vhost.KeyPath))
		if err != nil {
			_, _, _ = utils.HttpDelete(fmt.Sprintf("%s/api/vhost/%s/%s", model.ApiServerHost, model.AppMachineId, resp.Vhost.Id), nil)
			ctx.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  fmt.Sprintf("证书保存失败：%s", err.Error()),
			})
			return
		}
	}

	ctx.JSON(code, resp)
}

func ApiServerRemoveVhost(ctx *gin.Context) {
	var vhostId = ctx.Param("vhost_id")

	code, buf, _ := utils.HttpDelete(fmt.Sprintf("%s/api/vhost/%s/%s", model.ApiServerHost, model.AppMachineId, vhostId), nil)

	var resp gin.H
	_ = json.Unmarshal(buf, &resp)

	ctx.JSON(code, resp)
}

func ApiUsePortCheck(ctx *gin.Context) {
	code, buf, _ := utils.HttpPost(fmt.Sprintf("%s/api/use-port-check", model.ApiServerHost), nil)

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
	tmpUrl, _ := url.PathUnescape(ctx.Request.RequestURI)

	root, _ := filepath.Abs(filepath.Join("."))
	tmpFile, _ := filepath.Abs(filepath.Join(".", tmpUrl))
	_, err := os.Stat(tmpFile)
	if err == nil && strings.HasPrefix(tmpFile, root) {
		ctx.Status(http.StatusOK)
		ctx.File(tmpFile)
		return
	}

	ctx.JSON(http.StatusNotFound, gin.H{
		"code": 404,
		"msg":  "请求地址错误",
		"path": tmpUrl,
	})
}

func handlerVhostConfig(vhosts []model.Vhost) []v1.ProxyConfigurer {
	var proxyCfgs = make([]v1.ProxyConfigurer, 0)
	for _, vhost := range vhosts {
		if !vhost.Status {
			continue
		}
		pc := v1.NewProxyConfigurerByType(v1.ProxyType(vhost.Type))
		tmpTypedConfig := handlerVhostConfigTyped(pc, vhost)
		if tmpTypedConfig != nil {
			proxyCfgs = append(proxyCfgs, tmpTypedConfig)
		}
	}
	return proxyCfgs
}

func handlerVhostConfigTyped(pc v1.ProxyConfigurer, vhost model.Vhost) (proxyCfg v1.ProxyConfigurer) {
	host, port, _ := net.SplitHostPort(vhost.LocalAddr)
	tmpVhostName := fmt.Sprintf("%s-%s-%s(%s)", vhost.Type, port, vhost.Id, vhost.Name)
	switch tmpC := pc.(type) {
	case *v1.HTTPProxyConfig:
		tmpC.Name = tmpVhostName
		tmpC.Transport.BandwidthLimitMode = "client"

		tmpC.LocalIP = host
		tmpC.LocalPort = utils.StringToInt(port)

		tmpC.CustomDomains = make([]string, 0)
		tmpC.CustomDomains = append(tmpC.CustomDomains, vhost.CustomDomain)
		if strings.Contains(vhost.CnameDomain, ".") {
			tmpC.CustomDomains = append(tmpC.CustomDomains, vhost.CnameDomain)
		}

		proxyCfg = tmpC
	case *v1.HTTPSProxyConfig:
		certFile, keyFile, _ := parseCertToFile(vhost.Id, []byte(vhost.CrtPath), []byte(vhost.KeyPath))

		// 参考frp实际运行的配置数据结构填充
		tmpC.Name = tmpVhostName
		tmpC.Transport.BandwidthLimitMode = "client"

		tmpC.LocalIP = host
		tmpC.LocalPort = utils.StringToInt(port)

		tmpC.CustomDomains = make([]string, 0)
		tmpC.CustomDomains = append(tmpC.CustomDomains, vhost.CustomDomain)

		tmpC.Plugin.Type = "https2http"
		tmpC.Plugin.ClientPluginOptions = &v1.HTTPS2HTTPPluginOptions{
			Type:              "https2http",
			LocalAddr:         vhost.LocalAddr,
			HostHeaderRewrite: tmpC.LocalIP,
			RequestHeaders: v1.HeaderOperations{
				Set: map[string]string{
					"x-from-where": "frp",
				},
			},
			CrtPath: certFile,
			KeyPath: keyFile,
		}

		proxyCfg = tmpC
	case *v1.TCPProxyConfig:
		tmpC.Name = tmpVhostName
		tmpC.Transport.BandwidthLimitMode = "client"

		tmpC.LocalIP = host
		tmpC.LocalPort = utils.StringToInt(port)

		tmpC.RemotePort = vhost.RemotePort

		proxyCfg = tmpC
	case *v1.TCPMuxProxyConfig:
		tmpC.Name = tmpVhostName
		tmpC.Transport.BandwidthLimitMode = "client"

		tmpC.LocalIP = host
		tmpC.LocalPort = utils.StringToInt(port)

		tmpC.CustomDomains = make([]string, 0)
		tmpC.CustomDomains = append(tmpC.CustomDomains, vhost.CustomDomain)

		tmpC.Multiplexer = "httpconnect"

		proxyCfg = tmpC

	default:

	}
	return proxyCfg
}

func removeCertFile(vhostId string) {
	_ = os.Remove(utils.AppTempFile("certs", fmt.Sprintf("%s-cert.pem", vhostId)))
	_ = os.Remove(utils.AppTempFile("certs", fmt.Sprintf("%s-key.pem", vhostId)))
}

func parseCertToFile(vhostId string, certBuf, keyBuf []byte) (certFile, keyFile string, err error) {
	certFile = utils.AppTempFile("certs", fmt.Sprintf("%s-cert.pem", vhostId))
	keyFile = utils.AppTempFile("certs", fmt.Sprintf("%s-key.pem", vhostId))
	if !utils.FileExists(certFile) {
		if err = os.WriteFile(certFile, certBuf, fs.ModePerm); err != nil {
			log.Println("[CertFileSaveError] ", vhostId, certFile)
			return
		}
	}
	if !utils.FileExists(keyFile) {
		if err = os.WriteFile(keyFile, keyBuf, fs.ModePerm); err != nil {
			log.Println("[KeyFileSaveError] ", vhostId, certFile)
			return
		}
	}
	return
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

var svr *client.Service

func runFrpClient(serverAddr string, serverPort int, vhosts []model.Vhost) (err error) {
	var cfg = &v1.ClientCommonConfig{}
	cfg.Complete()

	cfg.ServerAddr = serverAddr
	cfg.ServerPort = serverPort

	var proxyCfgs = handlerVhostConfig(vhosts)
	var visitorCfgs = make([]v1.VisitorConfigurer, 0)

	utils.FrpCloseRecover(svr)

	if len(proxyCfgs) == 0 {
		return errors.New("没有vhost可用")
	}

	svr, err = client.NewService(client.ServiceOptions{
		Common:         cfg,
		ProxyCfgs:      proxyCfgs,
		VisitorCfgs:    visitorCfgs,
		ConfigFilePath: "",
	})
	if err != nil {
		return err
	}

	shouldGracefulClose := cfg.Transport.Protocol == "kcp" || cfg.Transport.Protocol == "quic"
	// Capture the exit signal if we use kcp or quic.
	if shouldGracefulClose {
		go utils.FrpTermSignal(svr)
	}
	err = svr.Run(context.Background())
	if err != nil {
		log.Println("[frpRunError]", err.Error())
	}

	return err
}

func NewClientVhost(localPort int) error {
	var body = utils.ToJsonString(gin.H{
		"type":       string(v1.ProxyTypeHTTP),
		"machine_id": model.AppMachineId,
		"local_addr": fmt.Sprintf("127.0.0.1:%d", localPort),
		"name":       fmt.Sprintf("http-%d(default)", localPort),
		"status":     true,
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
