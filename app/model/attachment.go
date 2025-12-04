package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attachment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AchievementID string             `bson:"achievementId" json:"achievementId"`
	FileName      string             `bson:"fileName" json:"fileName"`
	FileURL       string             `bson:"fileUrl" json:"fileUrl"`
	FileType      string             `bson:"fileType" json:"fileType"`
	UploadedAt    time.Time          `bson:"uploadedAt" json:"uploadedAt"`
}
