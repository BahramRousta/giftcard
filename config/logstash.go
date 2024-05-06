package config

type Logstash struct {
	Endpoint string `mapstructure:"endpoint" validate:"required"`
	Timeout  int    `mapstructure:"timeout"`
}
