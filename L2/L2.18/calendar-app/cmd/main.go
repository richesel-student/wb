package main

import (
	"log"
	"net/http"
	"os"

	"calendar-app/internal/calendar"
	"calendar-app/internal/handlers"
	"calendar-app/internal/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	service := calendar.NewService()
	handler := handlers.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", handler.CreateEvent)
	mux.HandleFunc("/update_event", handler.UpdateEvent)
	mux.HandleFunc("/delete_event", handler.DeleteEvent)
	mux.HandleFunc("/events_for_day", handler.EventsForDay)
	mux.HandleFunc("/events_for_week", handler.EventsForWeek)
	mux.HandleFunc("/events_for_month", handler.EventsForMonth)

	log.Println("Server running on :" + port)

	log.Fatal(http.ListenAndServe(":"+port, middleware.Logger(mux)))
}