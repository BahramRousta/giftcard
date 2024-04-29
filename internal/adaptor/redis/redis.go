package redis

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"time"
)

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
