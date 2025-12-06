package service

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticateService implements FR-001 (login)
func AuthenticateService(c *fiber.Ctx) error {
	var body struct {
		Username   string `json:"username"`
		Identifier string `json:"identifier"` // alias: username or email
		Password   string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Printf("[auth] body parse error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// choose identifier: prefer Username if provided, else Identifier
	ident := body.Username
	if ident == "" {
		ident = body.Identifier
	}
	if ident == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username and password required"})
	}

	// fetch user (repo will check username OR email)
	u, err := repository.GetUserByUsernameOrEmail(context.Background(), ident)
	if err != nil {
		log.Printf("[auth] repo error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
	}
	if u == nil {
		log.Printf("[auth] user not found for: %s", ident)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(body.Password)); err != nil {
		log.Printf("[auth] bcrypt compare fail for user %s: %v", u.Username, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// check active
	if !u.IsActive {
		log.Printf("[auth] user not active: %s", u.Username)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "user not active"})
	}

	// load permissions for role
	perms, err := repository.ListPermissionsByRole(context.Background(), u.RoleID)
	if err != nil {
		log.Printf("[auth] failed load permissions for role %s: %v", u.RoleID, err)
		perms = []string{}
	}

	// jwt settings
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "change-this-secret"
	}
	expHours := 24
	if v := os.Getenv("JWT_EXPIRES_HOURS"); v != "" {
		if hh, err := strconv.Atoi(v); err == nil && hh > 0 {
			expHours = hh
		}
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":         u.ID,
		"username":    u.Username,
		"email":       u.Email,
		"role_id":     u.RoleID,
		"permissions": perms,
		"iat":         now.Unix(),
		"exp":         now.Add(time.Duration(expHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("[auth] jwt sign error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed generate token"})
	}

	profile := fiber.Map{
		"id":          u.ID,
		"username":    u.Username,
		"email":       u.Email,
		"full_name":   u.FullName,
		"role_id":     u.RoleID,
		"is_active":   u.IsActive,
		"created_at":  u.CreatedAt,
		"updated_at":  u.UpdatedAt,
		"permissions": perms,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": tokenStr,
		"user":  profile,
	})
}
