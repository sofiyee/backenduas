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

/* =============================
   LOGIN (WITH REFRESH TOKEN)
============================= */
// Login godoc
// @Summary Login user
// @Description Login user dan mengembalikan access token & refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /auth/login [post]
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

	perms, _ := s.repo.GetPermissionsByRole(context.Background(), user.RoleID)

	// ACCESS TOKEN
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"role_id":   user.RoleID,
		"role_name": user.RoleName,
		"full_name": user.FullName,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, _ := access.SignedString([]byte(os.Getenv("API_KEY")))

	// REFRESH TOKEN
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, _ := refreshJWT.SignedString([]byte(os.Getenv("API_KEY")))

	// SAVE REFRESH TOKEN → DB
	s.repo.SaveRefreshToken(context.Background(), user.ID, refreshToken, time.Now().Add(7*24*time.Hour))

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        accessToken,
			"refreshToken": refreshToken,
			"user": fiber.Map{
				"id":          user.ID,
				"username":    user.Username,
				"fullName":    user.FullName,
				"role":        user.RoleName,
				"permissions": perms,
			},
		},
	})
}

/* =============================
   PROFILE
============================= */
// Profile godoc
// @Summary Get user profile
// @Description Get profile of currently logged-in user
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /auth/profile [get]
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

/* =============================
   LOGOUT (REVOKE REFRESH TOKEN)
============================= */
// Logout godoc
// @Summary Logout user
// @Description Logout user dan mencabut refresh token
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Param X-Refresh-Token header string true "Refresh Token"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /auth/logout [post]
func (s *AuthService) Logout(c *fiber.Ctx) error {
	refreshToken := c.Get("X-Refresh-Token")
	if refreshToken == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Refresh token tidak ditemukan"})
	}

	s.repo.DeleteRefreshToken(context.Background(), refreshToken)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logout berhasil. Token dicabut.",
	})
}

/* =============================
   REFRESH TOKEN
============================= */
// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate access token baru menggunakan refresh token
// @Tags Auth
// @Produce json
// @Param X-Refresh-Token header string true "Refresh Token"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /auth/refresh [post]
func (s *AuthService) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Get("X-Refresh-Token")
	if refreshToken == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Refresh token tidak ditemukan"})
	}

	valid, _ := s.repo.IsRefreshTokenValid(context.Background(), refreshToken)
	if !valid {
		return c.Status(401).JSON(fiber.Map{"error": "Refresh token sudah dicabut"})
	}

	// Parse refresh token
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_KEY")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Refresh token invalid"})
	}

	claims := token.Claims.(jwt.MapClaims)

	newClaims := jwt.MapClaims{
		"user_id":   claims["user_id"],
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	signed, _ := newToken.SignedString([]byte(os.Getenv("API_KEY")))

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": signed,
		},
	})
}
