package route

import (
	"clean-arch/database"
	mongoSrv "clean-arch/app/service"

	"github.com/gofiber/fiber/v2"
)

// RegisterMongoRoutes registers routes that use MongoDB (achievements, attachments)
func RegisterMongoRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	// Achievements
	v1.Post("/achievements", func(c *fiber.Ctx) error {
		return mongoSrv.CreateAchievementService(c, database.MongoDB)
	})
	v1.Get("/achievements", func(c *fiber.Ctx) error {
		return mongoSrv.ListAchievementsService(c, database.MongoDB)
	})
	v1.Get("/achievements/:id", func(c *fiber.Ctx) error {
		return mongoSrv.GetAchievementService(c, database.MongoDB)
	})
	v1.Put("/achievements/:id", func(c *fiber.Ctx) error {
		return mongoSrv.UpdateAchievementService(c, database.MongoDB)
	})
	v1.Delete("/achievements/:id", func(c *fiber.Ctx) error {
		return mongoSrv.DeleteAchievementService(c, database.MongoDB)
	})
	v1.Delete("/achievements/:id/permanent", func(c *fiber.Ctx) error {
		return mongoSrv.HardDeleteAchievementService(c, database.MongoDB)
	})

	// Attachments
	v1.Post("/achievements/:id/attachments", func(c *fiber.Ctx) error {
		return mongoSrv.AddAttachmentService(c, database.MongoDB)
	})
	v1.Get("/achievements/:id/attachments", func(c *fiber.Ctx) error {
		return mongoSrv.ListAttachmentsService(c, database.MongoDB)
	})
}
