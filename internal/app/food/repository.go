package food

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAllFoods(limit, offset int) ([]Food, int64, error) {
	var foods []Food
	var total int64

	if err := r.DB.Limit(limit).Offset(offset).Find(&foods).Error; err != nil {
		return nil, 0, err
	}

	r.DB.Model(&Food{}).Count(&total)
	return foods, total, nil
}

func (r *Repository) GetFoodByID(id string) (*Food, error) {
	var food Food
	if err := r.DB.First(&food, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *Repository) CreateFood(food *Food) error {
	return r.DB.Create(food).Error
}

func (r *Repository) UpdateFood(food *Food) error {
	return r.DB.Save(food).Error
}

func (r *Repository) DeleteFood(id string) error {
	return r.DB.Delete(&Food{}, "id = ?", id).Error
}
