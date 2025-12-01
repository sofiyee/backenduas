package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementMongo struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty"`
	StudentID       string                 `bson:"studentId"`
	Title           string                 `bson:"title"`
	Description     string                 `bson:"description"`
	AchievementType string                 `bson:"achievementType"`
	Details         map[string]interface{} `bson:"details"`
	Attachments     []AttachmentFile       `bson:"attachments"`
	Tags            []string               `bson:"tags"`
	Points          int                    `bson:"points"`
	CreatedAt       int64                  `bson:"createdAt"`
	UpdatedAt       int64                  `bson:"updatedAt"`
	Status          string                 `bson:"status"` // draft, deleted
}

type AttachmentFile struct {
	FileName   string `bson:"fileName"`
	FileURL    string `bson:"fileUrl"`
	FileType   string `bson:"fileType"`
	UploadedAt int64  `bson:"uploadedAt"`
}
