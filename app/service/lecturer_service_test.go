package service

import (
	"testing"

	"backenduas/app/model"
	"backenduas/app/repository"
)

func TestGetAllLecturersLogic(t *testing.T) {
	repo := repository.NewMockLecturerRepository()
	repo.Lecturers["1"] = model.LecturerDetail{
		ID:         "1",
		LecturerID: "L001",
		Department: "Informatika",
	}

	svc := NewLecturerLogicService(repo)

	data, err := svc.GetAllLecturersLogic()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) != 1 {
		t.Fatalf("expected 1 lecturer, got %d", len(data))
	}
}

func TestGetAdviseesLogic_Success(t *testing.T) {
	repo := repository.NewMockLecturerRepository()
	repo.Advisees["lect1"] = []model.Advisee{
		{
			StudentID:    "S001",
			StudentName:  "Sofie",
			ProgramStudy: "Informatika",
			AcademicYear: "2023",
		},
	}

	svc := NewLecturerLogicService(repo)

	data, err := svc.GetAdviseesLogic("lect1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) != 1 {
		t.Fatalf("expected 1 advisee, got %d", len(data))
	}
}

func TestGetAdviseesLogic_InvalidID(t *testing.T) {
	repo := repository.NewMockLecturerRepository()
	svc := NewLecturerLogicService(repo)

	_, err := svc.GetAdviseesLogic("")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetAdviseesLogic_NotFound(t *testing.T) {
	repo := repository.NewMockLecturerRepository()
	svc := NewLecturerLogicService(repo)

	_, err := svc.GetAdviseesLogic("unknown")
	if err == nil {
		t.Fatalf("expected error")
	}
}
