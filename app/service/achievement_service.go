package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"backenduas/app/model"
	"backenduas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct {
	pgRepo      *repository.AchievementPGRepository
	mongoRepo   *repository.AchievementMongoRepository
	studentRepo *repository.StudentRepository
}

func NewAchievementService(
	pg *repository.AchievementPGRepository,
	mg *repository.AchievementMongoRepository,
	st *repository.StudentRepository,
) *AchievementService {
	return &AchievementService{pgRepo: pg, mongoRepo: mg, studentRepo: st}
}

// =====================================================
//
//	FR-003 — CREATE PRESTASI (Mahasiswa)
//
// =====================================================
func (s *AchievementService) Create(c *fiber.Ctx) error {

	claims := c.Locals("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)
	ctx := context.Background()

	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Student record not found for this user"})
	}

	var req model.AchievementCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if req.Title == "" || req.AchievementType == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title & achievementType required"})
	}

	ach := model.AchievementMongo{
		StudentID:       studentID,
		Title:           req.Title,
		Description:     req.Description,
		AchievementType: req.AchievementType,
		Details:         req.Details,
		Tags:            req.Tags,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
		Status:          "draft",
	}

	mongoID, err := s.mongoRepo.Create(ctx, ach)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	ref := model.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          studentID,
		MongoAchievementID: mongoID.Hex(),
		Status:             "draft",
	}

	if err := s.pgRepo.CreateReference(ref); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":     "achievement created",
		"referenceId": ref.ID,
	})
}

// =====================================================
//
//	FR-004 — SUBMIT PRESTASI
//
// =====================================================
func (s *AchievementService) Submit(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement not found"})
	}

	if err := s.pgRepo.UpdateStatus(id, "submitted"); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	s.mongoRepo.Update(ctx, mongoID, bson.M{"status": "submitted"})

	return c.JSON(fiber.Map{"message": "achievement submitted"})
}

// =====================================================
//
//	FR-005 — DELETE PRESTASI DRAFT
//
// =====================================================
func (s *AchievementService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement not found"})
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "Only draft achievements can be deleted"})
	}

	mongoID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid MongoID"})
	}

	if err := s.mongoRepo.SoftDelete(ctx, mongoID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete MongoDB data"})
	}

	if err := s.pgRepo.UpdateStatus(id, "deleted"); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update PostgreSQL"})
	}

	return c.JSON(fiber.Map{"message": "Achievement deleted successfully"})
}

// =====================================================
//
//	FR-006 + FR-010 — LIST PRESTASI (Role-based List)
//
// =====================================================
func (s *AchievementService) GetAll(c *fiber.Ctx) error {

	claims := c.Locals("user").(jwt.MapClaims)
	role := claims["role_name"].(string)
	userID := claims["user_id"].(string)
	ctx := context.Background()

	switch role {

	// =====================================================
	//  MAHASISWA → hanya prestasi miliknya (tanpa join Mongo)
	// =====================================================
	case "Mahasiswa":
		studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "student not found"})
		}

		data, err := s.pgRepo.GetByStudentID(ctx, studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(data)

	// =====================================================
	//  DOSEN WALI → ambil semua prestasi seluruh mahasiswa bimbingan
	// =====================================================
	case "Dosen Wali":

		advisorID, err := s.studentRepo.GetLecturerIDByUserID(ctx, userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "lecturer not found"})
		}

		// daftar mahasiswa bimbingan
		studentIDs, err := s.studentRepo.GetStudentsByAdvisor(ctx, advisorID)
		if err != nil || len(studentIDs) == 0 {
			return c.JSON([]any{})
		}

		refs, err := s.pgRepo.GetByStudentIDs(ctx, studentIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// kumpulkan ObjectID Mongo
		mongoIDs := []primitive.ObjectID{}
		for _, r := range refs {
			oid, err := primitive.ObjectIDFromHex(r.MongoAchievementID)
			if err == nil {
				mongoIDs = append(mongoIDs, oid)
			}
		}

		// ambil detail mongo sekaligus
		mongoMap, err := s.mongoRepo.FindManyByIDs(ctx, mongoIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// gabungkan PG + Mongo
		response := []fiber.Map{}
		for _, ref := range refs {
			response = append(response, fiber.Map{
				"id":         ref.ID,
				"student_id": ref.StudentID,
				"status":     ref.Status,
				"detail":     mongoMap[ref.MongoAchievementID],
			})
		}

		return c.JSON(response)

	// =====================================================
	//  ADMIN → FR-010 (list semua prestasi + join Mongo)
	// =====================================================
	case "Admin":

		// 1. Ambil semua achievement references
		refs, err := s.pgRepo.GetAll(ctx)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if len(refs) == 0 {
			return c.JSON([]any{})
		}

		// 2. Convert ke ObjectID list
		mongoIDs := []primitive.ObjectID{}
		for _, r := range refs {
			oid, err := primitive.ObjectIDFromHex(r.MongoAchievementID)
			if err == nil {
				mongoIDs = append(mongoIDs, oid)
			}
		}

		// 3. Fetch all Mongo detail
		mongoMap, err := s.mongoRepo.FindManyByIDs(ctx, mongoIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// 4. Gabungkan PG + Mongo
		response := []fiber.Map{}
		for _, ref := range refs {
			response = append(response, fiber.Map{
				"id":         ref.ID,
				"student_id": ref.StudentID,
				"status":     ref.Status,
				"detail":     mongoMap[ref.MongoAchievementID],
			})
		}

		return c.JSON(response)

	// =====================================================
	//  ROLE INVALID
	// =====================================================
	default:
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
	}
}

// =====================================================
//
//	FR-007 — VERIFY PRESTASI (Dosen Wali)
//
// =====================================================
func (s *AchievementService) Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	// 1. Ambil reference dari PostgreSQL
	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement not found"})
	}

	// 2. Precondition: harus submitted
	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Only submitted achievements can be verified",
		})
	}

	// 3. Ambil user_id dari token
	claims := c.Locals("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	// 4. Convert user_id → lecturer_id (FK yang benar)
	lecturerID, err := s.studentRepo.GetLecturerIDByUserID(ctx, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "lecturer not found"})
	}

	// 5. Update PostgreSQL
	if err := s.pgRepo.Verify(id, lecturerID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 6. Update MongoDB
	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	s.mongoRepo.Update(ctx, mongoID, bson.M{
		"status":    "verified",
		"updatedAt": time.Now().Unix(),
	})

	// 7. Return success
	return c.JSON(fiber.Map{
		"message":     "achievement verified",
		"id":          id,
		"verified_by": lecturerID,
		"status":      "verified",
		"notification": model.Notification{
			Title:   "Prestasi Diverifikasi",
			Message: "Prestasi kamu telah diverifikasi oleh dosen wali.",
			SentAt:  time.Now().Unix(),
		},
	})
}

