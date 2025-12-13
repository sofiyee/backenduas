package main

import (
	"log"

	"backenduas/config"
	"backenduas/database"
	"backenduas/routes"
	"backenduas/app/repository"
	"backenduas/app/service"

	_ "backenduas/docs" 

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title Backend UAS API
// @version 1.0
// @description API Sistem Manajemen Prestasi Mahasiswa
// @termsOfService http://swagger.io/terms/

// @contact.name Sofie Kusuma Anggraini
// @contact.email sofie@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1


// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// 1. Load ENV
	config.LoadEnv()

	// 2. Connect PostgreSQL + MongoDB
	database.ConnectDatabases()

	// 3. Setup Fiber
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	// 🔥 Swagger endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)

	// === Init repository ===
	authRepo := repository.NewAuthRepository()
	userRepo := repository.NewUserRepository()
	studentRepo := repository.NewStudentRepository()
	lecturerRepo := repository.NewLecturerRepository()

	pgAchRepo := repository.NewAchievementPGRepository()
	mongoAchRepo := repository.NewAchievementMongoRepository()

	// === Init services ===
	authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)
	studentService := service.NewStudentService(studentRepo, pgAchRepo, mongoAchRepo)
	lecturerService := service.NewLecturerService(lecturerRepo)
	achievementService := service.NewAchievementService(pgAchRepo, mongoAchRepo, studentRepo)
	reportService := service.NewReportService(pgAchRepo, mongoAchRepo, studentRepo)

	// === Setup routes ===
	routes.SetupRoutes(
		app,
		authService,
		userService,
		studentService,
		lecturerService,
		achievementService,
		reportService,
	)

	// 5. Start server
	port := ":" + config.AppEnv.AppPort
	log.Println("🚀 Server running on port", port)
	app.Listen(port)
}
