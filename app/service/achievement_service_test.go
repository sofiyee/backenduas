package service

import (
	"testing"

	"backenduas/app/model"
	"backenduas/app/repository"
	
)

// ================= CREATE =================
func TestCreateAchievementLogic(t *testing.T) {
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()
	st := repository.NewMockStudentRepository()

	st.UserToStudent["u1"] = "s1"

	svc := NewAchievementLogicService(pg, mg, st)

	err := svc.Create("Mahasiswa", "u1", model.AchievementCreateRequest{
		Title:           "Juara 1",
		AchievementType: "Lomba",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateAchievementLogic_NotMahasiswa(t *testing.T) {
	svc := NewAchievementLogicService(nil, nil, nil)

	err := svc.Create("Admin", "u1", model.AchievementCreateRequest{})
	if err == nil {
		t.Fatalf("expected error")
	}
}

// ================= SUBMIT =================
func TestSubmitAchievementLogic(t *testing.T) {
	pg := repository.NewMockAchievementPGRepository()
	refID := pg.SeedDraft("s1")

	svc := NewAchievementLogicService(pg, nil, nil)

	err := svc.Submit(refID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================= DELETE =================
func TestDeleteAchievementLogic(t *testing.T) {
	pg := repository.NewMockAchievementPGRepository()
	refID := pg.SeedDraft("s1")

	svc := NewAchievementLogicService(pg, nil, nil)

	err := svc.Delete(refID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================= VERIFY =================
func TestVerifyAchievementLogic(t *testing.T) {
	pg := repository.NewMockAchievementPGRepository()
	st := repository.NewMockStudentRepository()

	refID := pg.SeedSubmitted("s1")
	st.UserToLecturer["uLect"] = "lect1"
	st.AdvisorMap["s1"] = "lect1"

	svc := NewAchievementLogicService(pg, nil, st)

	err := svc.Verify(refID, "uLect")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================= REJECT =================
func TestRejectAchievementLogic(t *testing.T) {
	pg := repository.NewMockAchievementPGRepository()
	st := repository.NewMockStudentRepository()

	refID := pg.SeedSubmitted("s1")
	st.UserToLecturer["uLect"] = "lect1"
	st.AdvisorMap["s1"] = "lect1"

	svc := NewAchievementLogicService(pg, nil, st)

	err := svc.Reject(refID, "uLect", "Kurang bukti")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================= UPDATE =================
func TestUpdateAchievementLogic(t *testing.T) {
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()
	st := repository.NewMockStudentRepository()

	// mapping user -> student
	st.UserToStudent["u1"] = "s1"

	// seed PG + ambil mongoID
	refID, oid := pg.SeedWithMongo("s1")

	// 🔥 WAJIB: seed Mongo data
	mg.Seed(oid)

	svc := NewAchievementLogicService(pg, mg, st)

	err := svc.Update(refID, "u1", model.AchievementUpdateInput{
		Title: "Updated Title",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}


// ================= HISTORY =================
func TestHistoryAchievementLogic(t *testing.T) {
	pg := repository.NewMockAchievementPGRepository()
	mg := repository.NewMockAchievementMongoRepository()

	refID, oid := pg.SeedWithMongo("s1")
	mg.Seed(oid)

	svc := NewAchievementLogicService(pg, mg, nil)

	history, err := svc.History(refID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(history) != 2 {
		t.Fatalf("expected 2 history items, got %d", len(history))
	}
}
