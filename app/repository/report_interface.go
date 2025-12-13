package repository

import (
	"context"

	"backenduas/app/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IReportPGRepository interface {
	GetByStudentID(ctx context.Context, studentID string) ([]model.AchievementReference, error)
	GetByStudentIDs(ctx context.Context, studentIDs []string) ([]model.AchievementReference, error)
}

type IReportMongoRepository interface {
	FindById(ctx context.Context, id primitive.ObjectID) (model.AchievementMongo, error)
}

type IReportStudentRepository interface {
	GetAllStudentIDs(ctx context.Context) ([]string, error)
	GetStudentIDByUserID(ctx context.Context, userID string) (string, error)
	GetLecturerIDByUserID(ctx context.Context, userID string) (string, error)
	GetStudentsByAdvisor(ctx context.Context, lecturerID string) ([]string, error)
	GetStudentsByIDs(ctx context.Context, ids []string) (map[string]model.StudentDetail, error)
	IsStudentUnderAdvisor(ctx context.Context, studentID, lecturerID string) (bool, error)
}
