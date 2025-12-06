package service

import (
	"context"

	"clean-arch/app/model"
	"clean-arch/app/repository"

	"github.com/gofiber/fiber/v2"
)

// CreateStudentService - admin creates a student
// expects JSON body matching model.Student (student_id, user_id, program_study, academic_year, advisor_id)
func CreateStudentService(c *fiber.Ctx) error {
	var s model.Student
	if err := c.BodyParser(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// minimal validation: user_id and student_id required
	if s.UserID == "" || s.StudentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "student_id and user_id required"})
	}

	if err := repository.CreateStudent(context.Background(), &s); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"studentId": s.ID})
}

// GetStudentService - GET /api/v1/students/:id
func GetStudentService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	s, err := repository.GetStudentByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if s == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(s)
}

// ListStudentsByAdvisorService - GET /api/v1/students/advisor/:advisorId
func ListStudentsByAdvisorService(c *fiber.Ctx) error {
	advisor := c.Params("advisorId")
	if advisor == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "advisorId required"})
	}
	out, err := repository.ListStudentsByAdvisor(context.Background(), advisor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(out)
}
