package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"clean-arch/app/model/postgre"
	"clean-arch/database"
	"github.com/google/uuid"
)

func CreateAchievementReference(ctx context.Context, ref *postgre.AchievementReference) error {
	if ref.ID == "" {
		ref.ID = uuid.New().String()
	}
	now := time.Now()
	ref.CreatedAt = now
	ref.UpdatedAt = now
	q := `INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := database.PostgresDB.ExecContext(ctx, q, ref.ID, ref.StudentID, ref.MongoAchievementID, ref.Status, ref.SubmittedAt, ref.VerifiedAt, ref.VerifiedBy, ref.RejectionNote, ref.CreatedAt, ref.UpdatedAt)
	return err
}

func GetAchievementReferenceByID(ctx context.Context, id string) (*postgre.AchievementReference, error) {
	q := `SELECT id,student_id,mongo_achievement_id,status,submitted_at,verified_at,verified_by,rejection_note,created_at,updated_at FROM achievement_references WHERE id=$1`
	row := database.PostgresDB.QueryRowContext(ctx, q, id)
	var ref postgre.AchievementReference
	var submitted, verified sql.NullTime
	var verifiedBy, rejectionNote sql.NullString
	if err := row.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &submitted, &verified, &verifiedBy, &rejectionNote, &ref.CreatedAt, &ref.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	if submitted.Valid { t := submitted.Time; ref.SubmittedAt = &t }
	if verified.Valid { t := verified.Time; ref.VerifiedAt = &t }
	if verifiedBy.Valid { v := verifiedBy.String; ref.VerifiedBy = &v }
	if rejectionNote.Valid { v := rejectionNote.String; ref.RejectionNote = &v }
	return &ref, nil
}

func UpdateAchievementReferenceStatus(ctx context.Context, id string, status string, updaterID *string, note *string) error {
	now := time.Now()
	q := `UPDATE achievement_references SET status=$1, verified_by=$2, rejection_note=$3, verified_at=$4, updated_at=$5 WHERE id=$6`
	var verifiedAt interface{}
	if status == "verified" {
		verifiedAt = now
	} else {
		verifiedAt = nil
	}
	_, err := database.PostgresDB.ExecContext(ctx, q, status, updaterID, note, verifiedAt, now, id)
	return err
}
