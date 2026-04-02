package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func Init() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@postgres:5432/postgres")
	if err != nil {
		panic(err)
	}
	Conn = conn
}
