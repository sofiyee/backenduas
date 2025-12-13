package repository

import (
	"context"
	"backenduas/app/model"
)

type IUserRepository interface {
	Create(ctx context.Context, u *model.User) error
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, userID string, roleID string) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
}
