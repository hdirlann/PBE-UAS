package postgre

import (
	"context"
	"time"

	"clean-arch/app/model/postgre"
	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
)

func CreateAchievementReferenceService(c *fiber.Ctx) error {
	var body struct {
		StudentID string `json:"studentId"`
		MongoID   string `json:"mongoId"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if body.StudentID == "" || body.MongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"studentId and mongoId required"})
	}
	ref := &postgre.AchievementReference{
		StudentID:          body.StudentID,
		MongoAchievementID: body.MongoID,
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := repository.CreateAchievementReference(context.Background(), ref); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"referenceId": ref.ID})
}

func SubmitAchievementReferenceService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"id required"}) }
	if err := repository.UpdateAchievementReferenceStatus(context.Background(), id, "submitted", nil, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status":"submitted"})
}

func VerifyAchievementReferenceService(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct{ VerifierID string `json:"verifierId"` }
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if body.VerifierID == "" { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"verifierId required"}) }
	if err := repository.UpdateAchievementReferenceStatus(context.Background(), id, "verified", &body.VerifierID, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status":"verified"})
}

func RejectAchievementReferenceService(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		VerifierID string `json:"verifierId"`
		Note       string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if body.VerifierID == "" || body.Note == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"verifierId and note required"})
	}
	if err := repository.UpdateAchievementReferenceStatus(context.Background(), id, "rejected", &body.VerifierID, &body.Note); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status":"rejected"})
}
