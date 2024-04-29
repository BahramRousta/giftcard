package config

import (
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	Username string `mapstructure:"redis.username" required:"true"`
	Password string `mapstructure:"redis.password" required:"true"`
	DB       int    `mapstructure:"redis.db" required:"true"`
	Host     string `mapstructure:"redis.host" required:"true"`
	Port     int    `mapstructure:"REDIS_PORT" required:"true"`
}

var redisPool *redis.Pool

func RedisInit() {
	redisPool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Printf("ERROR: fail init redis: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

func GetRedisConn() redis.Conn {
	return redisPool.Get()
}
