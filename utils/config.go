package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DbDriver      string `mapstructure:"DB_DRIVER"`
	DbSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfigFile(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	fmt.Println(err)
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
