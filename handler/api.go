package handler

import (
	"errors"
	"fmt"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/lixiang4u/frp-web/model"
	"github.com/lixiang4u/frp-web/utils"
	"net/http"
)

func ApiServerConfig(ctx *gin.Context) {
	code, buf, _ := utils.HttpGet(fmt.Sprintf("%s/api/config", model.ApiServerHost))

	var resp gin.H
	_ = json.Unmarshal(buf, &resp)

	ctx.JSON(code, resp)
}

func NewClientVhost() error {
	var body = utils.ToJsonString(gin.H{
		"type":       string(v1.ProxyTypeHTTP),
		"machine_id": model.AppMachineId,
		"local_addr": "type",
		"name":       "type",
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
	if v, ok := resp["code"]; !ok || v.(int) != http.StatusOK {
		msg, _ := resp["msg"]
		return errors.New(msg.(string))
	}

	return nil
}

func ClientVhostList() error {
	code, buf, err := utils.HttpGet(fmt.Sprintf("%s/api/vhosts", model.ApiServerHost))
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
	if v, ok := resp["code"]; !ok || v.(int) != http.StatusOK {
		msg, _ := resp["msg"]
		return errors.New(msg.(string))
	}

	return nil

}
