package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"clean-arch/database"
	"clean-arch/app/model"

	"github.com/google/uuid"
)

func CreateStudent(ctx context.Context, s *model.Student) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.CreatedAt = time.Now()
	q := `INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
	      VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := database.PostgresDB.ExecContext(ctx, q, s.ID, s.UserID, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID, s.CreatedAt)
	return err
}

func GetStudentByID(ctx context.Context, id string) (*model.Student, error) {
	var s model.Student
	q := `SELECT id,user_id,student_id,program_study,academic_year,advisor_id,created_at FROM students WHERE id=$1`
	row := database.PostgresDB.QueryRowContext(ctx, q, id)
	var advisor sql.NullString
	if err := row.Scan(&s.ID,&s.UserID,&s.StudentID,&s.ProgramStudy,&s.AcademicYear,&advisor,&s.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	if advisor.Valid { v := advisor.String; s.AdvisorID = &v }
	return &s, nil
}

func ListStudentsByAdvisor(ctx context.Context, advisorID string) ([]model.Student, error) {
	q := `SELECT id,user_id,student_id,program_study,academic_year,advisor_id,created_at FROM students WHERE advisor_id=$1`
	rows, err := database.PostgresDB.QueryContext(ctx, q, advisorID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []model.Student
	for rows.Next() {
		var s model.Student
		var advisor sql.NullString
		if err := rows.Scan(&s.ID,&s.UserID,&s.StudentID,&s.ProgramStudy,&s.AcademicYear,&advisor,&s.CreatedAt); err != nil {
			return nil, err
		}
		if advisor.Valid { v := advisor.String; s.AdvisorID = &v }
		out = append(out, s)
	}
	return out, nil
}
