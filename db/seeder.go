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

// ResetAndSeed akan menghapus tabel yang ada, membuat ulang, dan menjalankan seeder.
func ResetAndSeed(db *gorm.DB) error {
	log.Println("Attempting to reset and re-seed database...")

	// 1. Hapus tabel yang ada dengan urutan yang benar (dari yang punya foreign key)
	// Menggunakan Migrator().DropTable lebih aman daripada raw SQL
	err := db.Migrator().DropTable(&model.KRS{}, &model.User{}, &model.Course{})
	if err != nil {
		log.Printf("Failed to drop tables: %v\n", err)
		return err
	}
	log.Println("Old tables dropped successfully.")

	// 2. Buat ulang tabel (AutoMigrate)
	err = db.AutoMigrate(&model.User{}, &model.Course{}, &model.KRS{})
	if err != nil {
		log.Printf("Failed to migrate new tables: %v\n", err)
		return err
	}
	log.Println("Tables migrated successfully.")

	// 3. Panggil fungsi Seed untuk mengisi data awal
	// Kita memanggil seeder secara langsung, melewati pengecekan data
	SeedForced(db)

	log.Println("Database has been successfully reset and re-seeded.")
	return nil
}

// SeedForced adalah varian dari Seed yang tidak melakukan pengecekan data
func SeedForced(db *gorm.DB) {
	log.Println("Force seeding database...")

	hashedPasswordMhs, _ := bcrypt.GenerateFromPassword([]byte("mahasiswa123"), bcrypt.DefaultCost)
	hashedPasswordDosen, _ := bcrypt.GenerateFromPassword([]byte("dosen123"), bcrypt.DefaultCost)

	users := []model.User{
		{NIM: "12345", Nama: "Budi Santoso", Password: string(hashedPasswordMhs), Role: "mahasiswa"},
		{NIM: "67890", Nama: "Ani Yudhoyono", Password: string(hashedPasswordMhs), Role: "mahasiswa"},
		{NIM: "DOSEN01", Nama: "Dr. Retno Wulandari", Password: string(hashedPasswordDosen), Role: "dosen"},
	}
	db.Create(&users)

	courses := []model.Course{
		{KodeMK: "IF101", NamaMK: "Dasar Pemrograman", SKS: 3, Kuota: 5},
		{KodeMK: "IF102", NamaMK: "Struktur Data", SKS: 3, Kuota: 3},
		{KodeMK: "UM101", NamaMK: "Pendidikan Pancasila", SKS: 2, Kuota: 10},
	}
	db.Create(&courses)

	log.Println("Force seeding completed.")
}
