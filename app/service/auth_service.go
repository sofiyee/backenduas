package service

import (
	"context"
	"os"
	"time"

	"backenduas/app/model"
	"backenduas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo}
}

// =============================
//      LOGIN
// =============================
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad Request"})
	}

	user, err := s.repo.FindByUsername(context.Background(), req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "username atau password salah"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"error": "username atau password salah"})
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"role_id":   user.RoleID,
		"role_name": user.RoleName,
		"full_name": user.FullName,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(os.Getenv("API_KEY")))

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"token": signed,
			"user": fiber.Map{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"role_id":   user.RoleID,
				"role_name": user.RoleName,
			},
		},
	})
}

// =============================
//      PROFILE
// =============================
func (s *AuthService) Profile(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)

	userID := claims["user_id"].(string)
	user, err := s.repo.FindByID(context.Background(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"data": user,
	})
}

// =============================
//      LOGOUT
// =============================
func (s *AuthService) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Logout berhasil",
	})
}

// =============================
//      REFRESH TOKEN
// =============================
func (s *AuthService) Refresh(c *fiber.Ctx) error {
	oldClaims := c.Locals("user").(jwt.MapClaims)

	newClaims := jwt.MapClaims{
		"user_id":   oldClaims["user_id"],
		"role_id":   oldClaims["role_id"],
		"role_name": oldClaims["role_name"],
		"full_name": oldClaims["full_name"],
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	signed, _ := token.SignedString([]byte(os.Getenv("API_KEY")))

	return c.JSON(fiber.Map{
		"token": signed,
	})
}
