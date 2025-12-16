package service

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	mongoModel "clean-arch/app/model"
	repo "clean-arch/app/repository" // alias untuk package repository
	"clean-arch/database"
	"clean-arch/middleware"

	"github.com/gofiber/fiber/v2"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateAchievementService handles POST /api/v1/achievements
// @Summary Create achievement (Mongo)
// @Tags Achievements
// @Description Create an achievement document in MongoDB.
// @Accept json
// @Produce json
// @Param body body mongoModel.Achievement true "Achievement body"
// @Success 201 {object} mongoModel.Achievement
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /achievements [post]
func CreateAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	var req mongoModel.Achievement
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// ambil user dari JWT
	userID, _ := c.Locals(middleware.LocalsUserID).(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthenticated"})
	}

	// ambil student berdasarkan user_id
	student, err := repo.GetStudentByUserID(context.Background(), userID)
	if err != nil || student == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "student profile not found",
		})
	}

	// paksa studentId dari JWT
	req.StudentID = student.ID

	// 1️⃣ simpan ke MongoDB
	created, err := repo.CreateAchievement(db, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 2️⃣ BUAT reference PostgreSQL (status = draft)
	now := time.Now()
	ref := &mongoModel.AchievementReference{
		StudentID:          student.ID,
		MongoAchievementID: created.ID.Hex(),
		Status:             "draft",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := repo.CreateAchievementReference(context.Background(), ref); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 3️⃣ response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"achievement": created,
		"reference": fiber.Map{
			"id":     ref.ID,
			"status": ref.Status,
		},
	})
}


// GetAchievementService handles GET /api/v1/achievements/:id
// @Summary Get achievement by id
// @Tags Achievements
// @Description Get achievement document from Mongo by id.
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} mongoModel.Achievement
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /achievements/{id} [get]
func GetAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	a, err := repo.GetAchievementByID(db, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(a)
}

// UpdateAchievementService handles PUT /api/v1/achievements/:id
// @Summary Update achievement
// @Tags Achievements
// @Description Update an achievement document by id.
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param body body object true "Update fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /achievements/{id} [put]
func UpdateAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	var update map[string]interface{}
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	delete(update, "_id")
	delete(update, "createdAt")

	if err := repo.UpdateAchievement(db, id, bson.M(update)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "updated"})
}

// DeleteAchievementService (soft delete - mahasiswa)
// @Summary Delete draft achievement
// @Tags Achievements
// @Description Mahasiswa menghapus prestasi draft
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /achievements/{id} [delete]
func DeleteAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}

	userID, _ := c.Locals(middleware.LocalsUserID).(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthenticated"})
	}

	// ambil student
	student, err := repo.GetStudentByUserID(context.Background(), userID)
	if err != nil || student == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "student not found"})
	}

	// ambil reference
	ref, err := repo.GetAchievementReferenceByMongoID(context.Background(), mongoID)
	if err != nil || ref == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "reference not found"})
	}

	// ownership check
	if ref.StudentID != student.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not owner"})
	}

	// status harus draft
	if ref.Status != "draft" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "only draft achievement can be deleted",
		})
	}

	// soft delete mongo
	if err := repo.SoftDeleteAchievement(db, mongoID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// update postgres reference
	if err := repo.UpdateAchievementReferenceStatus(
		context.Background(),
		ref.ID,
		"deleted",
		nil,
		nil,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "achievement deleted"})
}

// HardDeleteAchievementService (admin)
// @Summary Hard delete achievement
// @Tags Achievements
// @Description Permanently remove achievement document (admin).
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /achievements/{id}/permanent [delete]
func HardDeleteAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	if err := repo.HardDeleteAchievement(db, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "removed"})
}

// ListAchievementsService
// @Summary List achievements
// @Tags Achievements
// @Description List achievements with pagination and filters.
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param studentId query string false "Student ID to filter"
// @Param type query string false "Achievement type"
// @Param search query string false "Search text"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /achievements [get]
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
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	out, total, err := repo.ListAchievements(db, filter, page, limit)
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

// -------------------- SRS workflow helpers (submit / verify / reject / history) --------------------

// SubmitAchievementService handles POST /achievements/:id/submit
// Flow: student submits a mongo achievement for verification -> create or update postgres reference
func SubmitAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}

	// ambil user id dari JWT
	userID, _ := c.Locals(middleware.LocalsUserID).(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthenticated"})
	}

	// 🔥 PENTING: ambil STUDENT berdasarkan user_id
	student, err := repo.GetStudentByUserID(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if student == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "student profile not found",
		})
	}

	// cari reference berdasarkan mongo achievement
	ref, err := repo.GetAchievementReferenceByMongoID(context.Background(), mongoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if ref == nil {
		// create new reference
		now := time.Now()
		newRef := &mongoModel.AchievementReference{
			StudentID:          student.ID, // ✅ INI YANG BENAR
			MongoAchievementID: mongoID,
			Status:             "submitted",
			SubmittedAt:        &now,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		if err := repo.CreateAchievementReference(context.Background(), newRef); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message":      "submitted",
			"referenceId":  newRef.ID,
			"studentId":    student.ID,
		})
	}

	// ownership check
	if ref.StudentID != student.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not owner"})
	}

	// update status
	if err := repo.UpdateAchievementReferenceStatus(
		context.Background(),
		ref.ID,
		"submitted",
		nil,
		nil,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":     "submitted",
		"referenceId": ref.ID,
	})
}


