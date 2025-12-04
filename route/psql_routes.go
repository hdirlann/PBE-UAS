package route

import (
	"github.com/gofiber/fiber/v2"

	pgSrv "clean-arch/app/service"
	"clean-arch/middleware"
)

func RegisterPsqlRoutes(app *fiber.App) {
	public := app.Group("/api/v1")
	public.Post("/auth/login", pgSrv.AuthenticateService)

	protected := app.Group("/api/v1", middleware.JWTMiddleware())

	// users (example permission strings)
	protected.Post("/users", middleware.RequirePermission("users.create"), pgSrv.CreateUserService)

	// roles
	protected.Post("/roles", middleware.RequirePermission("roles.create"), pgSrv.CreateRoleService)
	protected.Get("/roles/:name", pgSrv.GetRoleByNameService)

	// students
	protected.Post("/students", middleware.RequirePermission("students.create"), pgSrv.CreateStudentService)
	protected.Get("/students/:id", middleware.RequirePermission("students.view"), pgSrv.GetStudentService)
	protected.Get("/students/advisor/:advisorId", middleware.RequirePermission("students.view"), pgSrv.ListStudentsByAdvisorService)

	// lecturers
	protected.Post("/lecturers", middleware.RequirePermission("lecturers.create"), pgSrv.CreateLecturerService)
	protected.Get("/lecturers/:id", middleware.RequirePermission("lecturers.view"), pgSrv.GetLecturerService)

	// achievement refs workflow
	protected.Post("/refs", middleware.RequirePermission("refs.create"), pgSrv.CreateAchievementReferenceService)
	protected.Post("/refs/:id/submit", middleware.RequirePermission("refs.submit"), pgSrv.SubmitAchievementReferenceService)
	protected.Post("/refs/:id/verify", middleware.RequirePermission("refs.verify"), pgSrv.VerifyAchievementReferenceService)
	protected.Post("/refs/:id/reject", middleware.RequirePermission("refs.reject"), pgSrv.RejectAchievementReferenceService)
}
