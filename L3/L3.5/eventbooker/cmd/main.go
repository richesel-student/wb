package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

	"eventbooker/internal/config"
	"eventbooker/internal/handler"
	"eventbooker/internal/repository"
	"eventbooker/internal/service"
	"eventbooker/internal/worker"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		log.Println("waiting db...")
		time.Sleep(time.Second)
	}

	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.Start(ctx, svc)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: h.Routes(),
	}

	go func() {
		log.Println("server started :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	cancel()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	srv.Shutdown(ctxShutdown)
}
