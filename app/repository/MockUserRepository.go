package repository

import (
	"context"
	"errors"

	"backenduas/app/model"
)

type MockUserRepository struct {
	Data map[string]model.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Data: make(map[string]model.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, u *model.User) error {
	if u.Username == "" {
		return errors.New("username kosong")
	}
	m.Data[u.ID] = *u
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, u *model.User) error {
	if _, ok := m.Data[u.ID]; !ok {
		return errors.New("user tidak ditemukan")
	}
	m.Data[u.ID] = *u
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.Data[id]; !ok {
		return errors.New("user tidak ditemukan")
	}
	delete(m.Data, id)
	return nil
}

func (m *MockUserRepository) UpdateRole(ctx context.Context, userID string, roleID string) error {
	u, ok := m.Data[userID]
	if !ok {
		return errors.New("user tidak ditemukan")
	}
	u.RoleID = roleID
	m.Data[userID] = u
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	u, ok := m.Data[id]
	if !ok {
		return nil, errors.New("user tidak ditemukan")
	}
	return &u, nil
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	var res []model.User
	for _, u := range m.Data {
		res = append(res, u)
	}
	return res, nil
}
