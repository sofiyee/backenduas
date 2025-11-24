package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AllowRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		claims := c.Locals("user").(jwt.MapClaims)

		// DEBUG LOG
		fmt.Println("🔥 RBAC middleware dijalankan. CLAIMS:", claims)

		roleName, ok := claims["role_name"].(string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"error": "Role tidak ditemukan dalam token",
			})
		}

		roleName = strings.ToLower(roleName)

		for _, allowed := range roles {
			if roleName == strings.ToLower(allowed) {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden: Anda tidak memiliki akses",
		})
	}
}
