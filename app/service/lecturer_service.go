package service

import (
	"context"

	"clean-arch/app/model"
	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
)

// CreateLecturerService (pastikan ada) - contoh singkat
func CreateLecturerService(c *fiber.Ctx) error {
	var body struct {
		UserID     string `json:"userId"`
		LecturerID string `json:"lecturerId"`
		Department string `json:"department"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if body.UserID == "" || body.LecturerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "userId and lecturerId required"})
	}
	l := &model.Lecturer{
		UserID:     body.UserID,
		LecturerID: body.LecturerID,
		Department: body.Department,
	}
	if err := repository.CreateLecturer(context.Background(), l); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"lecturerId": l.ID})
}

// GetLecturerService - needed by your routes
func GetLecturerService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	l, err := repository.GetLecturerByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if l == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(l)
}
