package model

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/fatedier/frp/pkg/util/util"
	"github.com/lixiang4u/frp-web/utils"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var (
	ApiServerHost string
	AppMachineId  string
	AppInstance1  string
)

func init() {
	machineId, err := machineid.ID()
	if err != nil {
		log.Println("[MachineIdError]", err.Error())
		utils.WaitInputExit()
	}
	log.Println("[MachineId]", machineId)
	AppMachineId = machineId

	viper.SetConfigFile("config.toml")
	_ = viper.ReadInConfig()

	ApiServerHost = util.EmptyOr(strings.TrimRight(viper.GetString("app.server"), "/"), "http://api-frp.lixiang4u.xyz:7200")
	AppInstance1 = util.EmptyOr(strings.TrimSpace(viper.GetString("app.instance1")), "127.0.0.1:61234")

	if len(ApiServerHost) == 0 {
		log.Println("[configError] app.server 配置不能为空")
		utils.WaitInputExit()
	}
	if len(AppInstance1) == 0 {
		log.Println("[configError] app.instance1 配置不能为空")
		utils.WaitInputExit()
	}

	// 重新设置并写入配置
	viper.Set("app.server", ApiServerHost)
	viper.Set("app.instance1", AppInstance1)
	_ = viper.WriteConfig()
}
