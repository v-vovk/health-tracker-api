package group

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(limit, offset int) ([]Group, int64, error)
	GetByID(id string) (*Group, error)
	Create(group *Group) error
	Update(group *Group) error
	Delete(id string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) GetAll(limit, offset int) ([]Group, int64, error) {
	var groups []Group
	var total int64

	if err := r.db.Model(&Group{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(limit).Offset(offset).Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

func (r *repositoryImpl) GetByID(id string) (*Group, error) {
	var group Group
	if err := r.db.First(&group, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *repositoryImpl) Create(group *Group) error {
	return r.db.Create(group).Error
}

func (r *repositoryImpl) Update(group *Group) error {
	return r.db.Save(group).Error
}

func (r *repositoryImpl) Delete(id string) error {
	return r.db.Delete(&Group{}, "id = ?", id).Error
}
