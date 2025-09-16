package model

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	KodeMK string `gorm:"uniqueIndex;not null"`
	NamaMK string
	SKS    int
	Kuota  int
}
