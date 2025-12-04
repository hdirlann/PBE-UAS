package service

import (
	"context"

	"clean-arch/app/model"
	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// ===========================
// Create User (Admin Only)
// ===========================
func CreateUserService(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
		RoleID   string `json:"roleId"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	if body.Username == "" || body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "username, email, password required"})
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "failed to hash password"})
	}

	user := &model.User{
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: string(hashed),
		FullName:     body.FullName,
		RoleID:       body.RoleID,
		IsActive:     true,
	}

	if err := repository.CreateUser(context.Background(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).
		JSON(fiber.Map{"userId": user.ID})
}

// ===========================
// Get User by ID
// ===========================
func GetUserByIDService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "user id is required"})
	}

	user, err := repository.GetUserByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

// ===========================
// List Users
// ===========================
func ListUsersService(c *fiber.Ctx) error {
	users, err := repository.ListUsers(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}

// ===========================
// Delete User
// ===========================
func DeleteUserService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "user id is required"})
	}

	if err := repository.DeleteUser(context.Background(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "user deleted"})
}
