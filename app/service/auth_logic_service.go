package service

import (
	"context"
	"errors"
	"os"
	"time"

	"backenduas/app/model"
	"backenduas/app/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthLogicService struct {
	repo repository.IAuthRepository
}

func NewAuthLogicService(repo repository.IAuthRepository) *AuthLogicService {
	return &AuthLogicService{repo}
}

func (s *AuthLogicService) LoginLogic(
	username string,
	password string,
) (string, string, *model.User, error) {

	user, err := s.repo.FindByUsername(context.Background(), username)
	if err != nil {
		return "", "", nil, errors.New("username atau password salah")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", "", nil, errors.New("username atau password salah")
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"role_name": user.RoleName,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, _ := access.SignedString([]byte(os.Getenv("API_KEY")))

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, _ := refresh.SignedString([]byte(os.Getenv("API_KEY")))

	s.repo.SaveRefreshToken(context.Background(), user.ID, refreshToken, time.Now().Add(7*24*time.Hour))

	return accessToken, refreshToken, user, nil
}
