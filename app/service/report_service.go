package service

import (
    "backenduas/app/model"
    "backenduas/app/repository"
    "context"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportService struct {
    pgRepo      *repository.AchievementPGRepository
    mongoRepo   *repository.AchievementMongoRepository
    studentRepo *repository.StudentRepository
}

func NewReportService(pg *repository.AchievementPGRepository,
    mg *repository.AchievementMongoRepository,
    st *repository.StudentRepository) *ReportService {

    return &ReportService{
        pgRepo:      pg,
        mongoRepo:   mg,
        studentRepo: st,
    }
}

//
// ======================================================
//  GLOBAL STATISTICS
// ======================================================
func (s *ReportService) GlobalStatistics(c *fiber.Ctx) error {
    ctx := context.Background()

    claims := c.Locals("user").(jwt.MapClaims)
    role := claims["role_name"].(string)
    userID := claims["user_id"].(string)

    var studentIDs []string
    var err error

    switch role {

    case "Admin":
        studentIDs, err = s.studentRepo.GetAllStudentIDs(ctx)

    case "Dosen Wali":
        lecID, err1 := s.studentRepo.GetLecturerIDByUserID(ctx, userID)
        if err1 != nil {
            return c.Status(404).JSON(fiber.Map{"error": "lecturer not found"})
        }
        studentIDs, err = s.studentRepo.GetStudentsByAdvisor(ctx, lecID)

    case "Mahasiswa":
        sid, err1 := s.studentRepo.GetStudentIDByUserID(ctx, userID)
        if err1 != nil {
            return c.Status(404).JSON(fiber.Map{"error": "student not found"})
        }
        studentIDs = []string{sid}
    }

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // -------- Ambil data prestasi --------
    refs, err := s.pgRepo.GetByStudentIDs(ctx, studentIDs)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    stats := model.StudentAchievementStat{
        PerType:     map[string]int{},
        PerPeriod:   map[string]int{},
        Competition: map[string]int{},
        TopStudents: []model.TopStudentDetail{},
    }

    topCounter := map[string]int{}

    // -------- Hitung statistik --------
    for _, r := range refs {
        oid, err := primitive.ObjectIDFromHex(r.MongoAchievementID)
        if err != nil {
            continue
        }

        detail, err := s.mongoRepo.FindById(ctx, oid)
        if err != nil {
            continue
        }

        stats.PerType[detail.AchievementType]++

        year := time.Unix(detail.CreatedAt, 0).Format("2006")
        stats.PerPeriod[year]++

        for _, tag := range detail.Tags {
            stats.Competition[tag]++
        }

        topCounter[r.StudentID]++
    }

    // -------- Ambil data detail student untuk TOP STUDENTS --------
    studentMap, err := s.studentRepo.GetStudentsByIDs(ctx, studentIDs)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    for sid, cnt := range topCounter {
        if st, ok := studentMap[sid]; ok {
            stats.TopStudents = append(stats.TopStudents, model.TopStudentDetail{
                StudentID:    sid,
                FullName:     st.FullName,
                ProgramStudy: st.ProgramStudy,
                AcademicYear: st.AcademicYear,
                Count:        cnt,
            })
        }
    }

    return c.JSON(stats)
}

//
// ======================================================
//  STUDENT STATISTICS
// ======================================================
func (s *ReportService) StudentStatistics(c *fiber.Ctx) error {
    ctx := context.Background()

    studentID := c.Params("id")

    // Ambil role + userID dari token
    claims := c.Locals("user").(jwt.MapClaims)
    role := claims["role_name"].(string)
    userID := claims["user_id"].(string)

    // ================================================
    // 🔐 ROLE VALIDATION (Mahasiswa / Dosen / Admin)
    // ================================================

    switch role {

    case "Mahasiswa":
        // mahasiswa hanya boleh akses dirinya sendiri
        sid, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "student not found"})
        }
        if sid != studentID {
            return c.Status(403).JSON(fiber.Map{"error": "Forbidden: You can only view your own statistics"})
        }

    case "Dosen Wali":
        // dosen wali hanya boleh akses mahasiswa bimbingannya
        lecID, err := s.studentRepo.GetLecturerIDByUserID(ctx, userID)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "lecturer not found"})
        }

        ok, _ := s.studentRepo.IsStudentUnderAdvisor(ctx, studentID, lecID)
        if !ok {
            return c.Status(403).JSON(fiber.Map{"error": "Forbidden: This student is not your advisee"})
        }

    case "Admin":
        // bebas
    default:
        return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
    }

    // ================================================
    // Jika lolos validasi → ambil statistik student
    // ================================================
    refs, err := s.pgRepo.GetByStudentID(ctx, studentID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    stats := model.StudentAchievementStat{
        PerType:     map[string]int{},
        PerPeriod:   map[string]int{},
        Competition: map[string]int{},
        TopStudents: []model.TopStudentDetail{},
    }

    totalCount := 0

    for _, r := range refs {
        oid, _ := primitive.ObjectIDFromHex(r.MongoAchievementID)
        detail, err := s.mongoRepo.FindById(ctx, oid)
        if err != nil {
            continue
        }

        stats.PerType[detail.AchievementType]++

        year := time.Unix(detail.CreatedAt, 0).Format("2006")
        stats.PerPeriod[year]++

        for _, tag := range detail.Tags {
            stats.Competition[tag]++
        }

        totalCount++
    }

    // Ambil detail student
    studentMap, _ := s.studentRepo.GetStudentsByIDs(ctx, []string{studentID})
    st := studentMap[studentID]

    stats.TopStudents = append(stats.TopStudents, model.TopStudentDetail{
        StudentID:    studentID,
        FullName:     st.FullName,
        ProgramStudy: st.ProgramStudy,
        AcademicYear: st.AcademicYear,
        Count:        totalCount,
    })

    return c.JSON(stats)
}

