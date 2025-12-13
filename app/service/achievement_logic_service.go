package service

import (
	"context"
	"errors"
	"time"

	"backenduas/app/model"
	"backenduas/app/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementLogicService struct {
	pgRepo      repository.IAchievementPGRepository
	mongoRepo   repository.IAchievementMongoRepository
	studentRepo repository.IStudentRepository
}

func NewAchievementLogicService(
	pg repository.IAchievementPGRepository,
	mg repository.IAchievementMongoRepository,
	st repository.IStudentRepository,
) *AchievementLogicService {
	return &AchievementLogicService{pgRepo: pg, mongoRepo: mg, studentRepo: st}
}

// ================= CREATE =================
func (s *AchievementLogicService) Create(
	role string,
	userID string,
	req model.AchievementCreateRequest,
) error {

	if role != "Mahasiswa" {
		return errors.New("akses ditolak")
	}

	if req.Title == "" || req.AchievementType == "" {
		return errors.New("data tidak valid")
	}

	ctx := context.Background()
	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return errors.New("student tidak ditemukan")
	}

	ach := model.AchievementMongo{
		StudentID:       studentID,
		Title:           req.Title,
		AchievementType: req.AchievementType,
		Description:     req.Description,
		Status:          "draft",
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	}

	oid, _ := s.mongoRepo.Create(ctx, ach)

	ref := model.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          studentID,
		MongoAchievementID: oid.Hex(),
		Status:             "draft",
	}

	return s.pgRepo.CreateReference(ref)
}

// ================= SUBMIT =================
func (s *AchievementLogicService) Submit(id string) error {
	ctx := context.Background()

	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("achievement tidak ditemukan")
	}

	if ref.Status != "draft" {
		return errors.New("hanya draft yang bisa disubmit")
	}

	return s.pgRepo.UpdateStatus(id, "submitted")
}

// ================= DELETE =================
func (s *AchievementLogicService) Delete(id string) error {
	ctx := context.Background()

	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("achievement tidak ditemukan")
	}

	if ref.Status != "draft" {
		return errors.New("hanya draft yang bisa dihapus")
	}

	return s.pgRepo.UpdateStatus(id, "deleted")
}

// ================= VERIFY =================
func (s *AchievementLogicService) Verify(id string, lecturerUserID string) error {
	ctx := context.Background()

	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("achievement tidak ditemukan")
	}

	if ref.Status != "submitted" {
		return errors.New("achievement belum disubmit")
	}

	lecturerID, _ := s.studentRepo.GetLecturerIDByUserID(ctx, lecturerUserID)
	ok, _ := s.studentRepo.IsStudentUnderAdvisor(ctx, ref.StudentID, lecturerID)
	if !ok {
		return errors.New("bukan mahasiswa bimbingan")
	}

	return s.pgRepo.Verify(id, lecturerID)
}

// ================= REJECT =================
func (s *AchievementLogicService) Reject(id string, lecturerUserID string, note string) error {
	ctx := context.Background()

	if note == "" {
		return errors.New("catatan wajib diisi")
	}

	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("achievement tidak ditemukan")
	}

	lecturerID, _ := s.studentRepo.GetLecturerIDByUserID(ctx, lecturerUserID)
	ok, _ := s.studentRepo.IsStudentUnderAdvisor(ctx, ref.StudentID, lecturerID)
	if !ok {
		return errors.New("bukan mahasiswa bimbingan")
	}

	return s.pgRepo.Reject(id, note)
}

// ================= HISTORY =================
func (s *AchievementLogicService) History(id string) ([]map[string]any, error) {
	ctx := context.Background()

	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("achievement tidak ditemukan")
	}

	oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	detail, _ := s.mongoRepo.FindById(ctx, oid)

	history := []map[string]any{
		{"status": "created", "timestamp": detail.CreatedAt},
		{"status": ref.Status, "timestamp": detail.UpdatedAt},
	}

	return history, nil
}

// ================= UPDATE =================
func (s *AchievementLogicService) Update(id string, studentUserID string, req model.AchievementUpdateInput) error {
	ctx := context.Background()

	studentID, _ := s.studentRepo.GetStudentIDByUserID(ctx, studentUserID)
	ref, err := s.pgRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("achievement tidak ditemukan")
	}

	if ref.StudentID != studentID {
		return errors.New("bukan prestasi milik sendiri")
	}

	if ref.Status != "draft" {
		return errors.New("hanya draft yang bisa diubah")
	}

	oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	return s.mongoRepo.Update(ctx, oid, bson.M{
		"title": req.Title,
	})
}
