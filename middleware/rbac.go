package middleware

import (
	"context"
	"log"

	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
)

func hasPerm(perms []string, perm string) bool {
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}

// RequirePermission returns a fiber.Handler that enforces the given permission string.
func RequirePermission(perm string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1) fast path: permissions present in token
		if v := c.Locals(LocalsPermissions); v != nil {
			if perms, ok := v.([]string); ok {
				if hasPerm(perms, perm) {
					return c.Next()
				}
			}
		}

		// 2) get role id from locals
		rv := c.Locals(LocalsRoleID)
		if rv == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "missing role info"})
		}
		roleID, ok := rv.(string)
		if !ok || roleID == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "invalid role info"})
		}

		// 3) check cache
		if cp, ok := GetCachedPerms(roleID); ok {
			if hasPerm(cp, perm) {
				return c.Next()
			}
			log.Printf("[rbac] deny (cache) role=%s need=%s", roleID, perm)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
		}

		// 4) fallback load from repository
		perms, err := repository.ListPermissionsByRole(context.Background(), roleID)
		if err != nil {
			log.Printf("[rbac] failed load perms role=%s err=%v", roleID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "server error"})
		}
		// cache result for next calls
		SetCachedPerms(roleID, perms)

		if hasPerm(perms, perm) {
			return c.Next()
		}
		log.Printf("[rbac] deny role=%s need=%s", roleID, perm)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
	}
}
