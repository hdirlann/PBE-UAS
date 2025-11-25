package repository

import (
	"context"
	"time"

	mongoModel "clean-arch/app/model/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const achievementsCollection = "achievements"

// CreateAchievement inserts a new achievement document.
func CreateAchievement(db *mgo.Database, a *mongoModel.Achievement) (*mongoModel.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := db.Collection(achievementsCollection)
	a.ID = primitive.NewObjectID()
	a.CreatedAt = time.Now()
	a.UpdatedAt = a.CreatedAt

	_, err := col.InsertOne(ctx, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// GetAchievementByID fetches a single achievement by hex id.
func GetAchievementByID(db *mgo.Database, hexID string) (*mongoModel.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := db.Collection(achievementsCollection)
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return nil, err
	}
	var out mongoModel.Achievement
	if err := col.FindOne(ctx, bson.M{"_id": oid, "deletedAt": bson.M{"$exists": false}}).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateAchievement updates fields on an achievement.
func UpdateAchievement(db *mgo.Database, hexID string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := db.Collection(achievementsCollection)
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}
	update["updatedAt"] = time.Now()
	_, err = col.UpdateOne(ctx, bson.M{"_id": oid, "deletedAt": bson.M{"$exists": false}}, bson.M{"$set": update})
	return err
}

// SoftDeleteAchievement sets deletedAt timestamp.
func SoftDeleteAchievement(db *mgo.Database, hexID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := db.Collection(achievementsCollection)
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}
	_, err = col.UpdateOne(ctx, bson.M{"_id": oid, "deletedAt": bson.M{"$exists": false}}, bson.M{
		"$set": bson.M{"deletedAt": time.Now(), "updatedAt": time.Now()},
	})
	return err
}

// HardDeleteAchievement permanently removes a document (use carefully).
func HardDeleteAchievement(db *mgo.Database, hexID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := db.Collection(achievementsCollection)
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}
	_, err = col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// ListAchievements lists achievements by filters and supports pagination.
func ListAchievements(db *mgo.Database, filter bson.M, page, limit int64) ([]mongoModel.Achievement, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := db.Collection(achievementsCollection)

	if filter == nil {
		filter = bson.M{}
	}
	// Exclude deleted
	filter["deletedAt"] = bson.M{"$exists": false}

	total, err := col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.M{"createdAt": -1})

	cur, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var out []mongoModel.Achievement
	if err := cur.All(ctx, &out); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}
