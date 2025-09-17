package repository

import (
	"GoEcho1/model"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KRSRepository interface {
	GetStudentKRS(studentID uint) ([]model.KRS, error)
	GetAllKRS() ([]model.KRS, error)
	AddCourseToKRS(studentID, courseID uint) error
}

type krsRepository struct {
	db *gorm.DB
}

func NewKRSRepository(db *gorm.DB) KRSRepository {
	return &krsRepository{db}
}

func (r *krsRepository) GetStudentKRS(studentID uint) ([]model.KRS, error) {
	var krsList []model.KRS
	err := r.db.Preload("Course").Where("user_id = ?", studentID).Find(&krsList).Error
	return krsList, err
}

func (r *krsRepository) GetAllKRS() ([]model.KRS, error) {
	var krsList []model.KRS
	err := r.db.Preload("User").Preload("Course").Find(&krsList).Error
	return krsList, err
}

// Fungsi ini menggunakan transaksi dan row-level lock untuk mencegah race condition.
func (r *krsRepository) AddCourseToKRS(studentID, courseID uint) error {
	// Start db transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		var course model.Course
		// KUNCI: Mengunci baris (row) mata kuliah yang akan diambil.
		// Klausa `clause.Locking{Strength: "UPDATE"}` akan menghasilkan SQL `SELECT ... FOR UPDATE`.
		// Transaksi lain yang mencoba mengakses baris ini akan menunggu hingga transaksi ini selesai (commit/rollback).
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&course, courseID).Error; err != nil {
			return errors.New("mata kuliah tidak ditemukan")
		}

		// Tanpa klausa locking
		// if err := tx.First(&course, courseID).Error; err != nil {
		//     return errors.New("mata kuliah tidak ditemukan")
		// }

		// Cek kuota
		if course.Kuota <= 0 {
			return errors.New("kuota mata kuliah sudah habis")
		}

		// Cek apakah mahasiswa sudah mengambil matkul ini
		var existingKRS int64
		tx.Model(&model.KRS{}).Where("user_id = ? AND course_id = ?", studentID, courseID).Count(&existingKRS)
		if existingKRS > 0 {
			return errors.New("anda sudah mengambil mata kuliah ini")
		}

		// Kurangi kuota
		course.Kuota = course.Kuota - 1
		if err := tx.Save(&course).Error; err != nil {
			return err
		}

		// Tambahkan ke tabel KRS
		newKRS := model.KRS{
			UserID:   studentID,
			CourseID: courseID,
		}
		if err := tx.Create(&newKRS).Error; err != nil {
			return err
		}

		// Jika semua berhasil, commit transaksi
		return nil
	}) // Transaksi akan di-rollback otomatis jika ada error yang di-return
}
