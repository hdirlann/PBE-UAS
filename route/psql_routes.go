package route

import (
	"clean-arch/middleware"
	pgSrv "clean-arch/app/service"
	"github.com/gofiber/fiber/v2"
)

// RegisterPsqlRoutes registers routes that use Postgres-backed services.
// Pastikan RegisterPsqlRoutes dipanggil dari main.go dengan fiber app instance.
func RegisterPsqlRoutes(app *fiber.App) {
	// public routes
	public := app.Group("/api/v1")
	public.Post("/auth/login", pgSrv.AuthenticateService)

	// protected group (JWT required)
	protected := app.Group("/api/v1", middleware.JWTMiddleware())

	// Users (example permissions: users.create, users.list, users.view, users.delete)
	protected.Post("/users", middleware.RequirePermission("users.create"), pgSrv.CreateUserService)
	protected.Get("/users", middleware.RequirePermission("users.list"), pgSrv.ListUsersService)
	protected.Get("/users/:id", middleware.RequirePermission("users.view"), pgSrv.GetUserByIDService)
	protected.Delete("/users/:id", middleware.RequirePermission("users.delete"), pgSrv.DeleteUserService)

	// Roles (example)
	protected.Post("/roles", middleware.RequirePermission("roles.create"), pgSrv.CreateRoleService)
	protected.Get("/roles/:name", middleware.RequirePermission("roles.view"), pgSrv.GetRoleByNameService)

	// Students & Lecturers & Achievement refs - contoh
	protected.Post("/students", middleware.RequirePermission("students.create"), pgSrv.CreateStudentService)
	protected.Get("/students/:id", middleware.RequirePermission("students.view"), pgSrv.GetStudentService)
	protected.Get("/students/advisor/:advisorId", middleware.RequirePermission("students.view"), pgSrv.ListStudentsByAdvisorService)

	protected.Post("/lecturers", middleware.RequirePermission("lecturers.create"), pgSrv.CreateLecturerService)
	protected.Get("/lecturers/:id", middleware.RequirePermission("lecturers.view"), pgSrv.GetLecturerService)

	// Achievement reference workflow
	protected.Post("/refs", middleware.RequirePermission("refs.create"), pgSrv.CreateAchievementReferenceService)
	protected.Post("/refs/:id/submit", middleware.RequirePermission("refs.submit"), pgSrv.SubmitAchievementReferenceService)
	protected.Post("/refs/:id/verify", middleware.RequirePermission("refs.verify"), pgSrv.VerifyAchievementReferenceService)
	protected.Post("/refs/:id/reject", middleware.RequirePermission("refs.reject"), pgSrv.RejectAchievementReferenceService)
}
