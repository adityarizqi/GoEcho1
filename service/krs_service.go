package service

import (
	"GoEcho1/repository"
)

type KRSService interface {
	GetDashboardDataForStudent(studentID uint) (map[string]interface{}, error)
	GetDashboardDataForDosen() (map[string]interface{}, error)
	EnrollCourse(studentID, courseID uint) error
}

type krsService struct {
	krsRepo    repository.KRSRepository
	courseRepo repository.CourseRepository
	userRepo   repository.UserRepository
}

func NewKRSService(krsRepo repository.KRSRepository, courseRepo repository.CourseRepository, userRepo repository.UserRepository) KRSService {
	return &krsService{krsRepo, courseRepo, userRepo}
}

func (s *krsService) GetDashboardDataForStudent(studentID uint) (map[string]interface{}, error) {
	user, err := s.userRepo.GetUserByID(studentID)
	if err != nil {
		return nil, err
	}

	allCourses, err := s.courseRepo.GetAllCourses()
	if err != nil {
		return nil, err
	}

	takenCourses, err := s.krsRepo.GetStudentKRS(studentID)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"User":         user,
		"AllCourses":   allCourses,
		"TakenCourses": takenCourses,
	}

	return data, nil
}

func (s *krsService) GetDashboardDataForDosen() (map[string]interface{}, error) {
	allKRS, err := s.krsRepo.GetAllKRS()
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{
		"KRSList": allKRS,
	}
	return data, nil
}

func (s *krsService) EnrollCourse(studentID, courseID uint) error {
	return s.krsRepo.AddCourseToKRS(studentID, courseID)
}
