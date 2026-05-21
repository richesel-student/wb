package handlers

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
)


func GCPercentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		currentPercent := debug.SetGCPercent(-1)
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(fmt.Sprintf("Current GC percent: %d\n", currentPercent))); err != nil {
			logHTTPError(r, err)
		}
		debug.SetGCPercent(currentPercent) // restore
	case http.MethodPost:
		value := r.URL.Query().Get("value")
		if value == "" {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("Missing 'value' parameter\n")); err != nil {
				logHTTPError(r, err)
			}
			return
		}

		percent, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("Invalid 'value' parameter: must be integer\n")); err != nil {
				logHTTPError(r, err)
			}
			return
		}

	
		if percent < 1 || percent > 500 {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("GC percent must be between 1 and 500\n")); err != nil {
				logHTTPError(r, err)
			}
			return
		}

		debug.SetGCPercent(percent)
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(fmt.Sprintf("GC percent set to: %d\n", percent))); err != nil {
			logHTTPError(r, err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("Only GET and POST methods allowed\n")); err != nil {
			logHTTPError(r, err)
		}
	}
}


func logHTTPError(r *http.Request, err error) {
	log.Printf("HTTP write error on %s %s: %v", r.Method, r.URL.Path, err)
}