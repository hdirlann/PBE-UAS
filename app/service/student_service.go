package service

import (
	"context"
	"database/sql"

	"clean-arch/app/model"
	"clean-arch/app/repository"
	"clean-arch/database"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateStudentService - admin creates a student
// @Summary Create student
// @Tags Students
// @Description Admin creates a student profile.
// @Accept json
// @Produce json
// @Param body body object true "Student body" example({"studentId":"stu-2025-001","userId":"user-mhs","programStudy":"Teknik Informatika","academicYear":"2025","advisorId":"user-do-1"})
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /students [post]
func CreateStudentService(c *fiber.Ctx) error {
	var s model.Student
	if err := c.BodyParser(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if s.UserID == "" || s.StudentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "student_id and user_id required"})
	}

	if err := repository.CreateStudent(context.Background(), &s); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"studentId": s.ID})
}

// GetStudentService
// @Summary Get student by id
// @Tags Students
// @Description Get student by id.
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} model.Student
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /students/{id} [get]
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

// ListStudentsByAdvisorService
// @Summary List students by advisor
// @Tags Students
// @Description List students under an advisor.
// @Produce json
// @Param advisorId path string true "Advisor ID"
// @Success 200 {array} model.Student
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /students/advisor/{advisorId} [get]
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

// GetStudentAchievementsService - GET /students/:id/achievements
// Mengambil semua mongo achievement berdasarkan reference di Postgres
func GetStudentAchievementsService(c *fiber.Ctx) error {
	studentID := c.Params("id")
	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "student id required"})
	}
	q := `SELECT mongo_achievement_id FROM achievement_references WHERE student_id=$1`
	rows, err := database.PostgresDB.QueryContext(context.Background(), q, studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var mid sql.NullString
		if err := rows.Scan(&mid); err == nil && mid.Valid {
			ids = append(ids, mid.String)
		}
	}
	if len(ids) == 0 {
		return c.JSON([]interface{}{})
	}
	// fetch from mongo
	col := database.MongoDB.Collection("achievements")
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer cur.Close(context.Background())
	var out []map[string]interface{}
	if err := cur.All(context.Background(), &out); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(out)
}

// SetStudentAdvisorService - PUT /students/:id/advisor
func SetStudentAdvisorService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	var body struct {
		AdvisorID string `json:"advisorId"`
	}
	if err := c.BodyParser(&body); err != nil || body.AdvisorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "advisorId required"})
	}
	// try repository update function if present
	if err := repository.UpdateStudentAdvisor(context.Background(), id, body.AdvisorID); err != nil {
		// fallback to direct exec
		_, err2 := database.PostgresDB.ExecContext(context.Background(), `UPDATE students SET advisor_id=$1 WHERE id=$2`, body.AdvisorID, id)
		if err2 != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err2.Error()})
		}
	}
	return c.JSON(fiber.Map{"message": "advisor updated"})
}
