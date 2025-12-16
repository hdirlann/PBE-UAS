package service

import (
	"context"
	"time"

	"clean-arch/app/model"
	"clean-arch/app/repository"
	"clean-arch/middleware"
	"github.com/gofiber/fiber/v2"
)

// helper: ambil string dari map dengan beberapa alias
func getStringFromMap(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// CreateAchievementReferenceService
// @Summary Create achievement reference (Postgres ref -> Mongo doc id)
// @Tags AchievementReferences
// @Description Create a new achievement reference in Postgres that references a Mongo achievement.
// @Accept json
// @Produce json
// @Param body body object true "Reference body" example({"studentId":"stud-1","mongoAchievementId":"mongo-abc-123"})
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /refs [post]
func CreateAchievementReferenceService(c *fiber.Ctx) error {
	var body struct {
		StudentID          string `json:"studentId"`
		MongoAchievementID string `json:"mongoAchievementId"`
	}
	_ = c.BodyParser(&body)

	if body.StudentID == "" || body.MongoAchievementID == "" {
		var m map[string]interface{}
		if err := c.BodyParser(&m); err == nil {
			if body.StudentID == "" {
				body.StudentID = getStringFromMap(m, "studentId", "student_id", "student")
			}
			if body.MongoAchievementID == "" {
				body.MongoAchievementID = getStringFromMap(m, "mongoAchievementId", "mongo_achievement_id", "mongoId", "mongo_id")
			}
		}
	}

	if body.StudentID == "" || body.MongoAchievementID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "studentId and mongoAchievementId required"})
	}

	now := time.Now()
	ref := &model.AchievementReference{
		StudentID:          body.StudentID,
		MongoAchievementID: body.MongoAchievementID,
		Status:             "draft",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := repository.CreateAchievementReference(context.Background(), ref); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"referenceId": ref.ID,
		"status":      ref.Status,
		"studentId":   ref.StudentID,
		"mongoId":     ref.MongoAchievementID,
	})
}

// SubmitAchievementReferenceService
// @Summary Submit a draft reference
// @Tags AchievementReferences
// @Description Student submits a draft reference for verification.
// @Param id path string true "Reference ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /refs/{id}/submit [post]
func SubmitAchievementReferenceService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	if err := repository.UpdateAchievementReferenceStatus(context.Background(), id, "submitted", nil, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"id": id, "status": "submitted"})
}

// VerifyAchievementReferenceService
// @Summary Verify a reference
// @Tags AchievementReferences
// @Description Lecturer/admin verifies a submitted reference. Verifier taken from JWT.
// @Param id path string true "Reference ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /refs/{id}/verify [post]
func VerifyAchievementReferenceService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}

	verifierID, _ := c.Locals(middleware.LocalsUserID).(string)
	if verifierID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "verifier id missing in token"})
	}
	if err := repository.UpdateAchievementReferenceStatus(context.Background(), id, "verified", &verifierID, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"id": id, "status": "verified"})
}

// RejectAchievementReferenceService
// @Summary Reject a reference
// @Tags AchievementReferences
// @Description Lecturer/admin rejects a submitted reference with a note. Verifier taken from JWT.
// @Param id path string true "Reference ID"
// @Param body body object true "Reject body" example({"note":"dokumen kurang"})
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /refs/{id}/reject [post]
func RejectAchievementReferenceService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	var body struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if body.Note == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "note required"})
	}
	verifierID, _ := c.Locals(middleware.LocalsUserID).(string)
	if verifierID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "verifier id missing in token"})
	}
	if err := repository.UpdateAchievementReferenceStatus(context.Background(), id, "rejected", &verifierID, &body.Note); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"id": id, "status": "rejected"})
}

