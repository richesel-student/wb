package main

import (
	"bytes"
	"context"
	"encoding/json"
	_ "fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"warehouse-app/db"
	"warehouse-app/handlers"
	"warehouse-app/middleware"
)

func setupTestDB(t *testing.T) func() {
	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
	)
	if err != nil {
		t.Fatal(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	// env для config.go
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "testdb")

	// ждём пока postgres реально поднимется
	time.Sleep(2 * time.Second)

	db.Init()

	return func() {
		container.Terminate(ctx)
	}
}

func setupRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Get("/api/items", handlers.GetItems)
		r.Get("/api/items/{id}/history", handlers.GetHistory)

		r.With(middleware.RoleAllowed("admin", "manager")).Post("/api/items", handlers.CreateItem)
		r.With(middleware.RoleAllowed("admin", "manager")).Put("/api/items/{id}", handlers.UpdateItem)
		r.With(middleware.RoleAllowed("admin")).Delete("/api/items/{id}", handlers.DeleteItem)
	})

	return r
}

func TestFullFlow(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	server := httptest.NewServer(setupRouter())
	defer server.Close()

	client := &http.Client{}

	// LOGIN
	loginBody := []byte(`{"username":"admin","password":"admin"}`)
	resp, err := http.Post(server.URL+"/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		t.Fatal(err)
	}

	var loginResp map[string]string
	json.NewDecoder(resp.Body).Decode(&loginResp)

	token := loginResp["token"]
	if token == "" {
		t.Fatal("no token")
	}

	// CREATE
	req, _ := http.NewRequest("POST", server.URL+"/api/items",
		bytes.NewBuffer([]byte(`{"name":"apple","quantity":10}`)))

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create failed: %d", resp.StatusCode)
	}

	// GET
	req, _ = http.NewRequest("GET", server.URL+"/api/items", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	var items []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&items)

	if len(items) == 0 {
		t.Fatal("no items")
	}

	id := int(items[0]["id"].(float64))

	// UPDATE
	req, _ = http.NewRequest("PUT",
		server.URL+"/api/items/"+strconv.Itoa(id),
		bytes.NewBuffer([]byte(`{"name":"updated","quantity":20}`)))

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update failed: %d", resp.StatusCode)
	}

	// HISTORY
	req, _ = http.NewRequest("GET",
		server.URL+"/api/items/"+strconv.Itoa(id)+"/history",
		nil)

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	var history []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&history)

	if history[0]["action"] != "UPDATE" {
		t.Fatal("expected UPDATE")
	}

	// DELETE
	req, _ = http.NewRequest("DELETE",
		server.URL+"/api/items/"+strconv.Itoa(id),
		nil)

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete failed: %d", resp.StatusCode)
	}

	//  HISTORY again
	req, _ = http.NewRequest("GET",
		server.URL+"/api/items/"+strconv.Itoa(id)+"/history",
		nil)

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	json.NewDecoder(resp.Body).Decode(&history)

	if history[0]["action"] != "DELETE" {
		t.Fatal("expected DELETE")
	}
}
