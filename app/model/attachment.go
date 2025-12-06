package model

import "time"

type Attachment struct {
    ID            string    `bson:"_id,omitempty" json:"id"`
    AchievementID string    `bson:"achievement_id" json:"achievement_id"`
    FileName      string    `bson:"file_name" json:"file_name"`
    FileURL       string    `bson:"file_url" json:"file_url"`
    FileType      string    `bson:"file_type" json:"file_type"`
    CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}
