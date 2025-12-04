package repository

import (
	"context"
	"time"

	mongoModel "clean-arch/app/model"

	"go.mongodb.org/mongo-driver/bson"
	mgo "go.mongodb.org/mongo-driver/mongo"
)

const attachmentsCollection = "attachments"

// AddAttachment inserts an attachment record.
func AddAttachment(db *mgo.Database, a *mongoModel.Attachment) (*mongoModel.Attachment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := db.Collection(attachmentsCollection)
	a.UploadedAt = time.Now()
	_, err := col.InsertOne(ctx, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// ListAttachmentsByAchievement returns attachments for a given achievement id.
func ListAttachmentsByAchievement(db *mgo.Database, achievementID string) ([]mongoModel.Attachment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := db.Collection(attachmentsCollection)
	cur, err := col.Find(ctx, bson.M{"achievementId": achievementID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []mongoModel.Attachment
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
