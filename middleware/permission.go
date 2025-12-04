package middleware

import (
	"context"
	"sync"
	"time"

	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// simple in-memory cache for role permissions
type rolePermCacheEntry struct {
	perms  []string
	expire time.Time
}

var rolePermCache = struct {
	sync.RWMutex
	m map[string]rolePermCacheEntry
}{m: make(map[string]rolePermCacheEntry)}

// cache TTL
const cacheTTL = 5 * time.Minute

// RequirePermission returns a fiber.Handler that checks permission in token claims or loads from DB
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claimsI := c.Locals("user_claims")
		if claimsI == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing claims"})
		}

		var mc jwt.MapClaims
		switch t := claimsI.(type) {
		case jwt.MapClaims:
			mc = t
		case *jwt.MapClaims:
			mc = *t
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid claims"})
		}

		// first try permissions from claims
		if pI, ok := mc["permissions"]; ok && pI != nil {
			switch v := pI.(type) {
			case []string:
				for _, p := range v {
					if p == permission {
						return c.Next()
					}
				}
			case []interface{}:
				for _, it := range v {
					if s, ok := it.(string); ok && s == permission {
						return c.Next()
					}
				}
			}
		}

		// if not present in claims, fallback to load by role_id
		roleID, _ := mc["role_id"].(string)
		if roleID == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permission denied"})
		}

		// check cache
		rolePermCache.RLock()
		entry, ok := rolePermCache.m[roleID]
		rolePermCache.RUnlock()
		if ok && time.Now().Before(entry.expire) {
			for _, p := range entry.perms {
				if p == permission {
					return c.Next()
				}
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permission denied"})
		}

		// load from DB
		perms, err := postgres.ListPermissionsByRole(context.Background(), roleID)
		if err != nil {
			// jika error, kembalikan 500 atau deny; di sini kita pilih 500
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to load permissions"})
		}

		// update cache
		rolePermCache.Lock()
		rolePermCache.m[roleID] = rolePermCacheEntry{perms: perms, expire: time.Now().Add(cacheTTL)}
		rolePermCache.Unlock()

		for _, p := range perms {
			if p == permission {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permission denied"})
	}
}
