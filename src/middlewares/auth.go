package middlewares

import (
	"docker-control-go/src/configs"
	"docker-control-go/src/helpers"
	logger "docker-control-go/src/log"

	"github.com/gofiber/fiber/v2"
)

func Authorize(obj string, act string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userRole, ok := c.Locals("userRole").(string)
		if !ok {
			logger.Log.Info("❌ userRole not found in context")
			return helpers.ErrorResponse(c, 401, "Unauthorized: Missing user role", nil)
		}

		allowed, err := configs.Enforcer.Enforce(userRole, obj, act)
		if err != nil {
			logger.Log.Info("Authorization error:", err)
			return helpers.ErrorResponse(c, 401, "Authorization error", err)
		}

		if !allowed {
			logger.Log.Info("Permission denied")
			return helpers.ErrorResponse(c, 403, "Permission denied", nil)
		}

		return c.Next()
	}
}
