package main

import (
	"log"
	"net/http"

	"crypto-analytics/bybit"
	"crypto-analytics/db"
	"crypto-analytics/handler"
	"crypto-analytics/repository"
	"crypto-analytics/service"
)

// 🔥 CORS middleware
func enableCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		h(w, r)
	}
}

func main() {
	// 📦 DB
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// // 🔄 auto-migrate
	// err = db.Migrate(database)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// 🔗 зависимости
	repo := repository.NewRepo(database)
	bybitClient := bybit.NewClient()
	service := service.NewService(repo, bybitClient)
	h := handler.NewHandler(service)

	// ======================
	// 📡 API ROUTES
	// ======================

	// POST + GET
	http.HandleFunc("/items", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h.Create(w, r)
			return
		} else if r.Method == "GET" {
			h.GetAll(w, r)
			return
		}

		http.Error(w, "Method not allowed", 405)
	}))

	// PUT + DELETE
	http.HandleFunc("/items/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			h.Update(w, r)
			return
		} else if r.Method == "DELETE" {
			h.Delete(w, r)
			return
		}

		http.Error(w, "Method not allowed", 405)
	}))

	// 📊 analytics
	http.HandleFunc("/analytics", enableCORS(h.Analytics))

	// ======================
	// 🌐 UI (статические файлы)
	// ======================

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
	http.HandleFunc("/export", enableCORS(h.ExportCSV))
	http.HandleFunc("/chart", enableCORS(h.Chart))

	// ======================

	log.Println("Server running on http://localhost:8080")
	log.Println("Open UI: http://localhost:8080/web/")

	http.ListenAndServe(":8080", nil)
}
