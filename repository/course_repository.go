package repository

import (
	"GoEcho1/model"

	"gorm.io/gorm"
)

type CourseRepository interface {
	GetAllCourses() ([]model.Course, error)
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db}
}

func (r *courseRepository) GetAllCourses() ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Find(&courses).Error
	return courses, err
}
