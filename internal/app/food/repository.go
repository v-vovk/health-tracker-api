package food

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(limit, offset int) ([]Food, int64, error)
	GetByID(id string) (*Food, error)
	Create(food *Food) error
	Update(food *Food) error
	Delete(id string) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

func (r *repository) GetAll(limit, offset int) ([]Food, int64, error) {
	var foods []Food
	var total int64

	if err := r.DB.Limit(limit).Offset(offset).Find(&foods).Error; err != nil {
		return nil, 0, err
	}

	r.DB.Model(&Food{}).Count(&total)
	return foods, total, nil
}

func (r *repository) GetByID(id string) (*Food, error) {
	var food Food
	if err := r.DB.First(&food, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *repository) Create(food *Food) error {
	return r.DB.Create(food).Error
}

func (r *repository) Update(food *Food) error {
	return r.DB.Save(food).Error
}

func (r *repository) Delete(id string) error {
	return r.DB.Delete(&Food{}, "id = ?", id).Error
}
