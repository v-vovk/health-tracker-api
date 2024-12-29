package food

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllFoods(limit, offset int) ([]Food, int64, error)
	GetFoodByID(id string) (*Food, error)
	CreateFood(food *Food) error
	UpdateFood(food *Food) error
	DeleteFood(id string) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

func (r *repository) GetAllFoods(limit, offset int) ([]Food, int64, error) {
	var foods []Food
	var total int64

	if err := r.DB.Limit(limit).Offset(offset).Find(&foods).Error; err != nil {
		return nil, 0, err
	}

	r.DB.Model(&Food{}).Count(&total)
	return foods, total, nil
}

func (r *repository) GetFoodByID(id string) (*Food, error) {
	var food Food
	if err := r.DB.First(&food, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *repository) CreateFood(food *Food) error {
	return r.DB.Create(food).Error
}

func (r *repository) UpdateFood(food *Food) error {
	return r.DB.Save(food).Error
}

func (r *repository) DeleteFood(id string) error {
	return r.DB.Delete(&Food{}, "id = ?", id).Error
}
