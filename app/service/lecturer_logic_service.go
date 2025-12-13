package service

import (
	"context"
	"errors"

	"backenduas/app/model"
	"backenduas/app/repository"
)

type LecturerLogicService struct {
	repo repository.ILecturerRepository
}

func NewLecturerLogicService(repo repository.ILecturerRepository) *LecturerLogicService {
	return &LecturerLogicService{repo}
}

// ================= GET ALL =================
func (s *LecturerLogicService) GetAllLecturersLogic() ([]model.LecturerDetail, error) {
	return s.repo.GetAll(context.Background())
}

// ================= GET ADVISEES =================
func (s *LecturerLogicService) GetAdviseesLogic(lecturerID string) ([]model.Advisee, error) {
	if lecturerID == "" {
		return nil, errors.New("lecturer id kosong")
	}
	return s.repo.GetAdvisees(context.Background(), lecturerID)
}
