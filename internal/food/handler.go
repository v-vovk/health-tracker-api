package food

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Food struct {
	ID        string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) GetFoods(w http.ResponseWriter, r *http.Request) {
	var foods []Food
	h.DB.Find(&foods)
	json.NewEncoder(w).Encode(foods)
}

func (h *Handler) CreateFood(w http.ResponseWriter, r *http.Request) {
	var food Food
	if err := json.NewDecoder(r.Body).Decode(&food); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	h.DB.Create(&food)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(food)
}

func (h *Handler) GetFoodByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var food Food
	if err := h.DB.First(&food, "id = ?", id).Error; err != nil {
		http.Error(w, "Food not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(food)
}

func (h *Handler) UpdateFood(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var food Food
	if err := h.DB.First(&food, "id = ?", id).Error; err != nil {
		http.Error(w, "Food not found", http.StatusNotFound)
		return
	}

	var updatedData Food
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	food.Name = updatedData.Name
	h.DB.Save(&food)
	json.NewEncoder(w).Encode(food)
}

func (h *Handler) DeleteFood(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.DB.Delete(&Food{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Food not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
