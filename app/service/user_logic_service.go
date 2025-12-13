package service

import (
	"context"
	"errors"

	"backenduas/app/model"
	"backenduas/app/repository"
)

type UserLogicService struct {
	repo repository.IUserRepository
}

func NewUserLogicService(repo repository.IUserRepository) *UserLogicService {
	return &UserLogicService{repo}
}

// ================= CREATE =================
func (s *UserLogicService) CreateUserLogic(req model.User) error {
	if req.Username == "" || req.Password == "" {
		return errors.New("data tidak valid")
	}
	return s.repo.Create(context.Background(), &req)
}

// ================= UPDATE =================
func (s *UserLogicService) UpdateUserLogic(req model.User) error {
	if req.ID == "" {
		return errors.New("id wajib diisi")
	}
	return s.repo.Update(context.Background(), &req)
}

// ================= DELETE =================
func (s *UserLogicService) DeleteUserLogic(id string) error {
	if id == "" {
		return errors.New("id kosong")
	}
	return s.repo.Delete(context.Background(), id)
}

// ================= UPDATE ROLE =================
func (s *UserLogicService) UpdateRoleLogic(userID, roleID string) error {
	if roleID == "" {
		return errors.New("role kosong")
	}
	return s.repo.UpdateRole(context.Background(), userID, roleID)
}

// ================= GET =================
func (s *UserLogicService) GetUserByIDLogic(id string) (*model.User, error) {
	return s.repo.GetByID(context.Background(), id)
}

func (s *UserLogicService) GetAllUsersLogic() ([]model.User, error) {
	return s.repo.GetAll(context.Background())
}
