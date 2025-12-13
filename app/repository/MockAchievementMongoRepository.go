package repository

import (
	"context"
	"errors"
	"time"

	"backenduas/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAchievementMongoRepository struct {
	Data map[string]model.AchievementMongo // key = ObjectID.Hex()
}

func NewMockAchievementMongoRepository() *MockAchievementMongoRepository {
	return &MockAchievementMongoRepository{
		Data: make(map[string]model.AchievementMongo),
	}
}

// ===============================
// CREATE
// ===============================
func (m *MockAchievementMongoRepository) Create(
	ctx context.Context,
	ach model.AchievementMongo,
) (primitive.ObjectID, error) {

	oid := primitive.NewObjectID()
	ach.ID = oid
	ach.CreatedAt = time.Now().Unix()
	ach.UpdatedAt = ach.CreatedAt
	ach.Status = "draft"

	m.Data[oid.Hex()] = ach
	return oid, nil
}

// ===============================
// UPDATE
// ===============================
func (m *MockAchievementMongoRepository) Update(
	ctx context.Context,
	id primitive.ObjectID,
	update bson.M,
) error {

	ach, ok := m.Data[id.Hex()]
	if !ok {
		return errors.New("mongo data not found")
	}

	// manual apply update
	if v, ok := update["status"]; ok {
		ach.Status = v.(string)
	}

	ach.UpdatedAt = time.Now().Unix()
	m.Data[id.Hex()] = ach
	return nil
}

// ===============================
// SOFT DELETE
// ===============================
func (m *MockAchievementMongoRepository) SoftDelete(
	ctx context.Context,
	id primitive.ObjectID,
) error {

	ach, ok := m.Data[id.Hex()]
	if !ok {
		return errors.New("mongo data not found")
	}

	ach.Status = "deleted"
	ach.UpdatedAt = time.Now().Unix()
	m.Data[id.Hex()] = ach
	return nil
}

// ===============================
// FIND BY ID
// ===============================
func (m *MockAchievementMongoRepository) FindById(
	ctx context.Context,
	id primitive.ObjectID,
) (model.AchievementMongo, error) {

	if ach, ok := m.Data[id.Hex()]; ok {
		return ach, nil
	}
	return model.AchievementMongo{}, errors.New("mongo data not found")
}

// ===============================
// FIND MANY BY IDS
// ===============================
func (m *MockAchievementMongoRepository) FindManyByIDs(
	ctx context.Context,
	ids []primitive.ObjectID,
) (map[string]model.AchievementMongo, error) {

	result := make(map[string]model.AchievementMongo)

	for _, id := range ids {
		if ach, ok := m.Data[id.Hex()]; ok {
			result[id.Hex()] = ach
		}
	}

	return result, nil
}

// ===============================
// ADD ATTACHMENT
// ===============================
func (m *MockAchievementMongoRepository) AddAttachment(
	ctx context.Context,
	id primitive.ObjectID,
	file model.AttachmentFile,
) error {

	ach, ok := m.Data[id.Hex()]
	if !ok {
		return errors.New("mongo data not found")
	}

	// pastikan attachments array
	if ach.Attachments == nil {
		ach.Attachments = []model.AttachmentFile{}
	}

	ach.Attachments = append(ach.Attachments, file)
	ach.UpdatedAt = time.Now().Unix()

	m.Data[id.Hex()] = ach
	return nil
}

func (m *MockAchievementMongoRepository) Seed(id primitive.ObjectID) {
	m.Data[id.Hex()] = model.AchievementMongo{
		Title:     "Dummy",
		CreatedAt: 100,
		UpdatedAt: 200,
	}
}

