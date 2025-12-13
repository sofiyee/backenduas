package service

import (
    "backenduas/app/model"
    "backenduas/app/repository"
    "context"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentService struct {
    repo      *repository.StudentRepository
    achPGRepo *repository.AchievementPGRepository
    achMongo  *repository.AchievementMongoRepository
}

// 🔥 CONSTRUCTOR 
func NewStudentService(
    repo *repository.StudentRepository,
    achPG *repository.AchievementPGRepository,
    achMongo *repository.AchievementMongoRepository,
) *StudentService {
    return &StudentService{
        repo:      repo,
        achPGRepo: achPG,
        achMongo:  achMongo,
    }
}

// GET /students
// GetAll godoc
// @Summary Get all students
// @Description Mengambil daftar seluruh mahasiswa
// @Tags Student
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string][]model.StudentDetail
// @Failure 500 {object} map[string]any
// @Router /students [get]
func (s *StudentService) GetAll(c *fiber.Ctx) error {
    data, err := s.repo.GetAll(context.Background())
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"data": data})
}

// GET /students/:id
// GetByID godoc
// @Summary Get student by ID
// @Description Mengambil detail mahasiswa berdasarkan ID
// @Tags Student
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]model.StudentDetail
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /students/{id} [get]
func (s *StudentService) GetByID(c *fiber.Ctx) error {
    id := c.Params("id")
    data, err := s.repo.GetByID(context.Background(), id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Student tidak ditemukan"})
    }
    return c.JSON(fiber.Map{"data": data})
}

// POST /students
// Create godoc
// @Summary Create new student
// @Description Membuat data mahasiswa baru
// @Tags Student
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateStudentRequest true "Create student payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /students [post]
func (s *StudentService) Create(c *fiber.Ctx) error {
    var req model.CreateStudentRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    if req.UserID == "" || req.StudentID == "" {
        return c.Status(400).JSON(fiber.Map{"error": "user_id dan student_id wajib diisi"})
    }

    err := s.repo.Create(context.Background(), &model.Student{
        UserID:       req.UserID,
        StudentID:    req.StudentID,
        ProgramStudy: req.ProgramStudy,
        AcademicYear: req.AcademicYear,
        AdvisorID:    req.AdvisorID,
    })
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(201).JSON(fiber.Map{"message": "Student berhasil dibuat"})
}

// PUT /students/:id/advisor
// UpdateAdvisor godoc
// @Summary Update student advisor
// @Description Memperbarui dosen wali mahasiswa
// @Tags Student
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body object true "Advisor payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /students/{id}/advisor [put]
func (s *StudentService) UpdateAdvisor(c *fiber.Ctx) error {
    studentID := c.Params("id")

    body := struct {
        AdvisorID *string `json:"advisor_id"`
    }{}

    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    err := s.repo.UpdateAdvisor(context.Background(), studentID, body.AdvisorID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Advisor berhasil diperbarui"})
}

// =====================================================
// GET /students/:id/achievements
// =====================================================
// GetAchievements godoc
// @Summary Get student achievements
// @Description Mengambil daftar prestasi mahasiswa (role-based access)
// @Tags Student
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {array} model.StudentAchievement
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /students/{id}/achievements [get]
func (s *StudentService) GetAchievements(c *fiber.Ctx) error {
    ctx := context.Background()
    targetStudentID := c.Params("id")

    claims := c.Locals("user").(jwt.MapClaims)
    role := claims["role_name"].(string)
    userID := claims["user_id"].(string)

    // Mahasiswa hanya bisa lihat dirinya sendiri
    if role == "Mahasiswa" {
        sid, err := s.repo.GetStudentIDByUserID(ctx, userID)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "student not found"})
        }

        if sid != targetStudentID {
            return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
        }
    }

    // Dosen Wali hanya bisa lihat mahasiswa bimbingannya
    if role == "Dosen Wali" {
        lecturerID, err := s.repo.GetLecturerIDByUserID(ctx, userID)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "lecturer not found"})
        }

        ok, err := s.repo.IsStudentUnderAdvisor(ctx, targetStudentID, lecturerID)
        if err != nil || !ok {
            return c.Status(403).JSON(fiber.Map{"error": "Not your advisee"})
        }
    }

    // Ambil reference PostgreSQL
    refs, err := s.achPGRepo.GetByStudentID(ctx, targetStudentID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    results := []model.StudentAchievement{}

    // Ambil detail dari Mongo
    for _, ref := range refs {
        oid, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
        if err != nil {
            continue
        }

        detail, err := s.achMongo.FindById(ctx, oid)
        if err != nil {
            continue
        }

        results = append(results, model.StudentAchievement{
            ID:     ref.ID,
            Status: ref.Status,
            Detail: detail,
        })
    }

    return c.JSON(results)
}
