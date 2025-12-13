package repository

import (
	"context"
	"time"

	"backenduas/app/model"
)

type IAuthRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	GetPermissionsByRole(ctx context.Context, roleID string) ([]string, error)

	SaveRefreshToken(ctx context.Context, userID, token string, exp time.Time) error
	IsRefreshTokenValid(ctx context.Context, token string) (bool, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}
