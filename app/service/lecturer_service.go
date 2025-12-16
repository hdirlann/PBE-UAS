package service

import (
	"context"

	"clean-arch/app/model"
	"clean-arch/app/repository"
	"clean-arch/database"

	"github.com/gofiber/fiber/v2"
)

// CreateLecturerService
// @Summary Create lecturer
// @Tags Lecturers
// @Description Admin creates a lecturer profile.
// @Accept json
// @Produce json
// @Param body body object true "Lecturer body" example({"userId":"user-do-1","lecturerId":"LECT-2025-01","department":"Teknik"})
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /lecturers [post]
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

// GetLecturerService
// @Summary Get lecturer by id
// @Tags Lecturers
// @Description Get lecturer details by id.
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {object} model.Lecturer
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /lecturers/{id} [get]
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

// ListLecturersService - GET /lecturers
func ListLecturersService(c *fiber.Ctx) error {
	q := `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers ORDER BY created_at DESC`
	rows, err := database.PostgresDB.QueryContext(context.Background(), q)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()
	var out []model.Lecturer
	for rows.Next() {
		var l model.Lecturer
		if err := rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
			continue
		}
		out = append(out, l)
	}
	return c.JSON(out)
}

// GetLecturerAdviseesService - GET /lecturers/:id/advisees
func GetLecturerAdviseesService(c *fiber.Ctx) error {
	lectID := c.Params("id")
	if lectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "lecturer id required"})
	}
	l, err := repository.GetLecturerByID(context.Background(), lectID)
	if err != nil || l == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "lecturer not found"})
	}
	advisees, err := repository.ListStudentsByAdvisor(context.Background(), l.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(advisees)
}
