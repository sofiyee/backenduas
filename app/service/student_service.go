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

// 🔥 CONSTRUCTOR BARU (HARUS 3 PARAMETER)
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
