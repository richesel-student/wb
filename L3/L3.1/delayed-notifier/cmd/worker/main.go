package main

import (
	"encoding/json"
	"log"
	"math"
	"time"

	"delayed-notifier/internal/model"
	"delayed-notifier/internal/queue"
	"delayed-notifier/internal/repository"
	"delayed-notifier/pkg/db"
	"delayed-notifier/pkg/redis"
)

func main() {
	db.Init()
	redis.Init()
	queue.Init()

	ch := queue.Ch

	q, err := ch.QueueDeclare(
		"notifications",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,
		"notifications",
		"delayed",
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🚀 worker started")

	for msg := range msgs {
		log.Println("📩 received:", string(msg.Body))

		var n model.Notification
		if err := json.Unmarshal(msg.Body, &n); err != nil {
			log.Println("❌ json error:", err)

			if err := msg.Nack(false, false); err != nil {
				log.Println("nack error:", err)
			}
			continue
		}

		process(n)

		if err := msg.Ack(false); err != nil {
			log.Println("ack error:", err)
		}
	}
}

func process(n model.Notification) {
	dbN, err := repository.Get(n.ID)
	if err != nil {
		log.Println("❌ db error:", err)
		return
	}

	if dbN.Status == "canceled" || dbN.Status == "sent" {
		log.Println("⏭ skip:", n.ID)
		return
	}

	redis.Client.Set(redis.Ctx, "notify:"+n.ID, "processing", time.Hour)

	log.Println("📤 sending:", n.Message)
	err = nil // fake sender

	if err != nil {
		retries := dbN.Retries + 1
		log.Println("🔁 retry:", retries)

		if retries > 5 {
			if err := repository.UpdateStatus(n.ID, "failed"); err != nil {
				log.Println("update status error:", err)
			}
			return
		}

		delay := time.Duration(math.Pow(2, float64(retries))) * time.Second

		if err := repository.UpdateRetries(n.ID, retries); err != nil {
			log.Println("update retries error:", err)
		}

		if err := queue.Publish(n, delay); err != nil {
			log.Println("publish error:", err)
		}

		return
	}

	if err := repository.UpdateStatus(n.ID, "sent"); err != nil {
		log.Println("update status error:", err)
	}

	redis.Client.Set(redis.Ctx, "notify:"+n.ID, "sent", time.Hour)

	log.Println("✅ sent:", n.ID)
}
