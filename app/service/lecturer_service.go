package service

import (
	"context"
	"backenduas/app/repository"
	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	repo *repository.LecturerRepository
}

func NewLecturerService(repo *repository.LecturerRepository) *LecturerService {
	return &LecturerService{repo}
}

// GET /lecturers
func (s *LecturerService) GetAll(c *fiber.Ctx) error {
	data, err := s.repo.GetAll(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": data})
}

// GET /lecturers/:id/advisees
func (s *LecturerService) GetAdvisees(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := s.repo.GetAdvisees(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": data})
}
