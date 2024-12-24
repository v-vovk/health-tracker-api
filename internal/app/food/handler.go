package food

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/v-vovk/health-tracker-api/internal/infra/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Handler struct {
	Service   *Service
	Validator *validator.Validate
	Logger    *zap.Logger
}

func NewHandler(service *Service, validator *validator.Validate, logger *zap.Logger) *Handler {
	return &Handler{
		Service:   service,
		Validator: validator,
		Logger:    logger,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		} else {
			errors.WriteHTTPError(w, http.StatusBadRequest, "Invalid 'limit' parameter")
			h.Logger.Warn("Invalid 'limit' parameter", zap.String("limit", l))
			return
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		} else {
			errors.WriteHTTPError(w, http.StatusBadRequest, "Invalid 'offset' parameter")
			h.Logger.Warn("Invalid 'offset' parameter", zap.String("offset", o))
			return
		}
	}

	foods, total, err := h.Service.GetAllFoods(limit, offset)
	if err != nil {
		errors.WriteHTTPError(w, http.StatusInternalServerError, "Error retrieving foods")
		h.Logger.Error("Error retrieving foods", zap.Error(err))
		return
	}

	response := map[string]interface{}{
		"data":     foods,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
		"returned": len(foods),
	}

	h.Logger.Info("Retrieved foods", zap.Int("limit", limit), zap.Int("offset", offset), zap.Int("returned", len(foods)))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var food Food
	if err := json.NewDecoder(r.Body).Decode(&food); err != nil {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Invalid JSON input")
		h.Logger.Warn("Invalid JSON input", zap.Error(err))
		return
	}

	if err := h.Validator.Struct(food); err != nil {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Validation failed")
		h.Logger.Warn("Validation failed", zap.Error(err))
		return
	}

	if err := h.Service.CreateFood(&food); err != nil {
		errors.WriteHTTPError(w, http.StatusInternalServerError, "Error creating food")
		h.Logger.Error("Error creating food", zap.Error(err))
		return
	}

	h.Logger.Info("Created new food", zap.String("id", food.ID), zap.String("name", food.Name))
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(food); err != nil {
		h.Logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Missing food ID")
		h.Logger.Warn("Missing food ID in request")
		return
	}

	food, err := h.Service.GetFoodByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.WriteHTTPError(w, http.StatusNotFound, "Food not found")
			h.Logger.Warn("Food not found", zap.String("id", id))
			return
		}
		errors.WriteHTTPError(w, http.StatusInternalServerError, "Error retrieving food")
		h.Logger.Error("Error retrieving food", zap.String("id", id), zap.Error(err))
		return
	}

	h.Logger.Info("Retrieved food", zap.String("id", food.ID), zap.String("name", food.Name))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(food); err != nil {
		h.Logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Missing food ID")
		h.Logger.Warn("Missing food ID in request")
		return
	}

	var updatedData Food
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Invalid JSON input")
		h.Logger.Warn("Invalid JSON input for update", zap.Error(err))
		return
	}

	if err := h.Validator.Struct(updatedData); err != nil {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Validation failed")
		h.Logger.Warn("Validation failed for update", zap.Error(err))
		return
	}

	updatedData.ID = id
	if err := h.Service.UpdateFood(&updatedData); err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.WriteHTTPError(w, http.StatusNotFound, "Food not found")
			h.Logger.Warn("Food not found", zap.String("id", id))
			return
		}
		errors.WriteHTTPError(w, http.StatusInternalServerError, "Error updating food")
		h.Logger.Error("Error updating food", zap.String("id", id), zap.Error(err))
		return
	}

	h.Logger.Info("Updated food", zap.String("id", updatedData.ID), zap.String("name", updatedData.Name))
	if err := json.NewEncoder(w).Encode(updatedData); err != nil {
		h.Logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Missing food ID")
		h.Logger.Warn("Missing food ID in request")
		return
	}

	if err := h.Service.DeleteFood(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.WriteHTTPError(w, http.StatusNotFound, "Food not found")
			h.Logger.Warn("Food not found", zap.String("id", id))
			return
		}
		errors.WriteHTTPError(w, http.StatusInternalServerError, "Error deleting food")
		h.Logger.Error("Error deleting food", zap.String("id", id), zap.Error(err))
		return
	}

	h.Logger.Info("Deleted food", zap.String("id", id))
	w.WriteHeader(http.StatusNoContent)
}
