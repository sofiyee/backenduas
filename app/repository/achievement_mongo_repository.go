package repository

import (
	"context"
	"time"

	"backenduas/app/model"
	"backenduas/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementMongoRepository struct{}

func NewAchievementMongoRepository() *AchievementMongoRepository {
	return &AchievementMongoRepository{}
}

func (r *AchievementMongoRepository) collection() *mongo.Collection {
	return database.MongoDB.Collection("achievements")
}

func (r *AchievementMongoRepository) Create(ctx context.Context, ach model.AchievementMongo) (primitive.ObjectID, error) {
	ach.CreatedAt = time.Now().Unix()
	ach.UpdatedAt = time.Now().Unix()
	ach.Status = "draft"

	res, err := r.collection().InsertOne(ctx, ach)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *AchievementMongoRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now().Unix()

	_, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)

	return err
}

func (r *AchievementMongoRepository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"status":    "deleted",
				"updatedAt": time.Now().Unix(),
			},
		},
	)
	return err
}

func (r *AchievementMongoRepository) FindById(ctx context.Context, id primitive.ObjectID) (model.AchievementMongo, error) {
	var ach model.AchievementMongo
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&ach)
	return ach, err
}

func (r *AchievementMongoRepository) FindManyByIDs(ctx context.Context, ids []primitive.ObjectID) (map[string]model.AchievementMongo, error) {

	cursor, err := r.collection().Find(ctx, bson.M{
		"_id": bson.M{"$in": ids},
	})
	if err != nil {
		return nil, err
	}

	result := make(map[string]model.AchievementMongo)

	for cursor.Next(ctx) {
		var ach model.AchievementMongo
		cursor.Decode(&ach)
		result[ach.ID.Hex()] = ach
	}

	return result, nil
}
