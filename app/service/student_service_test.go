package service

import (
	"testing"

	"backenduas/app/model"
	"backenduas/app/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllStudentsLogic(t *testing.T) {
	st := repository.NewMockStudentRepository()
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()

	st.Students["s1"] = model.StudentDetail{ID: "s1", StudentID: "2021001"}
	st.Students["s2"] = model.StudentDetail{ID: "s2", StudentID: "2021002"}

	svc := NewStudentLogicService(st, pg, mg)

	data, err := svc.GetAllStudentsLogic()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) != 2 {
		t.Fatalf("expected 2 students, got %d", len(data))
	}
}

func TestGetStudentByIDLogic_Success(t *testing.T) {
	st := repository.NewMockStudentRepository()

	st.Students["s1"] = model.StudentDetail{
		ID:        "s1",
		StudentID: "2021001",
	}

	svc := NewStudentLogicService(st, nil, nil)

	data, err := svc.GetStudentByIDLogic("s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.StudentID != "2021001" {
		t.Fatalf("wrong student data")
	}
}

func TestGetStudentByIDLogic_NotFound(t *testing.T) {
	st := repository.NewMockStudentRepository()
	svc := NewStudentLogicService(st, nil, nil)

	_, err := svc.GetStudentByIDLogic("x")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCreateStudentLogic(t *testing.T) {
	st := repository.NewMockStudentRepository()
	svc := NewStudentLogicService(st, nil, nil)

	req := model.CreateStudentRequest{
		UserID:       "u1",
		StudentID:    "2021001",
		ProgramStudy: "Informatics",
		AcademicYear: "2021",
	}

	err := svc.CreateStudentLogic(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := st.Students["2021001"]; !ok {
		t.Fatalf("student not created")
	}
}

func TestUpdateAdvisorLogic(t *testing.T) {
	st := repository.NewMockStudentRepository()

	st.Students["s1"] = model.StudentDetail{
		ID: "s1",
	}

	advisorID := "lect1"
	svc := NewStudentLogicService(st, nil, nil)

	err := svc.UpdateAdvisorLogic("s1", &advisorID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if *st.Students["s1"].AdvisorID != "lect1" {
		t.Fatalf("advisor not updated")
	}
}


func TestGetAchievements_Mahasiswa_Self(t *testing.T) {
	st := repository.NewMockStudentRepository()
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()

	// mapping user -> student
	st.UserToStudent["u1"] = "s1"

	oid := primitive.NewObjectID()

	pg.Data["a1"] = model.AchievementReference{
		ID:                 "a1",
		StudentID:          "s1",
		MongoAchievementID: oid.Hex(),
		Status:             "draft",
	}

	mg.Data[oid.Hex()] = model.AchievementMongo{
		Title: "Juara 1",
	}

	svc := NewStudentLogicService(st, pg, mg)

	data, err := svc.GetAchievementsLogic("Mahasiswa", "u1", "s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) != 1 {
		t.Fatalf("expected 1 achievement, got %d", len(data))
	}
}

func TestGetAchievements_Mahasiswa_Other(t *testing.T) {
	st := repository.NewMockStudentRepository()
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()

	st.UserToStudent["u1"] = "s1"

	svc := NewStudentLogicService(st, pg, mg)

	_, err := svc.GetAchievementsLogic("Mahasiswa", "u1", "s2")
	if err == nil {
		t.Fatalf("expected forbidden error")
	}
}

func TestGetAchievements_DosenWali_Success(t *testing.T) {
	st := repository.NewMockStudentRepository()
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()

	st.UserToLecturer["uLect"] = "lect1"
	st.AdvisorMap["s1"] = "lect1"

	oid := primitive.NewObjectID()

	pg.Data["a1"] = model.AchievementReference{
		ID:                 "a1",
		StudentID:          "s1",
		MongoAchievementID: oid.Hex(),
		Status:             "verified",
	}

	mg.Data[oid.Hex()] = model.AchievementMongo{
		Title: "Juara 2",
	}

	svc := NewStudentLogicService(st, pg, mg)

	data, err := svc.GetAchievementsLogic("Dosen Wali", "uLect", "s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) != 1 {
		t.Fatalf("expected 1 achievement, got %d", len(data))
	}
}

func TestGetAchievements_DosenWali_Forbidden(t *testing.T) {
	st := repository.NewMockStudentRepository()
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()

	st.UserToLecturer["uLect"] = "lect1"
	st.AdvisorMap["s2"] = "lect1" 

	svc := NewStudentLogicService(st, pg, mg)

	_, err := svc.GetAchievementsLogic("Dosen Wali", "uLect", "s1")
	if err == nil {
		t.Fatalf("expected forbidden error")
	}
}
