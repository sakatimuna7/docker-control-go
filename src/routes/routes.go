package routes

import (
	"docker-control-go/src/controllers"
	"docker-control-go/src/helpers"
	middleware "docker-control-go/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return helpers.SuccessResponse(c, 200, "Welcome to Docker Control API", nil)
	})
	api.Get("/users", middleware.JWTMiddleware, middleware.ActivityLogger("Users fetched"), controllers.GetUsers)
	api.Post("/user", middleware.JWTMiddleware, middleware.ActivityLogger("User created"), controllers.CreateUser)
	api.Post("/user/login", controllers.UserLogin)
	app.Get("/api/docker/containers", middleware.JWTMiddleware, controllers.GetRunningConainters)
}
