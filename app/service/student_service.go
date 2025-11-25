package service

import (
	"context"
	"backenduas/app/model"
	"backenduas/app/repository"
	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{repo}
}

// GET /students
func (s *StudentService) GetAll(c *fiber.Ctx) error {
	data, err := s.repo.GetAll(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": data})
}

// GET /students/:id
func (s *StudentService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	data, err := s.repo.GetByID(context.Background(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Student tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"data": data})
}

// POST /students
func (s *StudentService) Create(c *fiber.Ctx) error {
	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if req.UserID == "" || req.StudentID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id dan student_id wajib diisi"})
	}

	if err := s.repo.Create(context.Background(), &model.Student{
		UserID:       req.UserID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    req.AdvisorID,
	}); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Student berhasil dibuat"})
}

// PUT /students/:id/advisor
func (s *StudentService) UpdateAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")

	body := struct {
		AdvisorID *string `json:"advisor_id"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := s.repo.UpdateAdvisor(context.Background(), studentID, body.AdvisorID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Advisor berhasil diperbarui"})
}
