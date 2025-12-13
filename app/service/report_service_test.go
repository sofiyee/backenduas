package service

import (
	"testing"
	"time"

	"backenduas/app/model"
	"backenduas/app/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGlobalStatistics_Admin(t *testing.T) {
	st := repository.NewMockStudentRepository()
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()

	// student
	st.Students["s1"] = model.StudentDetail{
		ID:           "s1",
		FullName:     "Sofie",
		ProgramStudy: "TI",
		AcademicYear: "2023",
	}

	oid := primitive.NewObjectID()

	pg.Data["a1"] = model.AchievementReference{
		ID:                 "a1",
		StudentID:          "s1",
		MongoAchievementID: oid.Hex(),
	}

	mg.Data[oid.Hex()] = model.AchievementMongo{
		AchievementType: "Lomba",
		CreatedAt:       time.Now().Unix(),
		Tags:            []string{"Nasional"},
	}

	svc := NewReportLogicService(pg, mg, st)

	stats, err := svc.GlobalStatisticsLogic("Admin", "admin-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.PerType["Lomba"] != 1 {
		t.Fatalf("expected 1 lomba")
	}

	if len(stats.TopStudents) != 1 {
		t.Fatalf("expected 1 top student")
	}
}

func TestStudentStatisticsLogic_Mahasiswa_Forbidden(t *testing.T) {
	// mock repo
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()
	st := repository.NewMockStudentRepository()

	// mapping user -> student
	st.UserToStudent["u1"] = "S001"

	svc := NewReportLogicService(pg, mg, st)

	_, err := svc.StudentStatisticsLogic(
		"Mahasiswa",
		"u1",
		"S002", // student lain
	)

	if err == nil {
		t.Fatalf("expected forbidden error")
	}

	if err.Error() != "forbidden" {
		t.Fatalf("unexpected error message: %v", err.Error())
	}
}

func TestGlobalStatisticsLogic_DosenWali(t *testing.T) {
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()
	st := repository.NewMockStudentRepository()

	// lecturer mapping
	st.UserToLecturer["uLect"] = "L001"
	st.AdvisorMap["S001"] = "L001"

	oid := primitive.NewObjectID()

	// postgres ref
	pg.Data["a1"] = model.AchievementReference{
		ID:                 "a1",
		StudentID:          "S001",
		MongoAchievementID: oid.Hex(),
		Status:             "verified",
	}

	// mongo detail
	mg.Data[oid.Hex()] = model.AchievementMongo{
		AchievementType: "Academic",
		Tags:            []string{"OSN"},
		CreatedAt:       time.Now().Unix(),
	}

	svc := NewReportLogicService(pg, mg, st)

	stats, err := svc.GlobalStatisticsLogic("Dosen Wali", "uLect")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.PerType["Academic"] != 1 {
		t.Fatalf("expected 1 academic achievement")
	}
}

