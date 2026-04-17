package queue

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Queue struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewWriter(broker, topic string) *Queue {
	return &Queue{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func NewReader(broker, topic, group string) *Queue {
	return &Queue{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: group,
		}),
	}
}

func (q *Queue) Send(ctx context.Context, id string) error {
	return q.writer.WriteMessages(ctx, kafka.Message{Value: []byte(id)})
}

func (q *Queue) Read(ctx context.Context) (string, error) {
	m, err := q.reader.ReadMessage(ctx)
	if err != nil {
		return "", err
	}
	return string(m.Value), nil
}

func (q *Queue) Close() {
	if q.reader != nil {
		q.reader.Close()
	}
	if q.writer != nil {
		q.writer.Close()
	}
}
