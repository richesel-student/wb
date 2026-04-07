package usecase

import (
	"fmt"
	"testing"

	"shortener/internal/domain"
)

// ===== MOCK REPO =====

type mockRepo struct {
	store map[string]string
}

func (m *mockRepo) Save(u domain.URL) error {
	m.store[u.Short] = u.Original
	return nil
}

func (m *mockRepo) Get(short string) (string, error) {
	val, ok := m.store[short]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return val, nil
}

func (m *mockRepo) SaveClick(c domain.Click) error {
	return nil
}

func (m *mockRepo) GetStats(short string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"total": 1,
	}, nil
}

// ===== MOCK CACHE =====

type mockCache struct {
	data map[string]string
}

func (m *mockCache) Get(key string) (string, error) {
	val, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return val, nil
}

func (m *mockCache) Set(key, value string) {
	m.data[key] = value
}

// ===== TESTS =====

func TestCreate(t *testing.T) {
	repo := &mockRepo{store: make(map[string]string)}
	uc := NewShortenerUseCase(repo, repo, nil)

	short, err := uc.Create("https://example.com", "")

	if err != nil {
		t.Fatal(err)
	}
	if short == "" {
		t.Fatal("short is empty")
	}
}

func TestCustomShort(t *testing.T) {
	repo := &mockRepo{store: make(map[string]string)}
	uc := NewShortenerUseCase(repo, repo, nil)

	short, _ := uc.Create("https://example.com", "custom")

	if short != "custom" {
		t.Fatal("custom short not used")
	}
}

func TestRedirect_WithCache(t *testing.T) {
	repo := &mockRepo{store: map[string]string{
		"abc": "https://example.com",
	}}

	cache := &mockCache{
		data: map[string]string{
			"abc": "https://cached.com",
		},
	}

	uc := NewShortenerUseCase(repo, repo, cache)

	url, _ := uc.Redirect("abc", "ua")

	if url != "https://cached.com" {
		t.Fatal("should use cache")
	}
}

func TestRedirect_DB(t *testing.T) {
	repo := &mockRepo{store: map[string]string{
		"abc": "https://example.com",
	}}

	uc := NewShortenerUseCase(repo, repo, nil)

	url, _ := uc.Redirect("abc", "ua")

	if url != "https://example.com" {
		t.Fatal("wrong redirect")
	}
}

func TestRedirect_NotFound(t *testing.T) {
	repo := &mockRepo{store: map[string]string{}}
	uc := NewShortenerUseCase(repo, repo, nil)

	_, err := uc.Redirect("nope", "ua")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAnalytics(t *testing.T) {
	repo := &mockRepo{store: make(map[string]string)}
	uc := NewShortenerUseCase(repo, repo, nil)

	data, _ := uc.Analytics("test")

	if data["total"] != 1 {
		t.Fatal("analytics broken")
	}
}
