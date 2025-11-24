package mongo

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
    ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
    StudentID       string                 `bson:"studentId" json:"studentId"` // UUID from Postgres students.id
    AchievementType string                 `bson:"achievementType" json:"achievementType"`
    Title           string                 `bson:"title" json:"title"`
    Description     string                 `bson:"description" json:"description"`
    Details         map[string]interface{} `bson:"details" json:"details"` // dynamic fields per SRS
    Attachments     []Attachment           `bson:"attachments" json:"attachments"`
    Tags            []string               `bson:"tags" json:"tags"`
    Points          *int                   `bson:"points,omitempty" json:"points,omitempty"`
    CreatedAt       time.Time              `bson:"createdAt" json:"createdAt"`
    UpdatedAt       time.Time              `bson:"updatedAt" json:"updatedAt"`
    Deleted         bool                   `bson:"deleted" json:"deleted"` // soft-delete flag for FR-005
}
