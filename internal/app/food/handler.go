package food

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/v-vovk/health-tracker-api/pkg/errors"
	"github.com/v-vovk/health-tracker-api/pkg/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Food struct {
	ID        string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `json:"name" gorm:"not null" validate:"required,min=2,max=50"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Handler struct {
	DB        *gorm.DB
	Validator *validator.Validate
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Extract limit and offset from query parameters
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// Default values if parameters are not provided
	limitInt := 10 // Default limit
	offsetInt := 0 // Default offset

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			limitInt = l
		} else {
			http.Error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
			log.Printf("Invalid 'limit' parameter: %s", limit)
			return
		}
	}

	if offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			offsetInt = o
		} else {
			http.Error(w, "Invalid 'offset' parameter", http.StatusBadRequest)
			log.Printf("Invalid 'offset' parameter: %s", offset)
			return
		}
	}

	// Fetch records with pagination
	var foods []Food
	if err := h.DB.Limit(limitInt).Offset(offsetInt).Find(&foods).Error; err != nil {
		http.Error(w, "Error retrieving foods", http.StatusInternalServerError)
		log.Printf("Error retrieving foods: %v", err)
		return
	}

	// Fetch total count for pagination info
	var total int64
	h.DB.Model(&Food{}).Count(&total)

	// Create response with pagination metadata
	response := map[string]interface{}{
		"data":     foods,
		"total":    total,
		"limit":    limitInt,
		"offset":   offsetInt,
		"returned": len(foods),
	}

	log.Printf("Retrieved %d foods with limit %d and offset %d", len(foods), limitInt, offsetInt)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var food Food
	if err := json.NewDecoder(r.Body).Decode(&food); err != nil {
		errors.JSONError(w, "Invalid JSON input", http.StatusBadRequest)
		logger.Log.Warn("Invalid JSON input", zap.Error(err))
		return
	}

	// Validate input
	if err := h.Validator.Struct(food); err != nil {
		errors.ValidationErrors(w, err, http.StatusBadRequest)
		logger.Log.Warn("Validation failed", zap.Error(err))
		return
	}

	// Save to database
	if err := h.DB.Create(&food).Error; err != nil {
		errors.JSONError(w, "Error creating food", http.StatusInternalServerError)
		logger.Log.Error("Error creating food", zap.Error(err))
		return
	}

	logger.Log.Info("Created new food", zap.String("id", food.ID), zap.String("name", food.Name))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(food)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing food ID", http.StatusBadRequest)
		log.Println("Missing food ID in request")
		return
	}

	var food Food
	if err := h.DB.First(&food, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Food not found", http.StatusNotFound)
			log.Printf("Food not found with ID: %s", id)
			return
		}
		http.Error(w, "Error retrieving food", http.StatusInternalServerError)
		log.Printf("Error retrieving food with ID %s: %v", id, err)
		return
	}
	log.Printf("Retrieved food: %+v", food)
	json.NewEncoder(w).Encode(food)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errors.JSONError(w, "Missing food ID", http.StatusBadRequest)
		log.Println("Missing food ID in request")
		return
	}

	var food Food
	if err := h.DB.First(&food, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.JSONError(w, "Food not found", http.StatusNotFound)
			log.Printf("Food not found with ID: %s", id)
			return
		}
		errors.JSONError(w, "Error retrieving food", http.StatusInternalServerError)
		log.Printf("Error retrieving food with ID %s: %v", id, err)
		return
	}

	var updatedData Food
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		errors.JSONError(w, "Invalid JSON input", http.StatusBadRequest)
		log.Printf("Invalid JSON input for updating food: %v", err)
		return
	}

	// Validate input
	if err := h.Validator.Struct(updatedData); err != nil {
		errors.ValidationErrors(w, err, http.StatusBadRequest)
		log.Printf("Validation failed for food update: %v", err)
		return
	}

	// Update record
	food.Name = updatedData.Name
	if err := h.DB.Save(&food).Error; err != nil {
		errors.JSONError(w, "Error updating food", http.StatusInternalServerError)
		log.Printf("Error updating food with ID %s: %v", id, err)
		return
	}

	log.Printf("Updated food: %+v", food)
	json.NewEncoder(w).Encode(food)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing food ID", http.StatusBadRequest)
		log.Println("Missing food ID in request")
		return
	}

	if err := h.DB.Delete(&Food{}, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Food not found", http.StatusNotFound)
			log.Printf("Food not found with ID: %s", id)
			return
		}
		http.Error(w, "Error deleting food", http.StatusInternalServerError)
		log.Printf("Error deleting food with ID %s: %v", id, err)
		return
	}
	log.Printf("Deleted food with ID: %s", id)
	w.WriteHeader(http.StatusNoContent)
}
