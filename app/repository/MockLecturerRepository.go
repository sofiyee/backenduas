package repository

import (
	"context"
	"errors"

	"backenduas/app/model"
)

type MockLecturerRepository struct {
	Lecturers map[string]model.LecturerDetail
	Advisees  map[string][]model.Advisee // lecturerID -> advisees
}

func NewMockLecturerRepository() *MockLecturerRepository {
	return &MockLecturerRepository{
		Lecturers: make(map[string]model.LecturerDetail),
		Advisees:  make(map[string][]model.Advisee),
	}
}

func (m *MockLecturerRepository) GetAll(ctx context.Context) ([]model.LecturerDetail, error) {
	var res []model.LecturerDetail
	for _, v := range m.Lecturers {
		res = append(res, v)
	}
	return res, nil
}

func (m *MockLecturerRepository) GetAdvisees(ctx context.Context, lecturerID string) ([]model.Advisee, error) {
	if data, ok := m.Advisees[lecturerID]; ok {
		return data, nil
	}
	return nil, errors.New("lecturer not found")
}
