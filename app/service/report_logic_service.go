package service

import (
	"context"
	"errors"
	"time"

	"backenduas/app/model"
	"backenduas/app/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportLogicService struct {
	pgRepo      repository.IAchievementPGRepository
	mongoRepo   repository.IAchievementMongoRepository
	studentRepo repository.IStudentRepository
}

func NewReportLogicService(
	pg repository.IAchievementPGRepository,
	mg repository.IAchievementMongoRepository,
	st repository.IStudentRepository,
) *ReportLogicService {
	return &ReportLogicService{
		pgRepo:      pg,
		mongoRepo:   mg,
		studentRepo: st,
	}
}

func (s *ReportLogicService) GlobalStatisticsLogic(
	role string,
	userID string,
) (model.StudentAchievementStat, error) {

	ctx := context.Background()

	var studentIDs []string
	var err error

	switch role {
	case "Admin":
		studentIDs, err = s.studentRepo.GetAllStudentIDs(ctx)

	case "Dosen Wali":
		lecID, err1 := s.studentRepo.GetLecturerIDByUserID(ctx, userID)
		if err1 != nil {
			return model.StudentAchievementStat{}, errors.New("lecturer not found")
		}
		studentIDs, err = s.studentRepo.GetStudentsByAdvisor(ctx, lecID)

	case "Mahasiswa":
		sid, err1 := s.studentRepo.GetStudentIDByUserID(ctx, userID)
		if err1 != nil {
			return model.StudentAchievementStat{}, errors.New("student not found")
		}
		studentIDs = []string{sid}
	default:
		return model.StudentAchievementStat{}, errors.New("forbidden")
	}

	if err != nil {
		return model.StudentAchievementStat{}, err
	}

	refs, err := s.pgRepo.GetByStudentIDs(ctx, studentIDs)
	if err != nil {
		return model.StudentAchievementStat{}, err
	}

	stats := model.StudentAchievementStat{
		PerType:     map[string]int{},
		PerPeriod:   map[string]int{},
		Competition: map[string]int{},
		TopStudents: []model.TopStudentDetail{},
	}

	counter := map[string]int{}

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

		counter[r.StudentID]++
	}

	stMap, _ := s.studentRepo.GetStudentsByIDs(ctx, studentIDs)
	for sid, cnt := range counter {
		if st, ok := stMap[sid]; ok {
			stats.TopStudents = append(stats.TopStudents, model.TopStudentDetail{
				StudentID:    sid,
				FullName:     st.FullName,
				ProgramStudy: st.ProgramStudy,
				AcademicYear: st.AcademicYear,
				Count:        cnt,
			})
		}
	}

	return stats, nil
}

func (s *ReportLogicService) StudentStatisticsLogic(
	role string,
	userID string,
	studentID string,
) (*model.StudentAchievementStat, error) {

	ctx := context.Background()

	// ================= ROLE VALIDATION =================
	switch role {

	case "Mahasiswa":
		sid, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("student not found")
		}
		if sid != studentID {
			return nil, errors.New("forbidden")
		}

	case "Dosen Wali":
		lectID, err := s.studentRepo.GetLecturerIDByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("lecturer not found")
		}

		ok, _ := s.studentRepo.IsStudentUnderAdvisor(ctx, studentID, lectID)
		if !ok {
			return nil, errors.New("forbidden")
		}

	case "Admin":
		// bebas

	default:
		return nil, errors.New("forbidden")
	}

	// ================= DATA =================
	refs, _ := s.pgRepo.GetByStudentID(ctx, studentID)

	stats := &model.StudentAchievementStat{
		PerType:     map[string]int{},
		PerPeriod:   map[string]int{},
		Competition: map[string]int{},
		TopStudents: []model.TopStudentDetail{},
	}

	count := 0

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

		count++
	}

	studentMap, _ := s.studentRepo.GetStudentsByIDs(ctx, []string{studentID})
	st := studentMap[studentID]

	stats.TopStudents = append(stats.TopStudents, model.TopStudentDetail{
		StudentID:    studentID,
		FullName:     st.FullName,
		ProgramStudy: st.ProgramStudy,
		AcademicYear: st.AcademicYear,
		Count:        count,
	})

	return stats, nil
}
