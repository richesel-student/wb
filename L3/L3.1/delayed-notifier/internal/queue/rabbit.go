package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Ch *amqp.Channel

func Init() {
	var conn *amqp.Connection
	var err error

	// бесконечный retry
	for {
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			break
		}

		log.Println("⏳ waiting for rabbitmq...")
		time.Sleep(2 * time.Second)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	//  exchange с delay
	err = ch.ExchangeDeclare(
		"delayed",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	Ch = ch

	log.Println("✅ rabbitmq connected")
}

// отправка сообщения
func Publish(n interface{}, delay time.Duration) error {
	body, err := json.Marshal(n)
	if err != nil {
		return err
	}

	return Ch.PublishWithContext(
		context.Background(),
		"delayed",
		"notifications",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers: amqp.Table{
				"x-delay": int(delay.Milliseconds()),
			},
		},
	)
}
