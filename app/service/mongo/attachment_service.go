package service

import (
	"clean-arch/app/model/mongo"
	"clean-arch/app/repository/mongo"

	"github.com/gofiber/fiber/v2"
	mgo "go.mongodb.org/mongo-driver/mongo"
)

// AddAttachmentService handles POST /api/v1/achievements/:id/attachments
// expects JSON body { fileName, fileUrl, fileType } OR you can adapt for multipart upload.
func AddAttachmentService(c *fiber.Ctx, db *mgo.Database) error {
	achID := c.Params("id")
	if achID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "achievement id required"})
	}

	var body struct {
		FileName string `json:"fileName"`
		FileURL  string `json:"fileUrl"`
		FileType string `json:"fileType"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	attach := &mongo.Attachment{
		AchievementID: achID,
		FileName:      body.FileName,
		FileURL:       body.FileURL,
		FileType:      body.FileType,
	}

	res, err := repository.AddAttachment(db, attach)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// ListAttachmentsService handles GET /api/v1/achievements/:id/attachments
func ListAttachmentsService(c *fiber.Ctx, db *mgo.Database) error {
	achID := c.Params("id")
	if achID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "achievement id required"})
	}
	out, err := repository.ListAttachmentsByAchievement(db, achID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(out)
}
