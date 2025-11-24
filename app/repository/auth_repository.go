package repository

import (
	"context"
	"errors"

	"backenduas/app/model"
	"backenduas/database"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {

	query := `
		SELECT 
    u.id,
    u.username,
    u.email,
    u.password_hash,
    u.full_name,
    u.role_id,
    r.name AS role_name,
    u.is_active,
    u.created_at,
    u.updated_at
	FROM users u
	JOIN roles r ON r.id = u.role_id
	WHERE u.username = $1

	`

	row := database.DB.QueryRow(ctx, query, username)

	var u model.User

	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.RoleID,
		&u.RoleName,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)


	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return &u, nil
}

func (r *AuthRepository) FindByID(ctx context.Context, id string) (*model.User, error) {

	query := `
		SELECT 
			u.id,
			u.username,
			u.email,
			u.full_name,
			u.role_id,
			r.name AS role_name,
			u.is_active,
			u.created_at,
			u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`

	row := database.DB.QueryRow(ctx, query, id)

	var u model.User

	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.FullName,
		&u.RoleID,
		&u.RoleName,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return &u, nil
}

