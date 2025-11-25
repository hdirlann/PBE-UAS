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

func CreateRole(ctx context.Context, r *postgre.Role) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	r.CreatedAt = time.Now()
	q := `INSERT INTO roles (id, name, description, created_at) VALUES ($1,$2,$3,$4)`
	_, err := database.PostgresDB.ExecContext(ctx, q, r.ID, r.Name, r.Desc, r.CreatedAt)
	return err
}

func GetRoleByName(ctx context.Context, name string) (*postgre.Role, error) {
	var r postgre.Role
	q := `SELECT id,name,description,created_at FROM roles WHERE name=$1`
	row := database.PostgresDB.QueryRowContext(ctx, q, name)
	if err := row.Scan(&r.ID, &r.Name, &r.Desc, &r.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	return &r, nil
}

func GetRoleByID(ctx context.Context, id string) (*postgre.Role, error) {
	var r postgre.Role
	q := `SELECT id,name,description,created_at FROM roles WHERE id=$1`
	row := database.PostgresDB.QueryRowContext(ctx, q, id)
	if err := row.Scan(&r.ID, &r.Name, &r.Desc, &r.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	return &r, nil
}
