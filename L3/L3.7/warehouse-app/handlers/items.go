package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"warehouse-app/db"
	"warehouse-app/middleware"
	"warehouse-app/utils"
)

type Item struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// АНТИПАТТЕРН: передаём user_id через SET LOCAL
func setUser(tx *sql.Tx, userID int) error {
	_, err := tx.Exec(fmt.Sprintf("SET LOCAL myapp.user_id = %d", userID))
	return err
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(middleware.UserContext).(*utils.Claims)

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if err := setUser(tx, claims.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := tx.Exec(
		"INSERT INTO items(name, quantity) VALUES($1,$2)",
		item.Name, item.Quantity,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id,name,quantity FROM items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ID, &i.Name, &i.Quantity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, i)
	}

	json.NewEncoder(w).Encode(items)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(middleware.UserContext).(*utils.Claims)

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if err := setUser(tx, claims.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := tx.Exec(
		"UPDATE items SET name=$1, quantity=$2, updated_at=now() WHERE id=$3",
		item.Name, item.Quantity, id,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := r.Context().Value(middleware.UserContext).(*utils.Claims)

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if err := setUser(tx, claims.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := tx.Exec("DELETE FROM items WHERE id=$1", id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	rows, err := db.DB.Query(`
		SELECT action, old_name, new_name, old_quantity, new_quantity, changed_by_user_id, changed_at
		FROM items_history WHERE item_id=$1 ORDER BY changed_at DESC`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var action, oldName, newName sql.NullString
		var oldQ, newQ sql.NullInt64
		var userID int
		var time string

		if err := rows.Scan(&action, &oldName, &newName, &oldQ, &newQ, &userID, &time); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result = append(result, map[string]interface{}{
			"action": action.String,
			"user":   userID,
			"time":   time,
		})
	}

	json.NewEncoder(w).Encode(result)
}
