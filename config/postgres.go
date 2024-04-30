package config

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
	Schema   string `mapstructure:"schema"`
	SSLMode  string `mapstructure:"ssl_mode"`
	TimeZone string `mapstructure:"timezone"`
}
