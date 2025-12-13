package repository

import (
	"context"
	"errors"

	"backenduas/app/model"
)

type MockStudentRepository struct {
	// data mock
	Students       map[string]model.StudentDetail      // studentID -> detail
	UserToStudent  map[string]string                   // userID -> studentID
	UserToLecturer map[string]string                   // userID -> lecturerID
	AdvisorMap     map[string]string                   // studentID -> lecturerID
}

func NewMockStudentRepository() *MockStudentRepository {
	return &MockStudentRepository{
		Students:       make(map[string]model.StudentDetail),
		UserToStudent:  make(map[string]string),
		UserToLecturer: make(map[string]string),
		AdvisorMap:     make(map[string]string),
	}
}



func (m *MockStudentRepository) GetAll(ctx context.Context) ([]model.StudentDetail, error) {
	out := []model.StudentDetail{}
	for _, v := range m.Students {
		out = append(out, v)
	}
	return out, nil
}

func (m *MockStudentRepository) GetByID(ctx context.Context, id string) (*model.StudentDetail, error) {
	if v, ok := m.Students[id]; ok {
		cp := v
		return &cp, nil
	}
	return nil, errors.New("student not found")
}

func (m *MockStudentRepository) Create(ctx context.Context, s *model.Student) error {
	m.Students[s.StudentID] = model.StudentDetail{
		ID:           s.StudentID,
		UserID:       s.UserID,
		StudentID:    s.StudentID,
		ProgramStudy: s.ProgramStudy,
		AcademicYear: s.AcademicYear,
		AdvisorID:    s.AdvisorID,
	}
	return nil
}

func (m *MockStudentRepository) UpdateAdvisor(ctx context.Context, studentID string, advisorID *string) error {
	if v, ok := m.Students[studentID]; ok {
		v.AdvisorID = advisorID
		m.Students[studentID] = v
		return nil
	}
	return errors.New("student not found")
}

func (m *MockStudentRepository) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	if sid, ok := m.UserToStudent[userID]; ok {
		return sid, nil
	}
	return "", errors.New("student not found")
}

func (m *MockStudentRepository) GetLecturerIDByUserID(ctx context.Context, userID string) (string, error) {
	if lid, ok := m.UserToLecturer[userID]; ok {
		return lid, nil
	}
	return "", errors.New("lecturer not found")
}

func (m *MockStudentRepository) GetStudentsByAdvisor(ctx context.Context, advisorID string) ([]string, error) {
	out := []string{}
	for sid, aid := range m.AdvisorMap {
		if aid == advisorID {
			out = append(out, sid)
		}
	}
	return out, nil
}

func (m *MockStudentRepository) IsStudentUnderAdvisor(ctx context.Context, studentID string, advisorID string) (bool, error) {
	if aid, ok := m.AdvisorMap[studentID]; ok {
		return aid == advisorID, nil
	}
	return false, nil
}

func (m *MockStudentRepository) GetAllStudentIDs(ctx context.Context) ([]string, error) {
	ids := []string{}
	for id := range m.Students {
		ids = append(ids, id)
	}
	return ids, nil
}

func (m *MockStudentRepository) GetStudentsByIDs(
	ctx context.Context,
	ids []string,
) (map[string]model.StudentDetail, error) {

	result := map[string]model.StudentDetail{}

	for _, id := range ids {
		if st, ok := m.Students[id]; ok {
			result[id] = st
		}
	}

	return result, nil
}
