package group

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/v-vovk/health-tracker-api/internal/infra/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	Service   Service
	Validator *validator.Validate
	Logger    *zap.Logger
}

func NewHandler(service Service, validator *validator.Validate, logger *zap.Logger) *Handler {
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

	groups, total, err := h.Service.GetAll(limit, offset)
	if err != nil {
		errors.WriteHTTPError(w, http.StatusInternalServerError, "Error retrieving groups")
		h.Logger.Error("Error retrieving groups", zap.Error(err))
		return
	}

	response := map[string]interface{}{
		"data":     groups,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
		"returned": len(groups),
	}

	h.Logger.Info("Retrieved groups", zap.Int("limit", limit), zap.Int("offset", offset), zap.Int("returned", len(groups)))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var group Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Invalid JSON input")
		h.Logger.Warn("Invalid JSON input", zap.Error(err))
		return
	}

	if err := h.Validator.Struct(group); err != nil {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Validation failed")
		h.Logger.Warn("Validation failed", zap.Error(err))
		return
	}

	if err := h.Service.Create(&group); err != nil {
		errors.WriteHTTPError(w, http.StatusInternalServerError, "Error creating group")
		h.Logger.Error("Error creating group", zap.Error(err))
		return
	}

	h.Logger.Info("Created new group", zap.String("id", group.ID), zap.String("name", group.Name))
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(group); err != nil {
		h.Logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Missing group ID")
		h.Logger.Warn("Missing group ID in request")
		return
	}

	group, err := h.Service.GetByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			errors.WriteHTTPError(w, http.StatusNotFound, "Group not found")
			h.Logger.Warn("Group not found", zap.String("id", id))
			return
		}
		errors.WriteHTTPError(w, http.StatusInternalServerError, "Error retrieving group")
		h.Logger.Error("Error retrieving group", zap.String("id", id), zap.Error(err))
		return
	}

	h.Logger.Info("Retrieved group", zap.String("id", group.ID), zap.String("name", group.Name))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(group); err != nil {
		h.Logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Missing group ID")
		h.Logger.Warn("Missing group ID in request")
		return
	}

	var updatedData Group
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
	if err := h.Service.Update(&updatedData); err != nil {
		if err.Error() == "record not found" {
			errors.WriteHTTPError(w, http.StatusNotFound, "Group not found")
			h.Logger.Warn("Group not found", zap.String("id", id))
			return
		}
		errors.WriteHTTPError(w, http.StatusInternalServerError, "Error updating group")
		h.Logger.Error("Error updating group", zap.String("id", id), zap.Error(err))
		return
	}

	h.Logger.Info("Updated group", zap.String("id", updatedData.ID), zap.String("name", updatedData.Name))
	if err := json.NewEncoder(w).Encode(updatedData); err != nil {
		h.Logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errors.WriteHTTPError(w, http.StatusBadRequest, "Missing group ID")
		h.Logger.Warn("Missing group ID in request")
		return
	}

	if err := h.Service.Delete(id); err != nil {
		if err.Error() == "record not found" {
			errors.WriteHTTPError(w, http.StatusNotFound, "Group not found")
			h.Logger.Warn("Group not found", zap.String("id", id))
			return
		}
		errors.WriteHTTPError(w, http.StatusInternalServerError, "Error deleting group")
		h.Logger.Error("Error deleting group", zap.String("id", id), zap.Error(err))
		return
	}

	h.Logger.Info("Deleted group", zap.String("id", id))
	w.WriteHeader(http.StatusNoContent)
}
