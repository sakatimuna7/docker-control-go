package middlewares

import (
	"docker-control-go/src/configs"
	logger "docker-control-go/src/log"

	"github.com/gofiber/fiber/v2"
)

func Authorize(obj string, act string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole").(string) // Ambil role dari context

		allowed, err := configs.Enforcer.Enforce(userRole, obj, act)
		if err != nil {
			logger.Log.Info("Authorization error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Authorization error"})
		}

		if !allowed {
			logger.Log.Info("Permission denied")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
		}

		return c.Next()
	}
}
