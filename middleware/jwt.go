package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	LocalsUserID      = "user_id"
	LocalsUsername    = "username"
	LocalsRoleID      = "role_id"
	LocalsPermissions = "permissions"
)

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

		if sub, ok := claims["sub"].(string); ok {
			c.Locals(LocalsUserID, sub)
		}
		if un, ok := claims["username"].(string); ok {
			c.Locals(LocalsUsername, un)
		}
		if rid, ok := claims["role_id"].(string); ok {
			c.Locals(LocalsRoleID, rid)
		}
		// permissions may be []interface{} or []string
		if perms, ok := claims["permissions"].([]interface{}); ok {
			out := make([]string, 0, len(perms))
			for _, p := range perms {
				if s, ok := p.(string); ok {
					out = append(out, s)
				}
			}
			c.Locals(LocalsPermissions, out)
		} else if permsS, ok := claims["permissions"].([]string); ok {
			c.Locals(LocalsPermissions, permsS)
		}

		return c.Next()
	}
}
