package model

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/lixiang4u/frp-web/utils"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var (
	ApiServerHost string
	AppMachineId  string
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
	if err := viper.ReadInConfig(); err != nil {
		log.Println("[configError]", err.Error())
		utils.WaitInputExit()
	}
	if len(viper.GetString("app.server")) == 0 {
		log.Println("[configError] app.server 配置不能为空")
		utils.WaitInputExit()
	}
	ApiServerHost = strings.TrimRight(viper.GetString("app.server"), "/")
}
