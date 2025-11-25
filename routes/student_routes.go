// routes/student_routes.go
package routes

import (
	"backenduas/app/service"
	"backenduas/middleware"

	"github.com/gofiber/fiber/v2"
)

func StudentRoutes(api fiber.Router, studentService *service.StudentService) {
	students := api.Group("/students", middleware.JWTProtected())

	// Admin only
	students.Get("/", middleware.AllowRoles("Admin"), studentService.GetAll)
	students.Post("/", middleware.AllowRoles("Admin"), studentService.Create)
	students.Put("/:id/advisor", middleware.AllowRoles("Admin"), studentService.UpdateAdvisor)

	// Admin + Dosen Wali (lihat detail mahasiswa)
	students.Get("/:id", middleware.AllowRoles("Admin", "Dosen Wali"), studentService.GetByID)

	// Nanti tambahkan:
	// students.Get("/:id/achievements", ...)
}
