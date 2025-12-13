package repository

import (
	"context"
	"errors"
	"time"

	"backenduas/app/model"
	"golang.org/x/crypto/bcrypt"
)

type MockAuthRepository struct {
	Users map[string]model.User
}

func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{
		Users: make(map[string]model.User),
	}
}

func (m *MockAuthRepository) SeedUser(id, username, password, role string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	m.Users[username] = model.User{
		ID:           id,
		Username:     username,
		PasswordHash: string(hash),
		RoleName:     role,
		RoleID:       "role-1",
		IsActive:     true,
	}
}

func (m *MockAuthRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	u, ok := m.Users[username]
	if !ok {
		return nil, errors.New("user tidak ditemukan")
	}
	return &u, nil
}

func (m *MockAuthRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	for _, u := range m.Users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user tidak ditemukan")
}

func (m *MockAuthRepository) GetPermissionsByRole(ctx context.Context, roleID string) ([]string, error) {
	return []string{"read", "write"}, nil
}

func (m *MockAuthRepository) SaveRefreshToken(ctx context.Context, userID, token string, exp time.Time) error {
	return nil
}

func (m *MockAuthRepository) IsRefreshTokenValid(ctx context.Context, token string) (bool, error) {
	return true, nil
}

func (m *MockAuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return nil
}
