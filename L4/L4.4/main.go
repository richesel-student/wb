package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gc-monitor/internal/handlers"
	"gc-monitor/internal/metrics"
)

func main() {
	// Initialize metrics collector
	metricsCollector := metrics.NewCollector()
	prometheus.MustRegister(metricsCollector)

	// Set up routes
	mux := http.NewServeMux()

	// Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())
	// GC percent control endpoint
	mux.HandleFunc("/gc-percent", handlers.GCPercentHandler)
	// Register pprof routes
	mux.Handle("/debug/pprof/", http.DefaultServeMux)

	
	// Create server with timeout configurations
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Graceful shutdown handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start HTTP server
	go func() {
		log.Printf("Starting server on port %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Printf("Shutdown signal received, shutting down gracefully...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop metrics collector ticker
	metricsCollector.Close()

	// Shutdown server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Printf("Server stopped gracefully")
}