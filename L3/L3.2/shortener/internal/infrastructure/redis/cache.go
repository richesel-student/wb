package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func InitRedis(addr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisCache) Set(key, value string) {
	r.client.Set(r.ctx, key, value, time.Hour)
}
