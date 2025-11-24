package service

import (
	"context"
	"backenduas/app/model"
	"backenduas/app/repository"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) GetAll(c *fiber.Ctx) error {
	users, err := s.repo.GetAll(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (s *UserService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.repo.GetByID(context.Background(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	return c.JSON(user)
}

func (s *UserService) Create(c *fiber.Ctx) error {
	var req model.User

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	// Wajib hash dari password input
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	req.PasswordHash = string(hash)

	if err := s.repo.Create(context.Background(), &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User berhasil dibuat"})
}


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

func (s *UserService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.repo.Delete(context.Background(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User berhasil dihapus"})
}

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
