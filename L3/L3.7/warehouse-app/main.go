package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"warehouse-app/db"
	"warehouse-app/handlers"
	"warehouse-app/middleware"
	"warehouse-app/utils"
)

func main() {
	db.Init()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	r.Post("/login", login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Get("/api/items", handlers.GetItems)
		r.Get("/api/items/{id}/history", handlers.GetHistory)

		r.With(middleware.RoleAllowed("admin", "manager")).Post("/api/items", handlers.CreateItem)
		r.With(middleware.RoleAllowed("admin", "manager")).Put("/api/items/{id}", handlers.UpdateItem)
		r.With(middleware.RoleAllowed("admin")).Delete("/api/items/{id}", handlers.DeleteItem)
	})

	http.ListenAndServe(":8080", r)
}

func login(w http.ResponseWriter, r *http.Request) {
	var creds map[string]string
	json.NewDecoder(r.Body).Decode(&creds)

	row := db.DB.QueryRow("SELECT id, role FROM users WHERE username=$1 AND password=$2",
		creds["username"], creds["password"])

	var id int
	var role string

	err := row.Scan(&id, &role)
	if err != nil {
		http.Error(w, "invalid", 401)
		return
	}

	token, _ := utils.GenerateToken(id, role)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
