package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"clean-arch/app/model"
	"clean-arch/database"

	"github.com/google/uuid"
)

// CreateAchievementReference inserts a new achievement reference row.
func CreateAchievementReference(ctx context.Context, r *model.AchievementReference) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	now := time.Now()
	if r.CreatedAt.IsZero() {
		r.CreatedAt = now
	}
	r.UpdatedAt = now

	q := `INSERT INTO achievement_references
	      (id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	var submitted, verified interface{}
	if r.SubmittedAt != nil {
		submitted = *r.SubmittedAt
	} else {
		submitted = nil
	}
	if r.VerifiedAt != nil {
		verified = *r.VerifiedAt
	} else {
		verified = nil
	}

	_, err := database.PostgresDB.ExecContext(ctx, q,
		r.ID, r.StudentID, r.MongoAchievementID, r.Status, submitted, verified, r.VerifiedBy, r.RejectionNote, r.CreatedAt, r.UpdatedAt,
	)
	return err
}

// GetAchievementReferenceByMongoID finds a reference row by mongo_achievement_id
func GetAchievementReferenceByMongoID(ctx context.Context, mongoID string) (*model.AchievementReference, error) {
	q := `SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
	      FROM achievement_references WHERE mongo_achievement_id=$1 LIMIT 1`
	row := database.PostgresDB.QueryRowContext(ctx, q, mongoID)

	var ref model.AchievementReference
	var submitted, verified sql.NullTime
	var verifiedBy, rejectionNote sql.NullString

	if err := row.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &submitted, &verified, &verifiedBy, &rejectionNote, &ref.CreatedAt, &ref.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if submitted.Valid {
		t := submitted.Time
		ref.SubmittedAt = &t
	}
	if verified.Valid {
		t := verified.Time
		ref.VerifiedAt = &t
	}
	if verifiedBy.Valid {
		v := verifiedBy.String
		ref.VerifiedBy = &v
	}
	if rejectionNote.Valid {
		v := rejectionNote.String
		ref.RejectionNote = &v
	}
	return &ref, nil
}

// UpdateAchievementReferenceStatus updates status and optional verifier/note
func UpdateAchievementReferenceStatus(ctx context.Context, referenceID, status string, verifierID *string, rejectionNote *string) error {
	now := time.Now()
	// We'll set fields depending on status:
	// - 'submitted': set status, submitted_at=now, updated_at
	// - 'verified': set status, verified_at=now, verified_by=verifierID, updated_at
	// - 'rejected': set status, rejection_note, verified_by=verifierID, updated_at
	// - other: just update status and updated_at
	var q string
	var args []interface{}

	switch status {
	case "submitted":
		q = `UPDATE achievement_references SET status=$1, submitted_at=$2, updated_at=$3 WHERE id=$4`
		args = []interface{}{status, now, now, referenceID}
	case "verified":
		q = `UPDATE achievement_references SET status=$1, verified_at=$2, verified_by=$3, updated_at=$4 WHERE id=$5`
		vby := sql.NullString{}
		if verifierID != nil {
			vby = sql.NullString{String: *verifierID, Valid: true}
		}
		args = []interface{}{status, now, vby, now, referenceID}
	case "rejected":
		q = `UPDATE achievement_references SET status=$1, rejection_note=$2, verified_by=$3, updated_at=$4 WHERE id=$5`
		rn := sql.NullString{}
		if rejectionNote != nil {
			rn = sql.NullString{String: *rejectionNote, Valid: true}
		}
		vby := sql.NullString{}
		if verifierID != nil {
			vby = sql.NullString{String: *verifierID, Valid: true}
		}
		args = []interface{}{status, rn, vby, now, referenceID}
	default:
		q = `UPDATE achievement_references SET status=$1, updated_at=$2 WHERE id=$3`
		args = []interface{}{status, now, referenceID}
	}

	_, err := database.PostgresDB.ExecContext(ctx, q, args...)
	return err
}
