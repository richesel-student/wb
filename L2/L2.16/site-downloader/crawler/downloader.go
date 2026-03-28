package crawler

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTP-клиент с таймаутом,
// чтобы запросы не зависали слишком долго
var client = &http.Client{
	Timeout: 10 * time.Second,
}

// Download выполняет HTTP-запрос с несколькими попытками
func Download(link string) ([]byte, string, error) {
	var lastErr error

	for i := 0; i < 3; i++ {
		resp, err := client.Get(link)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("status: %d", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = err
			continue
		}

		return body, resp.Header.Get("Content-Type"), nil
	}

	return nil, "", lastErr
}
