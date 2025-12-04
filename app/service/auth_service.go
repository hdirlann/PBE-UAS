package service

import (
	"context"
	"os"
	"strconv"
	"time"

	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticateService implements FR-001 (login: identifier + password)
func AuthenticateService(c *fiber.Ctx) error {
	var body struct {
		Identifier string `json:"identifier"` // username or email
		Password   string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input", "detail": err.Error()})
	}

	// ambil user (username atau email)
	u, err := repository.GetUserByUsernameOrEmail(context.Background(), body.Identifier)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "server error"})
	}
	if u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// cek aktif
	if !u.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "user not active"})
	}

	// ambil permissions berdasarkan role (bisa kosong)
	perms, err := repository.ListPermissionsByRole(context.Background(), u.RoleID)
	if err != nil {
		// jangan crash app hanya karena read perms gagal; fallback ke empty slice
		perms = []string{}
	}

	// JWT config
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "change-this-secret"
	}
	expHours := 24
	if v := os.Getenv("JWT_EXPIRES_HOURS"); v != "" {
		if hh, e := strconv.Atoi(v); e == nil && hh > 0 {
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