// =====================================================
//
//	FR-008 — REJECT PRESTASI (Dosen Wali)
//
// =====================================================
func (s *AchievementService) Reject(c *fiber.Ctx) error {

	id := c.Params("id")
	ctx := context.Background()

	// 1. Ambil reference
	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement not found"})
	}

	// 2. Precondition: harus submitted
	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Only submitted achievements can be rejected",
		})
	}

	// 3. Ambil dosen wali (verifier)
	claims := c.Locals("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	advisorID, err := s.studentRepo.GetLecturerIDByUserID(ctx, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "lecturer not found"})
	}

	// pastikan student ini adalah bimbingannya
	isBimbingan, _ := s.studentRepo.IsStudentUnderAdvisor(ctx, ref.StudentID, advisorID)
	if !isBimbingan {
		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden: this student is not your advisee",
		})
	}

	// 4. Input rejection note
	var body struct {
		Note string `json:"note"`
	}
	c.BodyParser(&body)

	if body.Note == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Rejection note is required"})
	}

	// 5. Update status PostgreSQL
	if err := s.pgRepo.Reject(id, body.Note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 6. Update MongoDB
	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	s.mongoRepo.Update(ctx, mongoID, bson.M{
		"status": "rejected",
	})

	// 7. (Optional) Simulasi Notification
	//
	_ = s.createNotificationDummy(ref.StudentID, "Prestasi kamu ditolak oleh dosen wali")

	// 8. Return response
	return c.JSON(fiber.Map{
		"message":        "achievement rejected",
		"id":             id,
		"status":         "rejected",
		"rejection_note": body.Note,
		"notification": model.Notification{
			Title:   "Prestasi Ditolak",
			Message: "Prestasi kamu ditolak. Catatan: " + body.Note,
			SentAt:  time.Now().Unix(),
		},
	})

}

// =====================================================
//
//	Dummy Notification — tidak disimpan, hanya simulasi
//
// =====================================================
func (s *AchievementService) createNotificationDummy(studentID string, msg string) model.Notification {
	return model.Notification{
		Title:   "Notifikasi Prestasi",
		Message: msg,
		SentAt:  time.Now().Unix(),
	}
}

