package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"clean-arch/app/model/postgre"
	"clean-arch/app/database"
	"github.com/google/uuid"
)

func CreateLecturer(ctx context.Context, l *postgre.Lecturer) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	l.CreatedAt = time.Now()
	q := `INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
		  VALUES ($1,$2,$3,$4,$5)`
	_, err := database.PostgresDB.ExecContext(ctx, q, l.ID, l.UserID, l.LecturerID, l.Department, l.CreatedAt)
	return err
}

func GetLecturerByID(ctx context.Context, id string) (*postgre.Lecturer, error) {
	var l postgre.Lecturer
	q := `SELECT id,user_id,lecturer_id,department,created_at FROM lecturers WHERE id=$1`
	row := database.PostgresDB.QueryRowContext(ctx, q, id)
	if err := row.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	return &l, nil
}
