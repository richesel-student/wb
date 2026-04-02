package repository

import (
	"context"
	"delayed-notifier/internal/model"
	"delayed-notifier/pkg/db"
)

func Create(n model.Notification) error {
	_, err := db.Conn.Exec(context.Background(),
		`INSERT INTO notifications (id, message, channel, recipient, send_at, status, retries, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		n.ID, n.Message, n.Channel, n.To, n.SendAt, n.Status, n.Retries, n.CreatedAt)

	return err
}

func Get(id string) (model.Notification, error) {
	var n model.Notification

	err := db.Conn.QueryRow(context.Background(),
		`SELECT id, message, channel, recipient, send_at, status, retries, created_at 
		 FROM notifications WHERE id=$1`, id).
		Scan(&n.ID, &n.Message, &n.Channel, &n.To, &n.SendAt, &n.Status, &n.Retries, &n.CreatedAt)

	return n, err
}

func UpdateStatus(id, status string) error {
	_, err := db.Conn.Exec(context.Background(),
		`UPDATE notifications SET status=$1 WHERE id=$2`, status, id)

	return err
}

func UpdateRetries(id string, retries int) error {
	_, err := db.Conn.Exec(context.Background(),
		`UPDATE notifications SET retries=$1 WHERE id=$2`, retries, id)

	return err
}
