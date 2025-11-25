package routes

import (
	"backenduas/app/service"
	"backenduas/middleware"
	"github.com/gofiber/fiber/v2"
)

func LecturerRoutes(api fiber.Router, lecturerService *service.LecturerService) {
	lect := api.Group("/lecturers", middleware.JWTProtected())

	// Admin & dosen wali bisa lihat ini
	lect.Get("/", middleware.AllowRoles("Admin", "Dosen Wali"), lecturerService.GetAll)

	// List mahasiswa bimbingan → hanya dosen wali yg bersangkutan + admin
	lect.Get("/:id/advisees", middleware.AllowRoles("Admin", "Dosen Wali"), lecturerService.GetAdvisees)
}
