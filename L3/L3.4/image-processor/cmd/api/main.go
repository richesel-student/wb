package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"image-processor/internal/handler"
	"image-processor/internal/queue"
	"image-processor/internal/repository"
	"image-processor/internal/service"
	"image-processor/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// --- Postgres ---
	config, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal("invalid POSTGRES_DSN:", err)
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("db connect error:", err)
	}

	if err := db.Ping(ctx); err != nil {
		log.Fatal("db ping error:", err)
	}

	repo := repository.New(db)

	// --- MinIO ---
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("minio init error:", err)
	}

	st := storage.New(minioClient, os.Getenv("MINIO_BUCKET"))

	// --- Kafka ---
	kafkaWriter := queue.NewWriter(os.Getenv("KAFKA_BROKER"), os.Getenv("KAFKA_TOPIC"))

	// --- Service & Handler ---
	svc := service.New(repo, st, kafkaWriter)
	h := handler.New(svc, repo, st)

	// --- Router ---
	mux := http.NewServeMux()

	// upload
	mux.HandleFunc("/upload", h.Upload)

	// image/{id}
	mux.HandleFunc("/image/", func(w http.ResponseWriter, r *http.Request) {
		// защита от /image без id
		if strings.TrimPrefix(r.URL.Path, "/image/") == "" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.GetImage(w, r)
		case http.MethodDelete:
			h.DeleteImage(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// list
	mux.HandleFunc("/images", h.List)

	// --- Static (UI) ---
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	// --- Server ---
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("API listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down API server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Println("Server Shutdown error:", err)
	}

	kafkaWriter.Close()
	db.Close()
}
