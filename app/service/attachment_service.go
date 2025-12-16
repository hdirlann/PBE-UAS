package service

import (
	"clean-arch/app/model"
	"clean-arch/app/repository"

	"github.com/gofiber/fiber/v2"
	mgo "go.mongodb.org/mongo-driver/mongo"
)

// AddAttachmentService handles POST /api/v1/achievements/:id/attachments
// @Summary Add attachment to achievement
// @Tags Attachments
// @Description Add an attachment (file meta) to an achievement document.
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param body body object true "Attachment body" example({"fileName":"dok.pdf","fileUrl":"https://...","fileType":"pdf"})
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /achievements/{id}/attachments [post]
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

	attach := &model.Attachment{
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

// ListAttachmentsService
// @Summary List attachments for achievement
// @Tags Attachments
// @Description List attachments for an achievement.
// @Param id path string true "Achievement ID"
// @Produce json
// @Success 200 {array} model.Attachment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /achievements/{id}/attachments [get]
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
