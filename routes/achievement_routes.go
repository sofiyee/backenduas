package routes

import (
	"backenduas/app/service"
	"backenduas/middleware"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(api fiber.Router, ach *service.AchievementService) {
	r := api.Group("/achievements", middleware.JWTProtected())

	// LIST (filtered by role)
	r.Get("/", middleware.AllowRoles("Mahasiswa", "Dosen Wali", "Admin"), ach.GetAll)

	// DETAIL
	r.Get("/:id", middleware.AllowRoles("Mahasiswa", "Dosen Wali", "Admin"), ach.GetByID)

	// CREATE (Mahasiswa)
	r.Post("/", middleware.AllowRoles("Mahasiswa"), ach.Create)

	// UPDATE (Mahasiswa)
	r.Put("/:id", middleware.AllowRoles("Mahasiswa"), ach.Update)

	// DELETE (Mahasiswa)
	r.Delete("/:id", middleware.AllowRoles("Mahasiswa"), ach.Delete)

	// SUBMIT (Mahasiswa)
	r.Post("/:id/submit", middleware.AllowRoles("Mahasiswa"), ach.Submit)

	// VERIFY (Dosen Wali)
	r.Post("/:id/verify", middleware.AllowRoles("Dosen Wali"), ach.Verify)

	// REJECT (Dosen Wali)
	r.Post("/:id/reject", middleware.AllowRoles("Dosen Wali"), ach.Reject)

	// HISTORY (Semua role)
	r.Get("/:id/history", middleware.AllowRoles("Mahasiswa", "Dosen Wali", "Admin"), ach.History)

	// UPLOAD ATTACHMENT (Mahasiswa)
	r.Post("/:id/attachments", middleware.AllowRoles("Mahasiswa"), ach.UploadAttachment)
}
