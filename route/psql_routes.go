package route

import (
	"github.com/gofiber/fiber/v2"

	pgSrv "clean-arch/app/service/postgre"
)

// RegisterPsqlRoutes registers routes that use PostgreSQL (users, roles, students, lecturers, refs)
func RegisterPsqlRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	// Auth & Users
	v1.Post("/auth/login", pgSrv.AuthenticateService)
	v1.Post("/users", pgSrv.CreateUserService)
	// Roles
	v1.Post("/roles", pgSrv.CreateRoleService)
	v1.Get("/roles/:name", pgSrv.GetRoleByNameService)

	// Students
	v1.Post("/students", pgSrv.CreateStudentService)
	v1.Get("/students/:id", pgSrv.GetStudentService)
	v1.Get("/students/advisor/:advisorId", pgSrv.ListStudentsByAdvisorService)

	// Lecturers
	v1.Post("/lecturers", pgSrv.CreateLecturerService)
	v1.Get("/lecturers/:id", pgSrv.GetLecturerService)

	// Achievement references (workflow)
	v1.Post("/refs", pgSrv.CreateAchievementReferenceService)
	v1.Post("/refs/:id/submit", pgSrv.SubmitAchievementReferenceService)
	v1.Post("/refs/:id/verify", pgSrv.VerifyAchievementReferenceService)
	v1.Post("/refs/:id/reject", pgSrv.RejectAchievementReferenceService)
}
