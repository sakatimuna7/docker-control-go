package routes

import (
	"docker-control-go/src/controllers"
	"docker-control-go/src/helpers"
	"docker-control-go/src/middlewares"

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
	users.Get("/", middlewares.JWTMiddleware, middlewares.Authorize("resource:user", "read"), controllers.GetUsers)
	users.Post("/", middlewares.JWTMiddleware, middlewares.Authorize("resource:user", "write"), controllers.CreateUser)
	users.Post("/login", controllers.UserLogin)

	// roles
	roles := api.Group("/roles", middlewares.JWTMiddleware)
	// roles.Get("/", middlewares.JWTMiddleware, controllers.GetRoles)
	roles.Post("add-permission", controllers.AddRolePermission)

	// Group untuk Docker
	docker := api.Group("/docker")
	docker.Get("/containers", middlewares.JWTMiddleware, controllers.GetRunningConainters)
}
