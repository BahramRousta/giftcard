package redis

import (
	"fmt"
	"giftCard/config"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var redisPool *redis.Pool

func RedisInit(cnf *config.Config) {
	redisPool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", cnf.Redis.Host, cnf.Redis.Port))
			if err != nil {
				log.Fatalf("ERROR: fail init redis: %s", err.Error())
			}
			return conn, err
		},
	}
}

func GetRedisConn() redis.Conn {
	return redisPool.Get()
}
