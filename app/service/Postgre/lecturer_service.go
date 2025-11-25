package postgre

import (
	"context"

	"clean-arch/app/model/postgre"
	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
)

func CreateLecturerService(c *fiber.Ctx) error {
	var l postgre.Lecturer
	if err := c.BodyParser(&l); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if l.LecturerID == "" || l.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"lecturer_id and user_id required"})
	}
	if err := repository.CreateLecturer(context.Background(), &l); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"lecturerId": l.ID})
}

func GetLecturerService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"id required"}) }
	l, err := repository.GetLecturerByID(context.Background(), id)
	if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()}) }
	if l == nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":"not found"}) }
	return c.JSON(l)
}
