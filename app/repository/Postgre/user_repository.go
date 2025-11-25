package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"clean-arch/app/database"
	"clean-arch/app/model/postgre"

	"github.com/google/uuid"
)

func CreateUser(ctx context.Context, u *postgre.User) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	q := `INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := database.PostgresDB.ExecContext(ctx, q,
		u.ID, u.Username, u.Email, u.PasswordHash, u.FullName, u.RoleID, u.IsActive, u.CreatedAt, u.UpdatedAt)
	return err
}

func GetUserByID(ctx context.Context, id string) (*postgre.User, error) {
	var u postgre.User
	q := `SELECT id,username,email,password_hash,full_name,role_id,is_active,created_at,updated_at FROM users WHERE id=$1`
	row := database.PostgresDB.QueryRowContext(ctx, q, id)
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*postgre.User, error) {
	var u postgre.User
	q := `SELECT id,username,email,password_hash,full_name,role_id,is_active,created_at,updated_at FROM users WHERE username=$1 OR email=$1`
	row := database.PostgresDB.QueryRowContext(ctx, q, identifier)
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	return &u, nil
}

func UpdateUser(ctx context.Context, u *postgre.User) error {
	u.UpdatedAt = time.Now()
	q := `UPDATE users SET username=$1,email=$2,password_hash=$3,full_name=$4,role_id=$5,is_active=$6,updated_at=$7 WHERE id=$8`
	_, err := database.PostgresDB.ExecContext(ctx, q, u.Username, u.Email, u.PasswordHash, u.FullName, u.RoleID, u.IsActive, u.UpdatedAt, u.ID)
	return err
}

func DeleteUser(ctx context.Context, id string) error {
	q := `DELETE FROM users WHERE id=$1`
	_, err := database.PostgresDB.ExecContext(ctx, q, id)
	return err
}
