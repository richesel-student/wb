package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"calendar-app/internal/handlers"
	"calendar-app/internal/logger"
	"calendar-app/internal/repository"
	"calendar-app/internal/worker"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Инициализация компонентов
	asyncLogger := logger.NewAsyncLogger(1000) // Инициализируем логгер сначала
	repo := repository.NewRepository()
	handler := handlers.NewHandler(repo, asyncLogger) // Передаем логгер в хендлер

	reminderWorker := worker.NewReminderWorker(asyncLogger)
	archiveWorker := worker.NewArchiveWorker(repo, asyncLogger)

	// Устанавливаем reminder worker в handler
	handler.SetReminderWorker(reminderWorker)

	// Контекст для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Запуск асинхронного логгера
	asyncLogger.Start(ctx)

	// Запуск reminder worker
	reminderWorker.Start(ctx)

	// Запуск archive worker
	archiveWorker.Start(ctx)

	mux := http.NewServeMux()

	// 1. Маршрут для работы со списком событий (коллекцией)
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateEvent(w, r)
		case http.MethodGet:
			handler.GetEvents(w, r)
		default:
			writeMethodNotAllowed(w, asyncLogger)
		}
	})

	// 2. Маршрут для работы с конкретным событием по ID
	mux.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetEvent(w, r)
		case http.MethodPut:
			handler.UpdateEvent(w, r)
		case http.MethodDelete:
			handler.DeleteEvent(w, r)
		default:
			writeMethodNotAllowed(w, asyncLogger)
		}
	})

	mux.HandleFunc("/health", handler.HealthCheck)

	// Обертываем mux в middleware логирования
	handlerWithLogging := asyncLogger.HTTPLogger(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handlerWithLogging,
	}

	// Запуск сервера
	go func() {
		asyncLogger.Info("Server starting", map[string]interface{}{
			"port": port,
		})

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			asyncLogger.Error("Server error", map[string]interface{}{
				"error": err.Error(),
			})
			os.Exit(1)
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	asyncLogger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		asyncLogger.Error("Graceful shutdown failed", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		asyncLogger.Info("Server stopped gracefully")
	}
}

// writeMethodNotAllowed теперь принимает логгер и безопасно обрабатывает внутреннюю ошибку
func writeMethodNotAllowed(w http.ResponseWriter, log *logger.AsyncLogger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"}); err != nil {
		log.Error("Failed to write method not allowed response", map[string]interface{}{
			"error": err.Error(),
		})
	}
}