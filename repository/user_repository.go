package repository

import (
	"GoEcho1/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByNIM(nim string) (model.User, error)
	GetUserByID(id uint) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByNIM(nim string) (model.User, error) {
	var user model.User
	err := r.db.Where("nim = ?", nim).First(&user).Error
	return user, err
}

func (r *userRepository) GetUserByID(id uint) (model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return user, err
}
