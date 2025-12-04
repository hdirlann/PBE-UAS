package middleware

import (
	"context"

	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
)

// RequirePermission returns a Fiber handler that enforces permission check
func RequirePermission(perm string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// check claims first
		if v := c.Locals(LocalsPermissions); v != nil {
			if perms, ok := v.([]string); ok {
				for _, p := range perms {
					if p == perm {
						return c.Next()
					}
				}
			}
		}

		// get role id
		roleIDv := c.Locals(LocalsRoleID)
		if roleIDv == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "missing role info"})
		}
		roleID, _ := roleIDv.(string)

		// try cache
		if cp, ok := GetCachedPerms(roleID); ok {
			for _, p := range cp {
				if p == perm {
					return c.Next()
				}
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
		}

		// load from DB via repository
		perms, err := repository.ListPermissionsByRole(context.Background(), roleID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "server error"})
		}
		SetCachedPerms(roleID, perms)

		for _, p := range perms {
			if p == perm {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
	}
}
