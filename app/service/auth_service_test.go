package service

import (
	"testing"

	"backenduas/app/repository"
)

func TestLoginLogic_Success(t *testing.T) {
	repo := repository.NewMockAuthRepository()
	repo.SeedUser("u1", "admin", "admin123", "Admin")

	svc := NewAuthLogicService(repo)

	_, _, user, err := svc.LoginLogic("admin", "admin123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Username != "admin" {
		t.Fatalf("wrong user returned")
	}
}

func TestLoginLogic_WrongPassword(t *testing.T) {
	repo := repository.NewMockAuthRepository()
	repo.SeedUser("u1", "admin", "admin123", "Admin")

	svc := NewAuthLogicService(repo)

	_, _, _, err := svc.LoginLogic("admin", "salah")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestLoginLogic_UserNotFound(t *testing.T) {
	repo := repository.NewMockAuthRepository()
	svc := NewAuthLogicService(repo)

	_, _, _, err := svc.LoginLogic("unknown", "123")
	if err == nil {
		t.Fatalf("expected error")
	}
}
