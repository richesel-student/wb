package handler

import (
	"crypto-analytics/model"
	"crypto-analytics/service"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var t model.Transaction

	// читаем JSON
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}

	// валидация
	if err := t.Validate(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// вызываем сервис
	err = h.service.CreateTransaction(t)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// ответ
	w.Write([]byte("ok"))
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	items, _ := h.service.GetAll(from, to)
	json.NewEncoder(w).Encode(items)
}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/items/")
	id, _ := strconv.Atoi(idStr)

	var t model.Transaction
	json.NewDecoder(r.Body).Decode(&t)

	h.service.Update(id, t)

	w.Write([]byte("updated"))
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/items/")
	id, _ := strconv.Atoi(idStr)

	h.service.Delete(id)

	w.Write([]byte("deleted"))
}

func (h *Handler) Chart(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	data, err := h.service.GetChart(from, to)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(data)
}
func (h *Handler) Analytics(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	data, _ := h.service.GetAnalytics(from, to)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	println("EXPORT CALLED")
	items, err := h.service.GetAll("", "")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=data.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"ID", "Asset", "Type", "Amount"})

	for _, t := range items {
		writer.Write([]string{
			strconv.Itoa(t.ID),
			t.AssetName,
			t.Type,
			fmt.Sprintf("%.2f", t.AmountUSD),
		})
	}
}
