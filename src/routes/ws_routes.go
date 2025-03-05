package routes

import (
	controllers "docker-control-go/src/controllers"
	"docker-control-go/src/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// SetupRoutesWS untuk WebSocket
func SetupRoutesWS(app *fiber.App) {
	// Middleware global untuk memastikan hanya WebSocket request yang diterima
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket untuk daftar container dengan autentikasi JWT
	app.Get("/ws/containers", websocket.New(middlewares.WebSocketAuthMiddleware(controllers.GetRunningContainersWS)))

	// WebSocket untuk Docker events (bisa tambahkan autentikasi jika perlu)
	app.Get("/ws/events", websocket.New(middlewares.WebSocketAuthMiddleware(controllers.DockerEventsWS)))

	//  WebSocket untuk daftar process
	app.Get("/ws/web-services", websocket.New(middlewares.WebSocketAuthMiddleware(controllers.PM2Controller)))
}
