package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// =====================
// PROD DB
// =====================

func InitSQLite() *sql.DB {
	return initDB("./data.db")
}

// =====================
// TEST DB (in-memory)
// =====================

func InitTestSQLite() *sql.DB {
	return initDB(":memory:")
}

// =====================
// COMMON INIT
// =====================

func initDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	createTables(db)
	return db
}

// =====================
// TABLES
// =====================

func createTables(db *sql.DB) {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original TEXT,
		short TEXT UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS clicks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short TEXT,
		user_agent TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		log.Fatal(err)
	}
}
