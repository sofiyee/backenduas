package repository

import (
	"context"
	"backenduas/app/model"
)

type ILecturerRepository interface {
	GetAll(ctx context.Context) ([]model.LecturerDetail, error)
	GetAdvisees(ctx context.Context, lecturerID string) ([]model.Advisee, error)
}
