package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"image-processor/internal/queue"
	"image-processor/internal/repository"
	"image-processor/internal/storage"
	"image-processor/internal/worker"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Postgres
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

	// MinIO/S3
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("minio init error:", err)
	}
	st := storage.New(minioClient, os.Getenv("MINIO_BUCKET"))

	// Kafka Reader
	kafkaReader := queue.NewReader(os.Getenv("KAFKA_BROKER"), os.Getenv("KAFKA_TOPIC"), "worker-group")

	proc := worker.NewProcessor(repo, st)

	if err := proc.Start(ctx, kafkaReader); err != nil {
		log.Printf("Worker stopped with error: %v", err)
	}

	log.Println("Shutting down worker...")
	kafkaReader.Close()
	db.Close()
}
