package route

import (
	"clean-arch/database"
	"clean-arch/middleware"
	svc "clean-arch/app/service"

	"github.com/gofiber/fiber/v2"
)

// RegisterAPIRoutes registers all API endpoints according to the SRS.
func RegisterAPIRoutes(app *fiber.App) {
	// Public group (no JWT) - only auth login & refresh
	public := app.Group("/api/v1")
	public.Post("/auth/login", svc.AuthenticateService)
	public.Post("/auth/refresh", svc.RefreshTokenService) // optional: allow token refresh without middleware

	// Protected group (JWT required)
	protected := app.Group("/api/v1", middleware.JWTMiddleware())

	// ----------------------
	// Auth (protected endpoints)
	// ----------------------
	protected.Post("/auth/logout", svc.LogoutService)
	protected.Get("/auth/profile", svc.ProfileService)

	// ----------------------
	// Users (Admin)
	// ----------------------
	// CRUD + assign role
	protected.Get("/users", middleware.RequirePermission("users.list"), svc.ListUsersService)
	protected.Get("/users/:id", middleware.RequirePermission("users.view"), svc.GetUserByIDService)
	protected.Post("/users", middleware.RequirePermission("users.create"), svc.CreateUserService)
	protected.Put("/users/:id", middleware.RequirePermission("users.update"), svc.UpdateUserService)
	protected.Delete("/users/:id", middleware.RequirePermission("users.delete"), svc.DeleteUserService)
	protected.Put("/users/:id/role", middleware.RequirePermission("users.assign_role"), svc.AssignRoleService)

	// ----------------------
	// Roles
	// ----------------------
	protected.Post("/roles", middleware.RequirePermission("roles.create"), svc.CreateRoleService)
	protected.Get("/roles/:name", middleware.RequirePermission("roles.view"), svc.GetRoleByNameService)

	// ----------------------
	// Achievements (Mongo) + reference workflow (Postgres)
	// - Many handlers accept database.MongoDB (mongo.Database) parameter via closure
	// ----------------------
	protected.Get("/achievements", middleware.RequirePermission("achievements.list"), func(c *fiber.Ctx) error {
		return svc.ListAchievementsService(c, database.MongoDB)
	})
	protected.Get("/achievements/:id", middleware.RequirePermission("achievements.view"), func(c *fiber.Ctx) error {
		return svc.GetAchievementService(c, database.MongoDB)
	})
	protected.Post("/achievements", middleware.RequirePermission("achievements.create"), func(c *fiber.Ctx) error {
		return svc.CreateAchievementService(c, database.MongoDB)
	})
	protected.Put("/achievements/:id", middleware.RequirePermission("achievements.update"), func(c *fiber.Ctx) error {
		return svc.UpdateAchievementService(c, database.MongoDB)
	})
	protected.Delete("/achievements/:id", middleware.RequirePermission("achievements.delete"), func(c *fiber.Ctx) error {
		return svc.DeleteAchievementService(c, database.MongoDB)
	})
	// Hard delete (admin)
	protected.Delete("/achievements/:id/permanent", middleware.RequirePermission("achievements.hard_delete"), func(c *fiber.Ctx) error {
		return svc.HardDeleteAchievementService(c, database.MongoDB)
	})

	// Submit / verify / reject flows (these operate by linking mongo doc -> postgres reference)
	protected.Post("/achievements/:id/submit", middleware.RequirePermission("achievements.submit"), func(c *fiber.Ctx) error {
		return svc.SubmitAchievementService(c, database.MongoDB)
	})
	protected.Post("/achievements/:id/verify", middleware.RequirePermission("achievements.verify"), func(c *fiber.Ctx) error {
		return svc.VerifyAchievementService(c, database.MongoDB)
	})
	protected.Post("/achievements/:id/reject", middleware.RequirePermission("achievements.reject"), func(c *fiber.Ctx) error {
		return svc.RejectAchievementService(c, database.MongoDB)
	})

	// Status history (reads from Postgres references)
	protected.Get("/achievements/:id/history", middleware.RequirePermission("achievements.history"), func(c *fiber.Ctx) error {
		return svc.GetAchievementHistoryService(c, database.MongoDB)
	})

	// Attachments upload & list (mongo-backed attachments collection or GridFS)
	protected.Post("/achievements/:id/attachments", middleware.RequirePermission("achievements.upload_attachment"), func(c *fiber.Ctx) error {
		return svc.AddAttachmentService(c, database.MongoDB)
	})
	protected.Get("/achievements/:id/attachments", middleware.RequirePermission("achievements.view_attachments"), func(c *fiber.Ctx) error {
		return svc.ListAttachmentsService(c, database.MongoDB)
	})

	// ----------------------
	// Achievement References (Postgres) - alternate entry (if needed)
	// ----------------------
	protected.Post("/refs", middleware.RequirePermission("refs.create"), svc.CreateAchievementReferenceService)
	protected.Post("/refs/:id/submit", middleware.RequirePermission("refs.submit"), svc.SubmitAchievementReferenceService)
	protected.Post("/refs/:id/verify", middleware.RequirePermission("refs.verify"), svc.VerifyAchievementReferenceService)
	protected.Post("/refs/:id/reject", middleware.RequirePermission("refs.reject"), svc.RejectAchievementReferenceService)

	// ----------------------
	// Students & Lecturers
	// ----------------------
	// Students
	protected.Post("/students", middleware.RequirePermission("students.create"), svc.CreateStudentService)
	protected.Get("/students", middleware.RequirePermission("students.list"), svc.ListStudentsByAdvisorService) // NOTE: this route could be adapted to list all or by query
	protected.Get("/students/:id", middleware.RequirePermission("students.view"), svc.GetStudentService)
	protected.Get("/students/:id/achievements", middleware.RequirePermission("students.read_achievements"), svc.GetStudentAchievementsService)
	protected.Put("/students/:id/advisor", middleware.RequirePermission("students.set_advisor"), svc.SetStudentAdvisorService)

	// Lecturers
	protected.Post("/lecturers", middleware.RequirePermission("lecturers.create"), svc.CreateLecturerService)
	protected.Get("/lecturers", middleware.RequirePermission("lecturers.list"), svc.ListLecturersService)
	protected.Get("/lecturers/:id", middleware.RequirePermission("lecturers.view"), svc.GetLecturerService)
	protected.Get("/lecturers/:id/advisees", middleware.RequirePermission("lecturers.view_advisees"), svc.GetLecturerAdviseesService)

	// ----------------------
	// Reports & Analytics
	// ----------------------
	protected.Get("/reports/statistics", middleware.RequirePermission("reports.read"), svc.StatisticsService)
	protected.Get("/reports/student/:id", middleware.RequirePermission("reports.read"), svc.StudentReportService)
}
