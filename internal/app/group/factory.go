package group

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewHandlerFactory(db *gorm.DB, logger *zap.Logger) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	validator := validator.New()
	groupLogger := logger.Named("GroupHandler")

	return NewHandler(service, validator, groupLogger)
}
