package routes

import (
	"backenduas/app/service"
	"backenduas/middleware"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(api fiber.Router, svc *service.ReportService) {
	r := api.Group("/reports", middleware.JWTProtected())

	// Global Statistics (Admin & Dosen Wali)
	r.Get("/statistics", 
		middleware.AllowRoles("Admin", "Dosen Wali", "Mahasiswa"), 
		svc.GlobalStatistics)

	// Per Student Statistics
	r.Get("/student/:id",
		middleware.AllowRoles("Admin", "Dosen Wali", "Mahasiswa"),
		svc.StudentStatistics)
}
