package db

import (
	"GoEcho1/model"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// Hanya seed jika tabel kosong
	if db.First(&model.User{}).Error != gorm.ErrRecordNotFound {
		log.Println("Database already seeded.")
		return
	}

	log.Println("Seeding database...")

	// Hashing passwords
	hashedPasswordMhs, _ := bcrypt.GenerateFromPassword([]byte("mahasiswa123"), bcrypt.DefaultCost)
	hashedPasswordDosen, _ := bcrypt.GenerateFromPassword([]byte("dosen123"), bcrypt.DefaultCost)

	// Seeding Users
	users := []model.User{
		{NIM: "12345", Nama: "Budi Santoso", Password: string(hashedPasswordMhs), Role: "mahasiswa"},
		{NIM: "67890", Nama: "Ani Yudhoyono", Password: string(hashedPasswordMhs), Role: "mahasiswa"},
		{NIM: "DOSEN01", Nama: "Dr. Retno Wulandari", Password: string(hashedPasswordDosen), Role: "dosen"},
	}
	db.Create(&users)

	// Seeding Courses
	// Kuota dibuat kecil untuk mempermudah tes race condition
	courses := []model.Course{
		{KodeMK: "IF101", NamaMK: "Dasar Pemrograman", SKS: 3, Kuota: 5},
		{KodeMK: "IF102", NamaMK: "Struktur Data", SKS: 3, Kuota: 3},
		{KodeMK: "UM101", NamaMK: "Pendidikan Pancasila", SKS: 2, Kuota: 10},
	}
	db.Create(&courses)

	log.Println("Seeding completed.")
}
