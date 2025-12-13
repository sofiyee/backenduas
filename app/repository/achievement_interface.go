package repository

import (
	"context"
	"backenduas/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
=====================================================
POSTGRES ACHIEVEMENT (REFERENCE)
=====================================================
*/
type IAchievementPGRepository interface {
	CreateReference(ref model.AchievementReference) error
	UpdateStatus(id string, status string) error
	Verify(id, verifier string) error
	Reject(id, note string) error

	GetByID(ctx context.Context, id string) (model.AchievementReference, error)
	GetByStudentID(ctx context.Context, studentID string) ([]model.AchievementReference, error)
	GetByStudentIDs(ctx context.Context, studentIDs []string) ([]model.AchievementReference, error)
	GetAll(ctx context.Context) ([]model.AchievementReference, error)
}

/*
=====================================================
MONGODB ACHIEVEMENT (DETAIL)
=====================================================
*/
type IAchievementMongoRepository interface {
	Create(ctx context.Context, data model.AchievementMongo) (primitive.ObjectID, error)

	FindById(ctx context.Context, id primitive.ObjectID) (model.AchievementMongo, error)
	FindManyByIDs(ctx context.Context, ids []primitive.ObjectID) (map[string]model.AchievementMongo, error)

	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	AddAttachment(ctx context.Context, id primitive.ObjectID, file model.AttachmentFile) error
	SoftDelete(ctx context.Context, id primitive.ObjectID) error
}
