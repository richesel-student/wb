package main

import (
	"log"
	"net/http"

	"comments-service/internal/handler"
	"comments-service/internal/repository"
	"comments-service/internal/service"
)

func main() {
	repo := repository.NewMemoryRepository()
	svc := service.NewCommentService(repo)
	h := handler.NewCommentHandler(svc)

	http.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			h.Get(w, r)
		}
	})

	http.HandleFunc("/comments/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			h.Delete(w, r)
		}
	})

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Println("🚀 Server started at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
