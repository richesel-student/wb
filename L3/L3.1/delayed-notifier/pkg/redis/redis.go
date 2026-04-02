package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Ctx = context.Background()

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	if err := Client.Ping(Ctx).Err(); err != nil {
		log.Fatal("❌ redis error:", err)
	}

	log.Println("✅ redis connected")
}
