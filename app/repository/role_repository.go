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

// CreateRole inserts a new role (sets created_at & updated_at)
func CreateRole(ctx context.Context, r *model.Role) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now

	q := `INSERT INTO roles (id, name, description, created_at, updated_at)
		  VALUES ($1,$2,$3,$4,$5)`
	_, err := database.PostgresDB.ExecContext(ctx, q, r.ID, r.Name, r.Desc, r.CreatedAt, r.UpdatedAt)
	return err
}

// GetRoleByName returns role by name
func GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var r model.Role
	q := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name=$1`
	row := database.PostgresDB.QueryRowContext(ctx, q, name)
	if err := row.Scan(&r.ID, &r.Name, &r.Desc, &r.CreatedAt, &r.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}

// GetRoleByID returns role by id
func GetRoleByID(ctx context.Context, id string) (*model.Role, error) {
	var r model.Role
	q := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id=$1`
	row := database.PostgresDB.QueryRowContext(ctx, q, id)
	if err := row.Scan(&r.ID, &r.Name, &r.Desc, &r.CreatedAt, &r.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}
