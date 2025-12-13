package repository

import (
	"context"
	"errors"
	"time"

	"backenduas/app/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAchievementPGRepository struct {
	Data map[string]model.AchievementReference
}

func NewMockAchievementPGRepository() *MockAchievementPGRepository {
	return &MockAchievementPGRepository{
		Data: make(map[string]model.AchievementReference),
	}
}

// ===================================
// SEED UNTUK TEST
// ===================================
func (m *MockAchievementPGRepository) SeedDraft(id string) string {
	ref := model.AchievementReference{
		ID:        id,
		StudentID: "s1",
		Status:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.Data[id] = ref
	return id
}

func (m *MockAchievementPGRepository) SeedSubmitted(id string) string {
	now := time.Now()
	ref := model.AchievementReference{
		ID:          id,
		StudentID:   "s1",
		Status:      "submitted",
		SubmittedAt: &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	m.Data[id] = ref
	return id
}

// ===================================
// IMPLEMENT INTERFACE
// ===================================
func (m *MockAchievementPGRepository) CreateReference(ref model.AchievementReference) error {
	ref.CreatedAt = time.Now()
	ref.UpdatedAt = time.Now()
	m.Data[ref.ID] = ref
	return nil
}

func (m *MockAchievementPGRepository) UpdateStatus(id string, status string) error {
	ref, ok := m.Data[id]
	if !ok {
		return errors.New("achievement not found")
	}
	ref.Status = status
	ref.UpdatedAt = time.Now()
	m.Data[id] = ref
	return nil
}

func (m *MockAchievementPGRepository) Verify(id, verifier string) error {
	ref, ok := m.Data[id]
	if !ok {
		return errors.New("achievement not found")
	}
	now := time.Now()
	ref.Status = "verified"
	ref.VerifiedBy = &verifier
	ref.VerifiedAt = &now
	ref.UpdatedAt = now
	m.Data[id] = ref
	return nil
}

func (m *MockAchievementPGRepository) Reject(id, note string) error {
	ref, ok := m.Data[id]
	if !ok {
		return errors.New("achievement not found")
	}
	ref.Status = "rejected"
	ref.RejectionNote = &note
	ref.UpdatedAt = time.Now()
	m.Data[id] = ref
	return nil
}

func (m *MockAchievementPGRepository) GetByID(ctx context.Context, id string) (model.AchievementReference, error) {
	ref, ok := m.Data[id]
	if !ok {
		return model.AchievementReference{}, errors.New("achievement not found")
	}
	return ref, nil
}

func (m *MockAchievementPGRepository) GetByStudentID(ctx context.Context, studentID string) ([]model.AchievementReference, error) {
	var out []model.AchievementReference
	for _, ref := range m.Data {
		if ref.StudentID == studentID {
			out = append(out, ref)
		}
	}
	return out, nil
}

func (m *MockAchievementPGRepository) GetByStudentIDs(ctx context.Context, studentIDs []string) ([]model.AchievementReference, error) {
	set := map[string]bool{}
	for _, id := range studentIDs {
		set[id] = true
	}
	var out []model.AchievementReference
	for _, ref := range m.Data {
		if set[ref.StudentID] {
			out = append(out, ref)
		}
	}
	return out, nil
}

func (m *MockAchievementPGRepository) GetAll(ctx context.Context) ([]model.AchievementReference, error) {
	var out []model.AchievementReference
	for _, ref := range m.Data {
		out = append(out, ref)
	}
	return out, nil
}

func (m *MockAchievementPGRepository) SeedWithMongo(studentID string) (string, primitive.ObjectID) {
	oid := primitive.NewObjectID()

	ref := model.AchievementReference{
		ID:                 "ref-" + studentID,
		StudentID:          studentID,
		MongoAchievementID: oid.Hex(),
		Status:             "draft",
	}

	m.Data[ref.ID] = ref
	return ref.ID, oid
}
