package main

import (
	"log"

	"backenduas/config"
	"backenduas/database"
	"backenduas/routes"
	"backenduas/app/repository"
	"backenduas/app/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Load ENV
	config.LoadEnv()

	// 2. Connect to PostgreSQL
	database.ConnectPostgre()

	// 3. Setup Fiber
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	// === Init repository & services ===
	authRepo := repository.NewAuthRepository()
	authService := service.NewAuthService(authRepo)

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)

	studentRepo := repository.NewStudentRepository()
	studentService := service.NewStudentService(studentRepo)

	lecturerRepo := repository.NewLecturerRepository()
	lecturerService := service.NewLecturerService(lecturerRepo)

	// 4. Setup routes (with dependency injection)
	routes.SetupRoutes(app, authService, userService, studentService, lecturerService)

	// 5. Start server
	port := ":" + config.AppEnv.AppPort
	log.Println("🚀 Server running on port", port)
	app.Listen(port)
}
