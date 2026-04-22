package worker

import (
	"context"
	"log"
	"time"

	"eventbooker/internal/service"
)

func Start(ctx context.Context, s *service.Service) {
	backoff := time.Second
	maxBackoff := time.Minute

	for {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		err := s.Cleanup(cleanupCtx)
		cancel()

		if err != nil {
			log.Println("cleanup error:", err)
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		backoff = time.Second

		select {
		case <-time.After(10 * time.Second):
		case <-ctx.Done():
			return
		}
	}
}
