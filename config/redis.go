package config

type Redis struct {
	Username string `mapstructure:"redis.username" required:"true"`
	Password string `mapstructure:"redis.password" required:"true"`
	DB       int    `mapstructure:"redis.db" required:"true"`
	Host     string `mapstructure:"redis.host" required:"true"`
	Port     int    `mapstructure:"REDIS_PORT" required:"true"`
}
