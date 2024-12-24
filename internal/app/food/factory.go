package food

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewHandlerFactory(db *gorm.DB, log *zap.Logger) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	foodLog := log.Named("FoodHandler")
	return NewHandler(service, validator.New(), foodLog)
}
