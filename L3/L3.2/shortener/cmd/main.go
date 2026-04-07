//

package main

import (
	"log"
	"net/http"
	"os"

	delivery "shortener/internal/delivery/http"
	"shortener/internal/infrastructure/db"
	"shortener/internal/infrastructure/redis"
	"shortener/internal/infrastructure/repository"
	"shortener/internal/usecase"
)

func main() {
	// DB
	database := db.InitSQLite()

	// Redis
	cache := redis.InitRedis(os.Getenv("REDIS_ADDR"))

	// Repo
	repo := repository.NewRepo(database)

	// Usecase
	uc := usecase.NewShortenerUseCase(repo, repo, cache)

	// Handler
	h := delivery.NewHandler(uc)

	http.HandleFunc("/shorten",
		delivery.Chain(h.Shorten, delivery.Logging, delivery.RateLimit),
	)

	http.HandleFunc("/s/",
		delivery.Chain(h.Redirect, delivery.Logging),
	)

	http.HandleFunc("/analytics/",
		delivery.Chain(h.Analytics, delivery.Logging, delivery.RateLimit),
	)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
