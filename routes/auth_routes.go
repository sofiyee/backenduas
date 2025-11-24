package routes

import (
	"backenduas/app/service"
	"backenduas/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(api fiber.Router, authService *service.AuthService) {

	auth := api.Group("/auth")

	auth.Post("/login", authService.Login)

	auth.Get("/profile",
		middleware.JWTProtected(),
		authService.Profile,
	)

	auth.Post("/logout",
		middleware.JWTProtected(),
		authService.Logout,
	)

	auth.Post("/refresh",
		middleware.JWTProtected(),
		authService.Refresh,
	)
}
