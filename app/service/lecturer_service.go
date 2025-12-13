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
// GetAllLecturers godoc
// @Summary Get all lecturers
// @Description Menampilkan daftar seluruh dosen wali
// @Tags Lecturer
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /lecturers [get]
// GET /lecturers
func (s *LecturerService) GetAll(c *fiber.Ctx) error {
	data, err := s.repo.GetAll(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": data})
}

// GET /lecturers/:id/advisees
// GetLecturerAdvisees godoc
// @Summary Get lecturer advisees
// @Description Menampilkan daftar mahasiswa bimbingan dari dosen wali
// @Tags Lecturer
// @Security BearerAuth
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /lecturers/{id}/advisees [get]
func (s *LecturerService) GetAdvisees(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := s.repo.GetAdvisees(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": data})
}