// =====================================================
//
//	FR-??? — DETAIL PRESTASI (Admin / Dosen Wali / Mahasiswa)
//
// =====================================================
func (s *AchievementService) GetByID(c *fiber.Ctx) error {

	id := c.Params("id")
	ctx := context.Background()

	// 1. Ambil reference dari PostgreSQL
	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement not found"})
	}

	// 2. Ambil claims user
	claims := c.Locals("user").(jwt.MapClaims)
	role := claims["role_name"].(string)
	userID := claims["user_id"].(string)

	// ----------------------------
	// 3. VALIDASI AKSES BERDASARKAN ROLE
	// ----------------------------

	switch role {

	// ======== MAHASISWA ========
	case "Mahasiswa":
		studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "student not found"})
		}

		if ref.StudentID != studentID {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden: not your achievement"})
		}

	// ======== DOSEN WALI ========
	case "Dosen Wali":
		advisorID, err := s.studentRepo.GetLecturerIDByUserID(ctx, userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "lecturer not found"})
		}

		// cek apakah student ini benar dibimbing dosen ini
		isBimbingan, err := s.studentRepo.IsStudentUnderAdvisor(ctx, ref.StudentID, advisorID)
		if err != nil || !isBimbingan {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden: not your student"})
		}

	// ======== ADMIN ========
	case "Admin":
		// bebas
	default:
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
	}

	// ----------------------------
	// 4. Ambil detail Achievement dari MongoDB
	// ----------------------------

	mongoID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid mongo id"})
	}

	mongoData, err := s.mongoRepo.FindById(ctx, mongoID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement detail not found in MongoDB"})
	}

	// ----------------------------
	// 5. Format response
	// ----------------------------

	return c.JSON(fiber.Map{
		"id":             ref.ID,
		"student_id":     ref.StudentID,
		"status":         ref.Status,
		"submitted_at":   ref.SubmittedAt,
		"verified_at":    ref.VerifiedAt,
		"verified_by":    ref.VerifiedBy,
		"rejection_note": ref.RejectionNote,
		"created_at":     ref.CreatedAt,
		"updated_at":     ref.UpdatedAt,
		"detail":         mongoData,
	})
}

func (s *AchievementService) Update(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	// Ambil claim mahasiswa
	claims := c.Locals("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	// Ambil student_id mahasiswa login
	sid, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "student not found"})
	}

	// Ambil reference di Postgre
	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
	}

	// Pastikan prestasi miliknya
	if ref.StudentID != sid {
		return c.Status(403).JSON(fiber.Map{"error": "cannot edit other student's achievement"})
	}

	// Pastikan status = DRAFT
	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only DRAFT achievements can be updated"})
	}

	// Parse body

	var body model.AchievementUpdateInput
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Build bson update sesuai field Mongo
	update := bson.M{
		"title":           body.Title,
		"description":     body.Description,
		"achievementType": body.AchievementType,
		"details":         body.Details, // <- important
		"tags":            body.Tags,
	}

	// Update Mongo
	oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err := s.mongoRepo.Update(ctx, oid, update); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "achievement updated"})
}

func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
    ctx := context.Background()
    id := c.Params("id")

    claims := c.Locals("user").(jwt.MapClaims)
    userID := claims["user_id"].(string)

    sid, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
    if err != nil {
        return c.Status(403).JSON(fiber.Map{"error": "student not found"})
    }

    ref, err := s.pgRepo.GetByID(ctx, id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
    }

    if ref.StudentID != sid {
        return c.Status(403).JSON(fiber.Map{"error": "cannot upload to other's achievement"})
    }

    if ref.Status != "draft" {
        return c.Status(400).JSON(fiber.Map{"error": "cannot upload attachment when not in DRAFT status"})
    }

    // Ambil file
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "file is required"})
    }

    // Simpan file lokal
    folder := fmt.Sprintf("./uploads/achievements/%s/", id)
    os.MkdirAll(folder, os.ModePerm)

    filePath := folder + file.Filename
    if err := c.SaveFile(file, filePath); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // Metadata file
    attachment := model.AttachmentFile{
        FileName:   file.Filename,
        FileURL:    filePath,
        FileType:   file.Header.Get("Content-Type"),
        UploadedAt: time.Now().Unix(),
    }

    // 🔥 INI YANG PALING PENTING: gunakan repo.AddAttachment()
    oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    if err := s.mongoRepo.AddAttachment(ctx, oid, attachment); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "message":   "file uploaded",
        "file_name": attachment.FileName,
        "file_url":  attachment.FileURL,
    })
}

func (s *AchievementService) History(c *fiber.Ctx) error {
    ctx := context.Background()
    id := c.Params("id")

    // Ambil user claims
    claims := c.Locals("user").(jwt.MapClaims)
    role := claims["role_name"].(string)
    userID := claims["user_id"].(string)

    // --- Ambil Postgre reference ---
    ref, err := s.pgRepo.GetByID(ctx, id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
    }

    // --- Authorization sesuai role ---
    switch role {
    case "Mahasiswa":
        sid, _ := s.studentRepo.GetStudentIDByUserID(ctx, userID)
        if ref.StudentID != sid {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }

    case "Dosen Wali":
        lecID, _ := s.studentRepo.GetLecturerIDByUserID(ctx, userID)
        ok, _ := s.studentRepo.IsStudentUnderAdvisor(ctx, ref.StudentID, lecID)
        if !ok {
            return c.Status(403).JSON(fiber.Map{"error": "not your advisee"})
        }

    case "Admin":
        // admin bebas akses
    }

    // --- Ambil detail dari Mongo ---
    oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
    detail, err := s.mongoRepo.FindById(ctx, oid)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "mongo data not found"})
    }

    // --- Bentuk "history" sederhana sesuai SRS ---
    history := []fiber.Map{
        {
            "status":    "created",
            "timestamp": detail.CreatedAt,
        },
        {
            "status":    ref.Status,
            "timestamp": detail.UpdatedAt,
        },
    }

    return c.JSON(fiber.Map{
        "achievement_id": id,
        "history":        history,
    })
}


