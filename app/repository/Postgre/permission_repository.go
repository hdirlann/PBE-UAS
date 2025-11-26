package repository

import (
	"context"
	"database/sql"
	"fmt"

	model "clean-arch/app/model/postgre" // import model Permission

	"github.com/google/uuid"
)

// DB is a package-level database handle. Set this from your app initialization.
var DB *sql.DB

// SetDB sets the package-level DB handle (call once during app startup).
func SetDB(db *sql.DB) {
	DB = db
}

// CreatePermission inserts a new permission record.
func CreatePermission(ctx context.Context, p *model.Permission) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	if p == nil {
		return fmt.Errorf("permission is nil")
	}
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	const q = `
		INSERT INTO permissions (id, name, resource, action, description)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := DB.ExecContext(ctx, q, p.ID, p.Name, p.Resource, p.Action, p.Description)
	return err
}

// GetPermissionByID returns a permission by ID.
func GetPermissionByID(ctx context.Context, id string) (*model.Permission, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	const q = `
		SELECT id, name, resource, action, description
		FROM permissions
		WHERE id = $1
	`
	row := DB.QueryRowContext(ctx, q, id)
	var p model.Permission
	if err := row.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// ListPermissions returns a list of permissions.
func ListPermissions(ctx context.Context, limit, offset int) ([]model.Permission, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	const q = `
		SELECT id, name, resource, action, description
		FROM permissions
		ORDER BY name
		LIMIT $1 OFFSET $2
	`
	rows, err := DB.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Permission
	for rows.Next() {
		var p model.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdatePermission updates fields of a permission by ID.
func UpdatePermission(ctx context.Context, p *model.Permission) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	if p == nil || p.ID == "" {
		return fmt.Errorf("permission or permission.ID is empty")
	}

	const q = `
		UPDATE permissions
		SET name = $1, resource = $2, action = $3, description = $4
		WHERE id = $5
	`
	res, err := DB.ExecContext(ctx, q, p.Name, p.Resource, p.Action, p.Description, p.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeletePermission deletes a permission by ID.
func DeletePermission(ctx context.Context, id string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	const q = `
		DELETE FROM permissions
		WHERE id = $1
	`
	res, err := DB.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}
