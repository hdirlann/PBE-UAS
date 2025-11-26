package service

import (
	"context"

	"clean-arch/app/model/postgre"
	"clean-arch/app/repository/postgre"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserService - admin creates a user
func CreateUserService(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
		RoleID   string `json:"roleId"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if body.Username == "" || body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username/email/password required"})
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed hash"})
	}
	u := &postgre.User{
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: string(hashed),
		FullName:     body.FullName,
		RoleID:       body.RoleID,
		IsActive:     true,
	}
	if err := repository.CreateUser(context.Background(), u); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"userId": u.ID})
}

// AuthenticateService - login with username/email + password
func AuthenticateService(c *fiber.Ctx) error {
	var body struct {
		Identifier string `json:"identifier"` // username or email
		Password   string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	u, err := repository.GetUserByUsernameOrEmail(context.Background(), body.Identifier)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	// TODO: generate JWT and include permissions. For now return profile.
	return c.JSON(fiber.Map{"user": u})
}
