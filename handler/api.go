package handler

import (
	"errors"
	"fmt"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/lixiang4u/frp-web/model"
	"github.com/lixiang4u/frp-web/utils"
	"log"
	"net/http"
	"net/url"
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
		"type":       string(v1.ProxyTypeHTTP),
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
	ctx.JSON(http.StatusOK, gin.H{"ok": "ok"})
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
	log.Println("[resp]", utils.ToJsonString(resp))

	return nil

}
