package mongo

import (
	"fmt"
	"strconv"
	"time"

	mongoModel "clean-arch/app/model/mongo"
	"clean-arch/app/repository"

	"github.com/gofiber/fiber/v2"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateAchievementService handles POST /api/v1/achievements
// body => JSON sesuai model mongo.Achievement (studentId, achievementType, title, description, details, tags, points)
func CreateAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	var req mongoModel.Achievement
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	created, err := repository.CreateAchievement(db, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Optionally: return created doc
	return c.Status(fiber.StatusCreated).JSON(created)
}

// GetAchievementService handles GET /api/v1/achievements/:id
func GetAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	a, err := repository.GetAchievementByID(db, id)
	if err != nil {
		// if not found, repository returns mongo.ErrNoDocuments
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(a)
}

// UpdateAchievementService handles PUT /api/v1/achievements/:id
func UpdateAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	var update map[string]interface{}
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Prevent updating immutable fields (optionally)
	delete(update, "_id")
	delete(update, "createdAt")

	if err := repository.UpdateAchievement(db, id, bson.M(update)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "updated"})
}

// DeleteAchievementService (soft delete) handles DELETE /api/v1/achievements/:id
func DeleteAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	if err := repository.SoftDeleteAchievement(db, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "deleted"})
}

// HardDeleteAchievementService (admin only) handles DELETE /api/v1/achievements/:id/permanent
func HardDeleteAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	if err := repository.HardDeleteAchievement(db, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "removed"})
}

// ListAchievementsService handles GET /api/v1/achievements
// query: page, limit, studentId, type, search (simple)
func ListAchievementsService(c *fiber.Ctx, db *mgo.Database) error {
	pageQ := c.Query("page", "1")
	limitQ := c.Query("limit", "10")
	studentID := c.Query("studentId", "")
	atype := c.Query("type", "")
	search := c.Query("search", "")

	page, _ := strconv.ParseInt(pageQ, 10, 64)
	limit, _ := strconv.ParseInt(limitQ, 10, 64)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filter := bson.M{}
	if studentID != "" {
		filter["studentId"] = studentID
	}
	if atype != "" {
		filter["achievementType"] = atype
	}
	if search != "" {
		// simple text search on title and description
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	out, total, err := repository.ListAchievements(db, filter, int(page), int(limit))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	meta := fiber.Map{
		"page":  page,
		"limit": limit,
		"total": total,
	}
	return c.JSON(fiber.Map{"data": out, "meta": meta})
}
