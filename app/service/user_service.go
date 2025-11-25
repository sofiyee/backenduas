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
