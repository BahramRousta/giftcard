package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Service  GiftCard `mapstructure:"service"`
	DataBase Postgres `mapstructure:"postgres"`
	Redis    Redis    `mapstructure:"redis"`
}

func LoadConfig() (*Config, error) {

	var config *Config

	dir, err := os.Getwd()
	viper.AddConfigPath(dir)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}

	err = viper.Unmarshal(&config)
	return config, nil
}
