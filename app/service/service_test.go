package service

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// ==========================
// STUDENT REPORT
// ==========================

func TestStudentReport_RouteNotMatch(t *testing.T) {
	app := fiber.New()

	app.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		return StudentReportService(c)
	})

	req := httptest.NewRequest("GET", "/reports/student/", nil)
	resp, _ := app.Test(req)

	// Router behavior
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

// ==========================
// STATISTICS
// ==========================

func TestStatistics_WithoutMiddleware(t *testing.T) {
	app := fiber.New()

	app.Get("/reports/statistics", func(c *fiber.Ctx) error {
		return StatisticsService(c)
	})

	req := httptest.NewRequest("GET", "/reports/statistics", nil)
	resp, _ := app.Test(req)

	// Tanpa JWT middleware → service tetap jalan
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

// ==========================
// REJECT ACHIEVEMENT
// ==========================

func TestRejectAchievement_RouteNotMatch(t *testing.T) {
	app := fiber.New()

	app.Post("/achievements/:id/reject", func(c *fiber.Ctx) error {
		return RejectAchievementService(c, nil)
	})

	req := httptest.NewRequest("POST", "/achievements//reject", nil)
	resp, _ := app.Test(req)

	// Router tidak match
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestRejectAchievement_MissingNote(t *testing.T) {
	app := fiber.New()

	app.Post("/achievements/:id/reject", func(c *fiber.Ctx) error {
		return RejectAchievementService(c, nil)
	})

	req := httptest.NewRequest("POST", "/achievements/abc/reject", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// ==========================
// STUDENT CREATE
// ==========================

func TestCreateStudent_InvalidBody(t *testing.T) {
	app := fiber.New()

	app.Post("/students", func(c *fiber.Ctx) error {
		return CreateStudentService(c)
	})

	req := httptest.NewRequest("POST", "/students", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
