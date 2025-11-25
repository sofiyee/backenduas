package routes

import (
	"backenduas/app/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authService *service.AuthService, userService *service.UserService, studentService *service.StudentService, lecturerService *service.LecturerService) {
	fmt.Println("🔥 REGISTERING ROUTES...")

	api := app.Group("/api/v1")

	// Auth Routes
	AuthRoutes(api, authService)

	// User CRUD Routes
	UserRoutes(api, userService)

	StudentRoutes(api, studentService)

	LecturerRoutes(api, lecturerService)

	fmt.Println("🔥 ROUTES REGISTERED")
}
