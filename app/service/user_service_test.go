package service

import (
	"testing"

	"backenduas/app/model"
	"backenduas/app/repository"
)

func TestCreateUserLogic_Success(t *testing.T) {
	repo := repository.NewMockUserRepository()
	svc := NewUserLogicService(repo)

	err := svc.CreateUserLogic(model.User{
		ID:       "1",
		Username: "sofie",
		Password: "123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateUserLogic_Invalid(t *testing.T) {
	repo := repository.NewMockUserRepository()
	svc := NewUserLogicService(repo)

	err := svc.CreateUserLogic(model.User{})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestUpdateRoleLogic(t *testing.T) {
	repo := repository.NewMockUserRepository()
	repo.Data["1"] = model.User{ID: "1", RoleID: "old"}

	svc := NewUserLogicService(repo)

	err := svc.UpdateRoleLogic("1", "new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.Data["1"].RoleID != "new" {
		t.Fatalf("role not updated")
	}
}

func TestDeleteUserLogic(t *testing.T) {
	repo := repository.NewMockUserRepository()
	repo.Data["1"] = model.User{ID: "1"}

	svc := NewUserLogicService(repo)

	err := svc.DeleteUserLogic("1")
	if err != nil {
		t.Fatalf("unexpected error")
	}
}

func TestGetAllUsersLogic(t *testing.T) {
	repo := repository.NewMockUserRepository()
	repo.Data["1"] = model.User{ID: "1"}
	repo.Data["2"] = model.User{ID: "2"}

	svc := NewUserLogicService(repo)

	users, err := svc.GetAllUsersLogic()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}
