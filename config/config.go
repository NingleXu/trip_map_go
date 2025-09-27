package config

import (
	"github.com/spf13/viper"
)

var (
	SysConfig Config
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&SysConfig); err != nil {
		panic(err)
	}
}

type Config struct {
	Mysql  Mysql `yaml:"mysql"`
	Server ServerConfig
}
