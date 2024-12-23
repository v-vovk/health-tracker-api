package main

import (
	"encoding/json"
	"fmt"
	"github.com/v-vovk/health-tracker-api/internal/food"
	"github.com/v-vovk/health-tracker-api/internal/middleware"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/v-vovk/health-tracker-api/internal/config"
	"github.com/v-vovk/health-tracker-api/internal/db"
)

func main() {
	cfg := config.LoadConfig()

	database := db.Connect(cfg)
	if database == nil {
		log.Fatal("Failed to connect to the database")
	}
	database.AutoMigrate(&food.Food{})

	r := chi.NewRouter()

	r.Use(middleware.JSONMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "Health Tracker API is running!"})
	})

	foodHandler := &food.Handler{DB: database}
	r.Mount("/foods", foodHandler.Routes())

	port := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Starting server on %s", port)

	http.ListenAndServe(port, r)
}
