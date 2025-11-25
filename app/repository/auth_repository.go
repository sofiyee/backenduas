package repository

import (
	"context"
	"errors"

	"backenduas/app/model"
	"backenduas/database"
	"time"
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

func (r *AuthRepository) GetPermissionsByRole(ctx context.Context, roleID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	rows, err := database.DB.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, userID, token string, exp time.Time) error {
	_, err := database.DB.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token, expired_at)
		VALUES ($1, $2, $3)
	`, userID, token, exp)
	return err
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := database.DB.Exec(ctx, `
		DELETE FROM refresh_tokens WHERE token = $1
	`, token)
	return err
}

func (r *AuthRepository) IsRefreshTokenValid(ctx context.Context, token string) (bool, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT COUNT(*) FROM refresh_tokens
		WHERE token = $1 AND expired_at > NOW()
	`, token)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}


