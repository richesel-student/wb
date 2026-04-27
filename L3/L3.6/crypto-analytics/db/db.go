package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	//connStr := "user=postgres password=1234 dbname=crypto sslmode=disable"
	connStr := "host=db user=postgres password=1234 dbname=crypto sslmode=disable"
	return sql.Open("postgres", connStr)
}
