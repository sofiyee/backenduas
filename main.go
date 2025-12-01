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

	// 2. Connect PostgreSQL + MongoDB
	database.ConnectDatabases()

	// 3. Setup Fiber
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	// === Init repository ===
	authRepo := repository.NewAuthRepository()
	userRepo := repository.NewUserRepository()
	studentRepo := repository.NewStudentRepository()
	lecturerRepo := repository.NewLecturerRepository()

	// Achievement repo
	pgAchRepo := repository.NewAchievementPGRepository()
	mongoAchRepo := repository.NewAchievementMongoRepository()

	// === Init services ===
	authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)
	studentService := service.NewStudentService(studentRepo)
	lecturerService := service.NewLecturerService(lecturerRepo)

	// Achievement service membutuhkan 3 repo: PG, Mongo, Student
	achievementService := service.NewAchievementService(pgAchRepo, mongoAchRepo, studentRepo)

	// 4. Setup routes
	routes.SetupRoutes(
		app,
		authService,
		userService,
		studentService,
		lecturerService,
		achievementService,
	)

	// 5. Start server
	port := ":" + config.AppEnv.AppPort
	log.Println("🚀 Server running on port", port)
	app.Listen(port)
}
