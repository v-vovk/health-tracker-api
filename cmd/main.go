package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/v-vovk/health-tracker-api/internal/app/food"
	"github.com/v-vovk/health-tracker-api/internal/infra/config"
	"github.com/v-vovk/health-tracker-api/internal/infra/db"
	"github.com/v-vovk/health-tracker-api/internal/infra/logger"
	"github.com/v-vovk/health-tracker-api/internal/infra/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	logger.InitLogger(cfg.Env)
	defer logger.Sync()

	logger.Log.Info("Starting Health Tracker API")

	logger.Log.Info("Environment Variables: " + cfg.Env)

	database := db.Connect(cfg)
	if database == nil {
		log.Fatal("Failed to connect to the database")
	}
	if err := database.AutoMigrate(&food.Food{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.JSONMiddleware)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.RecoveryMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "Health Tracker API is running!"})
	})

	foodHandler := food.NewHandlerFactory(database, logger.Log)
	r.Mount("/foods", foodHandler.Routes())

	port := fmt.Sprintf(":%s", cfg.AppPort)
	logger.Log.Info("Server is starting", zap.String("port", port))

	log.Fatal(http.ListenAndServe(port, r))
}
