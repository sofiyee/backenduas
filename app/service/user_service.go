package service

import (
	"context"

	"backenduas/app/model"
	"backenduas/app/repository"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

// ===============================
// GET ALL
// ===============================
// GetAll godoc
// @Summary Get all users
// @Description Mengambil daftar seluruh user
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]any
// @Router /users [get]
func (s *UserService) GetAll(c *fiber.Ctx) error {
	users, err := s.repo.GetAll(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// ===============================
// GET BY ID
// ===============================
// GetByID godoc
// @Summary Get user by ID
// @Description Mengambil detail user berdasarkan ID
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /users/{id} [get]
func (s *UserService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.repo.GetByID(context.Background(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	return c.JSON(user)
}

// ===============================
// CREATE USER + AUTO INSERT
// ===============================
// Create godoc
// @Summary Create new user
// @Description Membuat user baru dan otomatis membuat data mahasiswa atau dosen wali sesuai role
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.User true "Create user payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /users [post]
func (s *UserService) Create(c *fiber.Ctx) error {
	var req model.User

	// Parse body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	// Generate UUID untuk user
	req.ID = uuid.NewString()

	// Hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	req.PasswordHash = string(hash)

	// Insert user + auto insert ke students/lecturers
	if err := s.repo.Create(context.Background(), &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User berhasil dibuat"})
}

// ===============================
// UPDATE USER
// ===============================
// Update godoc
// @Summary Update user
// @Description Memperbarui data user (username, email, full name)
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body model.User true "Update user payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /users/{id} [put]
func (s *UserService) Update(c *fiber.Ctx) error {
	var req model.User

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	req.ID = c.Params("id")

	if err := s.repo.Update(context.Background(), &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User berhasil diperbarui"})
}

// ===============================
// DELETE USER
// ===============================
// Delete godoc
// @Summary Delete user
// @Description Menghapus user berdasarkan ID
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]any
// @Router /users/{id} [delete]
func (s *UserService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.repo.Delete(context.Background(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User berhasil dihapus"})
}

// ===============================
// UPDATE ROLE
// ===============================
// UpdateRole godoc
// @Summary Update user role
// @Description Mengubah role user
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body object true "Update role payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /users/{id}/role [put]
func (s *UserService) UpdateRole(c *fiber.Ctx) error {
	userID := c.Params("id")

	body := struct {
		RoleID string `json:"role_id"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := s.repo.UpdateRole(context.Background(), userID, body.RoleID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Role user berhasil diperbarui"})
}
