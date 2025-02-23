package routes

import (
	controllers "docker-control-go/src/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// SetupRoutes untuk Container API
func SetupRoutesWS(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket untuk event Docker container
	app.Get("/ws/containers", websocket.New(controllers.GetRunningContainersWS))

	// WebSocket lain untuk event Docker
	app.Get("/ws", websocket.New(controllers.DockerEventsWS))
}
