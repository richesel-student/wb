package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shortener/internal/infrastructure/db"
	"shortener/internal/infrastructure/repository"
	"shortener/internal/usecase"
)

// helper
func setup() *Handler {
	database := db.InitTestSQLite()
	repo := repository.NewRepo(database)
	uc := usecase.NewShortenerUseCase(repo, repo, nil)

	return NewHandler(uc)
}

// =====================
// TEST: CREATE SHORT
// =====================

func TestShorten(t *testing.T) {
	h := setup()

	body := []byte(`{"url":"https://example.com"}`)
	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Shorten(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var res map[string]string
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	if res["short"] == "" {
		t.Fatal("short url is empty")
	}
}

// =====================
// TEST: REDIRECT
// =====================

func TestRedirect(t *testing.T) {
	h := setup()

	// create
	body := []byte(`{"url":"https://example.com"}`)
	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.Shorten(w, req)

	var res map[string]string
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	short := res["short"]

	// redirect
	req2 := httptest.NewRequest("GET", "/s/"+short, nil)
	w2 := httptest.NewRecorder()

	h.Redirect(w2, req2)

	if w2.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", w2.Code)
	}
}

// =====================
// TEST: ANALYTICS
// =====================

func TestAnalytics(t *testing.T) {
	h := setup()

	// create
	body := []byte(`{"url":"https://example.com"}`)
	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.Shorten(w, req)

	var res map[string]string
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	short := res["short"]

	// simulate click
	req2 := httptest.NewRequest("GET", "/s/"+short, nil)
	w2 := httptest.NewRecorder()
	h.Redirect(w2, req2)

	// analytics
	req3 := httptest.NewRequest("GET", "/analytics/"+short, nil)
	w3 := httptest.NewRecorder()
	h.Analytics(w3, req3)

	if w3.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w3.Code)
	}
}

func TestRedirect_NotFound_Handler(t *testing.T) {
	h := setup()

	req := httptest.NewRequest("GET", "/s/unknown", nil)
	w := httptest.NewRecorder()

	h.Redirect(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestShorten_EmptyBody(t *testing.T) {
	h := setup()

	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer([]byte(`{}`)))
	w := httptest.NewRecorder()

	h.Shorten(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAnalytics_Empty(t *testing.T) {
	h := setup()

	req := httptest.NewRequest("GET", "/analytics/unknown", nil)
	w := httptest.NewRecorder()

	h.Analytics(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200")
	}
}

func TestShorten_InvalidJSON(t *testing.T) {
	h := setup()

	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer([]byte(`invalid`)))
	w := httptest.NewRecorder()

	h.Shorten(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestShorten_Custom(t *testing.T) {
	h := setup()

	body := []byte(`{"url":"https://example.com","custom":"abc123"}`)
	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Shorten(w, req)

	var res map[string]string
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	if res["short"] != "abc123" {
		t.Fatal("custom short not used")
	}
}
