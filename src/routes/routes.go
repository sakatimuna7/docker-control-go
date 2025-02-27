package routes

import (
	"docker-control-go/src/controllers"
	"docker-control-go/src/helpers"
	middleware "docker-control-go/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Welcome message
	api.Get("/", func(c *fiber.Ctx) error {
		return helpers.SuccessResponse(c, 200, "Welcome to Docker Control API", nil)
	})

	// Group untuk User
	users := api.Group("/users")
	users.Get("/", middleware.JWTMiddleware, controllers.GetUsers)
	users.Post("/", middleware.JWTMiddleware, controllers.CreateUser)
	users.Post("/login", controllers.UserLogin)

	// Group untuk Docker
	docker := api.Group("/docker")
	docker.Get("/containers", middleware.JWTMiddleware, controllers.GetRunningConainters)
}
