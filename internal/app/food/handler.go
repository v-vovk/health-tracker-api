package food

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/v-vovk/health-tracker-api/internal/infra/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Food struct {
	ID        string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `json:"name" gorm:"not null" validate:"required,min=1,max=50"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Handler struct {
	DB        *gorm.DB
	Validator *validator.Validate
	Logger    *zap.Logger
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	limitInt := 10
	offsetInt := 0

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			limitInt = l
		} else {
			errors.JSONError(w, "Invalid 'limit' parameter", http.StatusBadRequest)
			h.Logger.Warn("Invalid 'limit' parameter", zap.String("limit", limit))
			return
		}
	}

	if offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			offsetInt = o
		} else {
			errors.JSONError(w, "Invalid 'offset' parameter", http.StatusBadRequest)
			h.Logger.Warn("Invalid 'offset' parameter", zap.String("offset", offset))
			return
		}
	}

	var foods []Food
	if err := h.DB.Limit(limitInt).Offset(offsetInt).Find(&foods).Error; err != nil {
		errors.JSONError(w, "Error retrieving foods", http.StatusInternalServerError)
		h.Logger.Error("Error retrieving foods", zap.Error(err))
		return
	}

	var total int64
	h.DB.Model(&Food{}).Count(&total)

	response := map[string]interface{}{
		"data":     foods,
		"total":    total,
		"limit":    limitInt,
		"offset":   offsetInt,
		"returned": len(foods),
	}

	h.Logger.Info("Retrieved foods", zap.Int("limit", limitInt), zap.Int("offset", offsetInt), zap.Int("returned", len(foods)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var food Food
	if err := json.NewDecoder(r.Body).Decode(&food); err != nil {
		errors.JSONError(w, "Invalid JSON input", http.StatusBadRequest)
		h.Logger.Warn("Invalid JSON input", zap.Error(err))
		return
	}

	if err := h.Validator.Struct(food); err != nil {
		errors.ValidationErrors(w, err, http.StatusBadRequest)
		h.Logger.Warn("Validation failed", zap.Error(err))
		return
	}

	if err := h.DB.Create(&food).Error; err != nil {
		errors.JSONError(w, "Error creating food", http.StatusInternalServerError)
		h.Logger.Error("Error creating food", zap.Error(err))
		return
	}

	h.Logger.Info("Created new food", zap.String("id", food.ID), zap.String("name", food.Name))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(food)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errors.JSONError(w, "Missing food ID", http.StatusBadRequest)
		h.Logger.Warn("Missing food ID in request")
		return
	}

	var food Food
	if err := h.DB.First(&food, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.JSONError(w, "Food not found", http.StatusNotFound)
			h.Logger.Warn("Food not found", zap.String("id", id))
			return
		}
		errors.JSONError(w, "Error retrieving food", http.StatusInternalServerError)
		h.Logger.Error("Error retrieving food", zap.String("id", id), zap.Error(err))
		return
	}

	h.Logger.Info("Retrieved food", zap.String("id", food.ID), zap.String("name", food.Name))
	json.NewEncoder(w).Encode(food)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errors.JSONError(w, "Missing food ID", http.StatusBadRequest)
		h.Logger.Warn("Missing food ID in request")
		return
	}

	var food Food
	if err := h.DB.First(&food, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.JSONError(w, "Food not found", http.StatusNotFound)
			h.Logger.Warn("Food not found", zap.String("id", id))
			return
		}
		errors.JSONError(w, "Error retrieving food", http.StatusInternalServerError)
		h.Logger.Error("Error retrieving food", zap.String("id", id), zap.Error(err))
		return
	}

	var updatedData Food
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		errors.JSONError(w, "Invalid JSON input", http.StatusBadRequest)
		h.Logger.Warn("Invalid JSON input for update", zap.Error(err))
		return
	}

	if err := h.Validator.Struct(updatedData); err != nil {
		errors.ValidationErrors(w, err, http.StatusBadRequest)
		h.Logger.Warn("Validation failed for update", zap.Error(err))
		return
	}

	food.Name = updatedData.Name
	if err := h.DB.Save(&food).Error; err != nil {
		errors.JSONError(w, "Error updating food", http.StatusInternalServerError)
		h.Logger.Error("Error updating food", zap.String("id", id), zap.Error(err))
		return
	}

	h.Logger.Info("Updated food", zap.String("id", food.ID), zap.String("name", food.Name))
	json.NewEncoder(w).Encode(food)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errors.JSONError(w, "Missing food ID", http.StatusBadRequest)
		h.Logger.Warn("Missing food ID in request")
		return
	}

	if err := h.DB.Delete(&Food{}, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.JSONError(w, "Food not found", http.StatusNotFound)
			h.Logger.Warn("Food not found", zap.String("id", id))
			return
		}
		errors.JSONError(w, "Error deleting food", http.StatusInternalServerError)
		h.Logger.Error("Error deleting food", zap.String("id", id), zap.Error(err))
		return
	}

	h.Logger.Info("Deleted food", zap.String("id", id))
	w.WriteHeader(http.StatusNoContent)
}
