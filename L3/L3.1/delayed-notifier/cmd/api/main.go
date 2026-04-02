package main

import (
	"time"

	"delayed-notifier/internal/model"
	"delayed-notifier/internal/queue"
	"delayed-notifier/internal/repository"
	"delayed-notifier/pkg/db"
	"delayed-notifier/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	db.Init()
	redis.Init()
	queue.Init()

	r := gin.Default()

	// UI
	r.Static("/ui", "./web")

	// CREATE
	r.POST("/notify", func(c *gin.Context) {
		var n model.Notification

		if err := c.ShouldBindJSON(&n); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		n.ID = uuid.New().String()
		n.Status = "pending"
		n.CreatedAt = time.Now()

		// DB
		if err := repository.Create(n); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Redis (queued)
		redis.Client.Set(redis.Ctx, "notify:"+n.ID, "queued", time.Hour)

		// delay
		delay := time.Until(n.SendAt)
		if n.SendAt.IsZero() {
			delay = 5 * time.Second
		}

		// RabbitMQ
		if err := queue.Publish(n, delay); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"id": n.ID})
	})

	// GET
	r.GET("/notify/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Redis first
		val, err := redis.Client.Get(redis.Ctx, "notify:"+id).Result()
		if err == nil {
			c.JSON(200, gin.H{"status": val})
			return
		}

		// fallback DB
		n, err := repository.Get(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}

		c.JSON(200, gin.H{"status": n.Status})
	})

	// DELETE
	r.DELETE("/notify/:id", func(c *gin.Context) {
		id := c.Param("id")

		if err := repository.UpdateStatus(id, "canceled"); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		redis.Client.Set(redis.Ctx, "notify:"+id, "canceled", time.Hour)

		c.JSON(200, gin.H{"status": "canceled"})
	})

	// RUN
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
