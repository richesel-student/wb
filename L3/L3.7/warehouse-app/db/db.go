package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"warehouse-app/config"
)

var DB *sql.DB

func Init() {
	cfg := config.Load()

	var err error
	DB, err = sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	runMigrations()
}

func runMigrations() {
	query := `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username TEXT UNIQUE,
	password TEXT,
	role TEXT
);

INSERT INTO users (username, password, role)
VALUES
('admin','admin','admin'),
('manager','manager','manager'),
('viewer','viewer','viewer')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS items (
	id SERIAL PRIMARY KEY,
	name TEXT,
	quantity INT,
	updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS items_history (
	id SERIAL PRIMARY KEY,
	item_id INT,
	action TEXT,
	old_name TEXT,
	new_name TEXT,
	old_quantity INT,
	new_quantity INT,
	changed_by_user_id INT,
	changed_at TIMESTAMP DEFAULT now()
);

-- 🔴 АНТИПАТТЕРН: история через триггер
CREATE OR REPLACE FUNCTION log_item_changes()
RETURNS TRIGGER AS $$
DECLARE
	user_id INT;
BEGIN
	user_id := current_setting('myapp.user_id', true)::INT;

	IF TG_OP = 'INSERT' THEN
		INSERT INTO items_history(item_id, action, new_name, new_quantity, changed_by_user_id)
		VALUES (NEW.id, 'INSERT', NEW.name, NEW.quantity, user_id);
		RETURN NEW;

	ELSIF TG_OP = 'UPDATE' THEN
		INSERT INTO items_history(item_id, action, old_name, new_name, old_quantity, new_quantity, changed_by_user_id)
		VALUES (NEW.id, 'UPDATE', OLD.name, NEW.name, OLD.quantity, NEW.quantity, user_id);
		RETURN NEW;

	ELSIF TG_OP = 'DELETE' THEN
		INSERT INTO items_history(item_id, action, old_name, old_quantity, changed_by_user_id)
		VALUES (OLD.id, 'DELETE', OLD.name, OLD.quantity, user_id);
		RETURN OLD;
	END IF;

	RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS items_trigger ON items;

CREATE TRIGGER items_trigger
AFTER INSERT OR UPDATE OR DELETE ON items
FOR EACH ROW EXECUTE FUNCTION log_item_changes();
`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("migration error:", err)
	}

	fmt.Println("Migrations applied")
}
