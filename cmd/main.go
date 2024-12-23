package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/v-vovk/health-tracker-api/internal/config"
	"github.com/v-vovk/health-tracker-api/internal/db"
)

func main() {
	cfg := config.LoadConfig()

	database := db.Connect(cfg)
	log.Printf("Database connection established: %v", database)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health Tracker API is running!"))
	})

	port := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Starting server on %s", port)

	http.ListenAndServe(port, r)
}
