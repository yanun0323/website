package config

import (
	"website/util"

	"github.com/spf13/viper"
)

const (
	SITE_ADDRESS = "https://www.yanunyang.com/"
)

func Init(configName string) error {
	name := configName
	if len(configName) > 0 {
		name = "config"
	}
	viper.AddConfigPath(util.GetAbsPath("config"))
	viper.AddConfigPath(util.GetAbsPath())
	viper.AutomaticEnv()
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}
