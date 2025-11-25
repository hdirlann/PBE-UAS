package repository

import (
	"context"

	"clean-arch/app/database"
)

func AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	q := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`
	_, err := database.PostgresDB.ExecContext(ctx, q, roleID, permissionID)
	return err
}

func ListPermissionsByRole(ctx context.Context, roleID string) ([]string, error) {
	q := `SELECT p.name FROM permissions p JOIN role_permissions rp ON rp.permission_id = p.id WHERE rp.role_id=$1`
	rows, err := database.PostgresDB.QueryContext(ctx, q, roleID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil { return nil, err }
		out = append(out, name)
	}
	return out, nil
}
