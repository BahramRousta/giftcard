package redis

import (
	"context"
	"encoding/json"
	"giftcard/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
	"time"
)

type Store struct {
	db *redis.Client
}

func NewRedis(lc fx.Lifecycle) *Store {
	rds := Store{}
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			log.Println("redis connected")
			return rds.connect(*config.C())
		},
		OnStop: func(c context.Context) error {
			log.Println("redis closed")
			return rds.db.Close()
		},
	})
	return &rds
}

func (r *Store) connect(confs config.Config) error {
	var err error

	r.db = redis.NewClient(&redis.Options{
		DB:   confs.Redis.DB,
		Addr: "localhost:6379",
	})

	if err = r.db.Ping(context.Background()).Err(); err != nil {
		zap.L().Error(err.Error())
	}

	return err
}

func (r *Store) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return r.db.Set(ctx, key, p, duration).Err()
}

// Get meth, get value with key
func (r *Store) Get(ctx context.Context, key string, dest interface{}) error {
	p, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(p), &dest)
}
