package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	NIM      string `gorm:"uniqueIndex;not null"`
	Nama     string
	Password string
	Role     string `gorm:"default:'mahasiswa'"`
}