// VerifyAchievementService handles POST /achievements/:id/verify
// Flow: lecturer verifies a submitted reference (by mongo id)
func VerifyAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	verifierID, _ := c.Locals(middleware.LocalsUserID).(string)
	if verifierID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthenticated"})
	}

	ref, err := repo.GetAchievementReferenceByMongoID(context.Background(), mongoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if ref == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "reference not found"})
	}
	if err := repo.UpdateAchievementReferenceStatus(context.Background(), ref.ID, "verified", &verifierID, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "verified", "referenceId": ref.ID})
}

// RejectAchievementService handles POST /achievements/:id/reject
// Flow: lecturer rejects with a note
func RejectAchievementService(c *fiber.Ctx, db *mgo.Database) error {
	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}
	var body struct{ Note string `json:"note"` }
	if err := c.BodyParser(&body); err != nil || body.Note == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "note required"})
	}
	verifierID, _ := c.Locals(middleware.LocalsUserID).(string)
	if verifierID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthenticated"})
	}
	ref, err := repo.GetAchievementReferenceByMongoID(context.Background(), mongoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if ref == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "reference not found"})
	}
	if err := repo.UpdateAchievementReferenceStatus(context.Background(), ref.ID, "rejected", &verifierID, &body.Note); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "rejected", "referenceId": ref.ID})
}

// GetAchievementHistoryService handles GET /achievements/:id/history
func GetAchievementHistoryService(c *fiber.Ctx, db *mgo.Database) error {
	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}

	q := `SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
	      FROM achievement_references WHERE mongo_achievement_id=$1 ORDER BY created_at DESC LIMIT 20`
	rows, err := database.PostgresDB.QueryContext(context.Background(), q, mongoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var out []map[string]interface{}
	for rows.Next() {
		var id, studentID, mongoIDv sql.NullString
		var status sql.NullString
		var submitted, verified sql.NullTime
		var verifiedBy, rejectionNote sql.NullString
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &studentID, &mongoIDv, &status, &submitted, &verified, &verifiedBy, &rejectionNote, &createdAt, &updatedAt); err != nil {
			continue
		}
		m := map[string]interface{}{
			"id":             id.String,
			"studentId":      studentID.String,
			"mongoId":        mongoIDv.String,
			"status":         status.String,
			"submitted_at":   nil,
			"verified_at":    nil,
			"verified_by":    nil,
			"rejection_note": nil,
			"created_at":     nil,
			"updated_at":     nil,
		}
		if submitted.Valid {
			m["submitted_at"] = submitted.Time
		}
		if verified.Valid {
			m["verified_at"] = verified.Time
		}
		if verifiedBy.Valid {
			m["verified_by"] = verifiedBy.String
		}
		if rejectionNote.Valid {
			m["rejection_note"] = rejectionNote.String
		}
		if createdAt.Valid {
			m["created_at"] = createdAt.Time
		}
		if updatedAt.Valid {
			m["updated_at"] = updatedAt.Time
		}
		out = append(out, m)
	}
	return c.JSON(out)
}

// -------------------- Reporting & Analytics --------------------

// StatisticsService - basic stats (counts)
// Returns total achievements (mongo) and total references (postgres).
func StatisticsService(c *fiber.Ctx) error {
	ctx := context.Background()

	// Mongo: count achievements collection
	var mongoCount int64 = 0
	if database.MongoDB != nil {
		col := database.MongoDB.Collection("achievements")
		cnt, err := col.CountDocuments(ctx, bson.M{})
		if err == nil {
			mongoCount = cnt
		}
	}

	// Postgres: count achievement_references
	var pgCount int64 = 0
	if database.PostgresDB != nil {
		row := database.PostgresDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM achievement_references`)
		_ = row.Scan(&pgCount) // ignore error and return 0 if fails
	}

	out := fiber.Map{
		"total_achievements_mongo": mongoCount,
		"total_references_pg":      pgCount,
		"generated_at":             time.Now(),
	}
	return c.JSON(out)
}

// StudentReportService - basic per-student report
// Returns student profile and count of their references and achievements.
func StudentReportService(c *fiber.Ctx) error {
	studentID := c.Params("id")
	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "student id required"})
	}
	ctx := context.Background()

	// Get student profile (Postgres)
	st, err := repo.GetStudentByID(ctx, studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if st == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}

	// Count refs in Postgres
	var refsCount int64 = 0
	if database.PostgresDB != nil {
		row := database.PostgresDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM achievement_references WHERE student_id=$1`, studentID)
		_ = row.Scan(&refsCount)
	}

	// Count referenced mongo achievements for this student
	var mongoAchievementsCount int64 = 0
	if database.PostgresDB != nil && database.MongoDB != nil {
		rows, err := database.PostgresDB.QueryContext(ctx, `SELECT mongo_achievement_id FROM achievement_references WHERE student_id=$1`, studentID)
		if err == nil {
			defer rows.Close()
			var ids []string
			for rows.Next() {
				var mid sql.NullString
				if err := rows.Scan(&mid); err == nil && mid.Valid {
					ids = append(ids, mid.String)
				}
			}
			if len(ids) > 0 {
				col := database.MongoDB.Collection("achievements")
				cnt, err := col.CountDocuments(ctx, bson.M{"_id": bson.M{"$in": ids}})
				if err == nil {
					mongoAchievementsCount = cnt
				}
			}
		}
	}

	resp := fiber.Map{
		"student":                  st,
		"references_count":         refsCount,
		"mongo_achievements_count": mongoAchievementsCount,
		"generated_at":             time.Now(),
	}
	return c.JSON(resp)
}
