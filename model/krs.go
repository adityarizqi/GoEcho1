package model

import "gorm.io/gorm"

type KRS struct {
	gorm.Model
	UserID   uint
	User     User `gorm:"foreignKey:UserID"`
	CourseID uint
	Course   Course `gorm:"foreignKey:CourseID"`
}
