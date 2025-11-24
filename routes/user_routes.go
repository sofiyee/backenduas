package routes

import (
	"backenduas/app/service"
	"backenduas/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(api fiber.Router, userService *service.UserService) {

	u := api.Group("/users",
		middleware.JWTProtected(),
		middleware.AllowRoles("Admin"),
	)

	u.Get("/", userService.GetAll)
	u.Get("/:id", userService.GetByID)
	u.Post("/", userService.Create)
	u.Put("/:id", userService.Update)
	u.Delete("/:id", userService.Delete)
	u.Put("/:id/role", userService.UpdateRole)
}
