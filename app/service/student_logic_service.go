package service

import (
	"context"
	"errors"

	"backenduas/app/model"
	"backenduas/app/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentLogicService struct {
	stRepo repository.IStudentRepository
	pgRepo repository.IAchievementPGRepository
	mgRepo repository.IAchievementMongoRepository
}

func NewStudentLogicService(
	st repository.IStudentRepository,
	pg repository.IAchievementPGRepository,
	mg repository.IAchievementMongoRepository,
) *StudentLogicService {
	return &StudentLogicService{
		stRepo: st,
		pgRepo: pg,
		mgRepo: mg,
	}
}

//
// =====================================================
// LOGIC: GET ALL STUDENTS
// =====================================================
//
func (s *StudentLogicService) GetAllStudentsLogic() ([]model.StudentDetail, error) {
	ctx := context.Background()
	return s.stRepo.GetAll(ctx)
}

//
// =====================================================
// LOGIC: GET STUDENT BY ID
// =====================================================
//
func (s *StudentLogicService) GetStudentByIDLogic(id string) (*model.StudentDetail, error) {
	ctx := context.Background()

	data, err := s.stRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("student tidak ditemukan")
	}

	return data, nil
}

//
// =====================================================
// LOGIC: CREATE STUDENT (ADMIN)
// =====================================================
//
func (s *StudentLogicService) CreateStudentLogic(req model.CreateStudentRequest) error {
	ctx := context.Background()

	if req.UserID == "" || req.StudentID == "" {
		return errors.New("user_id dan student_id wajib diisi")
	}

	err := s.stRepo.Create(ctx, &model.Student{
		UserID:       req.UserID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    req.AdvisorID,
	})

	return err
}

//
// =====================================================
// LOGIC: UPDATE ADVISOR
// =====================================================
//
func (s *StudentLogicService) UpdateAdvisorLogic(
	studentID string,
	advisorID *string,
) error {

	ctx := context.Background()

	err := s.stRepo.UpdateAdvisor(ctx, studentID, advisorID)
	if err != nil {
		return errors.New("student tidak ditemukan")
	}

	return nil
}

//
// =====================================================
// LOGIC: GET STUDENT ACHIEVEMENTS
// =====================================================
//
func (s *StudentLogicService) GetAchievementsLogic(
	role string,
	userID string,
	targetStudentID string,
) ([]model.StudentAchievement, error) {

	ctx := context.Background()

	// ===============================
	// ROLE: MAHASISWA
	// ===============================
	if role == "Mahasiswa" {
		sid, err := s.stRepo.GetStudentIDByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("student not found")
		}

		if sid != targetStudentID {
			return nil, errors.New("forbidden")
		}
	}

	// ===============================
	// ROLE: DOSEN WALI
	// ===============================
	if role == "Dosen Wali" {
		lectID, err := s.stRepo.GetLecturerIDByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("lecturer not found")
		}

		ok, _ := s.stRepo.IsStudentUnderAdvisor(ctx, targetStudentID, lectID)
		if !ok {
			return nil, errors.New("not your advisee")
		}
	}

	// ===============================
	// AMBIL PRESTASI
	// ===============================
	refs, err := s.pgRepo.GetByStudentID(ctx, targetStudentID)
	if err != nil {
		return nil, err
	}

	result := []model.StudentAchievement{}

	for _, ref := range refs {
		oid, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		if err != nil {
			continue
		}

		detail, err := s.mgRepo.FindById(ctx, oid)
		if err != nil {
			continue
		}

		result = append(result, model.StudentAchievement{
			ID:     ref.ID,
			Status: ref.Status,
			Detail: detail,
		})
	}

	return result, nil
}
