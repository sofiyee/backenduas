package routes

import (
	"backenduas/app/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	authService *service.AuthService,
	userService *service.UserService,
	studentService *service.StudentService,
	lecturerService *service.LecturerService,
	achievementService *service.AchievementService, // ⬅️ TAMBAH INI
) {
	fmt.Println("🔥 REGISTERING ROUTES...")

	api := app.Group("/api/v1")

	// Auth Routes
	AuthRoutes(api, authService)

	// User CRUD Routes
	UserRoutes(api, userService)

	// Student Routes
	StudentRoutes(api, studentService)

	// Lecturer Routes
	LecturerRoutes(api, lecturerService)

	// Achievement Routes (NEW)
	AchievementRoutes(api, achievementService) 

	fmt.Println("🔥 ROUTES REGISTERED")
}
