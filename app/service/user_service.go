package service

import (
	"context"
	"time"

	"clean-arch/app/model"
	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserService
// @Summary Create user (admin)
// @Tags Users
// @Description Create a new user (admin only).
// @Accept json
// @Produce json
// @Param body body object true "User body" example({"username":"newuser","email":"a@b.com","password":"secret","fullName":"New User","roleId":"role-mahasiswa"})
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /users [post]
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

// GetUserByIDService
// @Summary Get user by id
// @Tags Users
// @Description Get user profile by user id.
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /users/{id} [get]
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

// ListUsersService
// @Summary List users
// @Tags Users
// @Description List users (admin).
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /users [get]
func ListUsersService(c *fiber.Ctx) error {
	users, err := repository.ListUsers(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}

// DeleteUserService
// @Summary Delete user
// @Tags Users
// @Description Delete user by id (admin).
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /users/{id} [delete]
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

// UpdateUserService - partial update (PUT /users/:id)
// Tambahan untuk SRS: update user fields (admin)
func UpdateUserService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	var payload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"fullName"`
		RoleID   string `json:"roleId"`
		IsActive *bool  `json:"isActive"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	u, err := repository.GetUserByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if u == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	if payload.Username != "" {
		u.Username = payload.Username
	}
	if payload.Email != "" {
		u.Email = payload.Email
	}
	if payload.FullName != "" {
		u.FullName = payload.FullName
	}
	if payload.RoleID != "" {
		u.RoleID = payload.RoleID
	}
	if payload.IsActive != nil {
		u.IsActive = *payload.IsActive
	}
	u.UpdatedAt = time.Now()
	if err := repository.UpdateUser(context.Background(), u); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

// AssignRoleService - set roleId for a user
func AssignRoleService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	var body struct {
		RoleID string `json:"roleId"`
	}
	if err := c.BodyParser(&body); err != nil || body.RoleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "roleId required"})
	}
	u, err := repository.GetUserByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if u == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	u.RoleID = body.RoleID
	u.UpdatedAt = time.Now()
	if err := repository.UpdateUser(context.Background(), u); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "role updated"})
}
