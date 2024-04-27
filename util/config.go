package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	BASE_URL      string `mapstructure:"BASE_URL"`
	CLIENT_ID     string `mapstructure:"CLIENT_ID"`
	CLIENT_SECRET string `mapstructure:"CLIENT_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
