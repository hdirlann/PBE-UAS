package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Locals keys (dipakai di middleware lain / handlers)
const (
	LocalsUserID      = "user_id"
	LocalsUsername    = "username"
	LocalsRoleID      = "role_id"
	LocalsPermissions = "permissions"
)

// JWTMiddleware extracts Bearer token, validates it, and stores claims in c.Locals
func JWTMiddleware() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-this-secret"
	}

	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header"})
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// only HMAC supported
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}

		// copy expected claims into locals
		if sub, ok := claims["sub"].(string); ok {
			c.Locals(LocalsUserID, sub)
		}
		if un, ok := claims["username"].(string); ok {
			c.Locals(LocalsUsername, un)
		}
		if rid, ok := claims["role_id"].(string); ok {
			c.Locals(LocalsRoleID, rid)
		}

		// normalize permissions into []string and always set it (even empty)
		var permsOut []string
		if rawP, ok := claims["permissions"]; ok && rawP != nil {
			switch v := rawP.(type) {
			case []interface{}:
				out := make([]string, 0, len(v))
				for _, it := range v {
					if s, ok := it.(string); ok {
						out = append(out, s)
					} else {
						out = append(out, fmt.Sprintf("%v", it))
					}
				}
				permsOut = out
			case []string:
				permsOut = v
			case string:
				// support comma separated string
				for _, s := range strings.Split(v, ",") {
					permsOut = append(permsOut, strings.TrimSpace(s))
				}
			default:
				// fallback: try to stringify
				permsOut = []string{fmt.Sprintf("%v", v)}
			}
		} else {
			permsOut = []string{}
		}
		c.Locals(LocalsPermissions, permsOut)

		return c.Next()
	}
}
