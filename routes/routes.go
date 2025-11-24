package routes

import (
	"backenduas/app/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authService *service.AuthService, userService *service.UserService) {
	fmt.Println("🔥 REGISTERING ROUTES...")

	api := app.Group("/api/v1")

	// Auth Routes
	AuthRoutes(api, authService)

	// User CRUD Routes
	UserRoutes(api, userService)

	fmt.Println("🔥 ROUTES REGISTERED")
}
