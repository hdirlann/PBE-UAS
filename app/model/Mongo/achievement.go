package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	StudentID       string                 `bson:"studentId" json:"studentId"`
	AchievementType string                 `bson:"achievementType" json:"achievementType"` // academic, competition, organization, publication, certification, other
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`
	Details         map[string]interface{} `bson:"details,omitempty" json:"details,omitempty"` // dynamic fields
	Attachments     []Attachment           `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Tags            []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	Points          *int                   `bson:"points,omitempty" json:"points,omitempty"`
	CreatedAt       time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time              `bson:"updatedAt" json:"updatedAt"`
	DeletedAt       *time.Time             `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}
