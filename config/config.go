package config

import (
	"github.com/spf13/viper"
	"os"
)

var (
	// Global config
	confs = Config{}
)

type Config struct {
	Service  GiftCard `mapstructure:"service"`
	DataBase Postgres `mapstructure:"postgres"`
	Redis    Redis    `mapstructure:"redis"`
	Jaeger   tracer   `yaml:"jaeger"`
	Logstash Logstash `yaml:"logstash"`
	//Debug    bool   `mapstructure:"debug"`
}

func init() {
	dir, _ := os.Getwd()

	viper.AddConfigPath(dir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	loadConfigs()
}

func loadConfigs() {
	if err := viper.Unmarshal(&confs); err != nil {
		panic(err)
	}
}

func C() *Config {
	return &confs
}
