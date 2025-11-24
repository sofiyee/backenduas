package repository

import (
	"context"
	"backenduas/app/model"
	"backenduas/database"
)

type UserRepository struct {}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	rows, err := database.DB.Query(ctx, `
		SELECT u.id, u.username, u.email, u.full_name,
			   u.role_id, r.name AS role_name,
			   u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName,
			&u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT u.id, u.username, u.email, u.full_name,
			   u.role_id, r.name AS role_name,
			   u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`, id)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName,
		&u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	_, err := database.DB.Exec(ctx, `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, true, NOW(), NOW())
	`, u.Username, u.Email, u.PasswordHash, u.FullName, u.RoleID)

	return err
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	_, err := database.DB.Exec(ctx, `
		UPDATE users SET username=$1, email=$2, full_name=$3, updated_at=NOW()
		WHERE id = $4
	`, u.Username, u.Email, u.FullName, u.ID)

	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID string, roleID string) error {
	_, err := database.DB.Exec(ctx, `
		UPDATE users SET role_id=$1, updated_at=NOW()
		WHERE id = $2
	`, roleID, userID)
	return err
}
