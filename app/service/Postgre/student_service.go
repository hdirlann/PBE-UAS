package postgre

import (
	"context"

	"clean-arch/app/model/postgre"
	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
)

func CreateStudentService(c *fiber.Ctx) error {
	var s postgre.Student
	if err := c.BodyParser(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if s.StudentID == "" || s.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "student_id and user_id required"})
	}
	if err := repository.CreateStudent(context.Background(), &s); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"studentId": s.ID})
}

func GetStudentService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"id required"}) }
	s, err := repository.GetStudentByID(context.Background(), id)
	if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()}) }
	if s == nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":"not found"}) }
	return c.JSON(s)
}

func ListStudentsByAdvisorService(c *fiber.Ctx) error {
	advisor := c.Params("advisorId")
	if advisor == "" { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"advisorId required"}) }
	out, err := repository.ListStudentsByAdvisor(context.Background(), advisor)
	if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(out)
}
