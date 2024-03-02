package model

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/spf13/viper"
	"log"
	"os"
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
		os.Exit(1)
	}
	log.Println("[MachineId]", machineId)
	AppMachineId = machineId

	viper.SetConfigFile("config.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("[configError]", err.Error())
		os.Exit(1)
	}
	ApiServerHost = strings.TrimRight(viper.GetString("app.server"), "/")
}
