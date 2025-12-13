package repository

import (
	"context"
	"backenduas/app/model"
)

type IStudentRepository interface {
	GetAll(ctx context.Context) ([]model.StudentDetail, error)
	GetByID(ctx context.Context, id string) (*model.StudentDetail, error)
	Create(ctx context.Context, s *model.Student) error
	UpdateAdvisor(ctx context.Context, studentID string, advisorID *string) error
	GetStudentIDByUserID(ctx context.Context, userID string) (string, error)
	GetLecturerIDByUserID(ctx context.Context, userID string) (string, error)
	IsStudentUnderAdvisor(ctx context.Context, studentID, lecturerID string) (bool, error)

	GetAllStudentIDs(ctx context.Context) ([]string, error)
	GetStudentsByAdvisor(ctx context.Context, advisorID string) ([]string, error)
	GetStudentsByIDs(ctx context.Context, ids []string) (map[string]model.StudentDetail, error)

	
}
