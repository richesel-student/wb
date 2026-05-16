package logger

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type LogEntry struct {
	Time    time.Time
	Level   string
	Message string
	Fields  map[string]interface{}
}

type AsyncLogger struct {
	entries chan LogEntry
}

func NewAsyncLogger(bufferSize int) *AsyncLogger {
	return &AsyncLogger{
		entries: make(chan LogEntry, bufferSize),
	}
}

func (l *AsyncLogger) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Обрабатываем оставшиеся логи перед выходом
				l.drain()
				log.Println("Async logger stopped")
				return
			case entry := <-l.entries:
				l.writeEntry(entry)
			}
		}
	}()
}

func (l *AsyncLogger) drain() {
	close(l.entries)
	for entry := range l.entries {
		l.writeEntry(entry)
	}
}

func (l *AsyncLogger) writeEntry(entry LogEntry) {
	fieldsStr := ""
	if len(entry.Fields) > 0 {
		for k, v := range entry.Fields {
			fieldsStr += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	log.Printf("[%s] %s%s\n", entry.Level, entry.Message, fieldsStr)
}

func (l *AsyncLogger) Info(message string, fields ...map[string]interface{}) {
	l.log("INFO", message, fields...)
}

func (l *AsyncLogger) Error(message string, fields ...map[string]interface{}) {
	l.log("ERROR", message, fields...)
}

func (l *AsyncLogger) Warn(message string, fields ...map[string]interface{}) {
	l.log("WARN", message, fields...)
}

func (l *AsyncLogger) Debug(message string, fields ...map[string]interface{}) {
	l.log("DEBUG", message, fields...)
}

func (l *AsyncLogger) log(level, message string, fields ...map[string]interface{}) {
	entry := LogEntry{
		Time:    time.Now(),
		Level:   level,
		Message: message,
	}

	if len(fields) > 0 {
		entry.Fields = fields[0]
	}

	select {
	case l.entries <- entry:
		// успешно добавлено в очередь
	default:
		// если очередь переполнена, пишем синхронно
		log.Printf("WARN: logger queue full, writing sync: [%s] %s\n", level, message)
		l.writeEntry(entry)
	}
}

// HTTPLogger middleware для логирования HTTP запросов
func (l *AsyncLogger) HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Используем response writer для захвата статуса
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		l.Info("HTTP request", map[string]interface{}{
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   rw.statusCode,
			"duration": duration.String(),
			"ip":       r.RemoteAddr,
		})
	})
}

// responseWriter для захвата статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}