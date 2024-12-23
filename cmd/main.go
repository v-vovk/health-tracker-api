package main

import (
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
	database.AutoMigrate(&food.Food{})

	r := chi.NewRouter()

	r.Use(middleware.JSONMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health Tracker API is running!"))
	})

	foodHandler := &food.Handler{DB: database}
	r.Route("/foods", func(r chi.Router) {
		r.Get("/", foodHandler.GetFoods)
		r.Post("/", foodHandler.CreateFood)
		r.Get("/{id}", foodHandler.GetFoodByID)
		r.Put("/{id}", foodHandler.UpdateFood)
		r.Delete("/{id}", foodHandler.DeleteFood)
	})

	port := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Starting server on %s", port)

	http.ListenAndServe(port, r)
}
