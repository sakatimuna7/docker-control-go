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

	// file manager
	file := api.Group("/files", middlewares.JWTMiddleware)
	file.Get("/list", middlewares.Authorize("resource:file-list", "read"), controllers.ListFilesHandler)

	// Group untuk User
	users := api.Group("/users")
	users.Get("/", middlewares.JWTMiddleware, middlewares.Authorize("resource:user", "read"), controllers.GetUsers)
	users.Get("/:id", middlewares.JWTMiddleware, middlewares.Authorize("resource:user", "read"), controllers.GetUserByID)
	users.Get("/:username/username", middlewares.JWTMiddleware, middlewares.Authorize("resource:user", "read"), controllers.GetUserByUsername)
	users.Post("/create", middlewares.JWTMiddleware, middlewares.Authorize("resource:user", "create"), controllers.CreateUser)
	users.Patch("/password", middlewares.JWTMiddleware, middlewares.Authorize("resource:user", "update"), controllers.UpdateUserPassword)
	// Auth
	users.Post("/login", controllers.UserLogin)
	users.Delete("/logout", middlewares.JWTMiddleware, controllers.UserLogout)

	// roles
	// roles := api.Group("/roles", middlewares.JWTMiddleware)
	// roles.Get("/", middlewares.JWTMiddleware, controllers.GetRoles)
	// roles.Post("add-permission", controllers.AddRolePermission)

	// permission group
	permissions := api.Group("/permissions", middlewares.JWTMiddleware)
	permissions.Get("/", middlewares.Authorize("resource:permissions", "read"), controllers.GetAllPermissions)
	permissions.Post("/container/add", middlewares.Authorize("resource:permissions-containers", "create"), controllers.AddProcessPermission)
	permissions.Delete("/container/remove", middlewares.Authorize("resource:permissions-containers", "delete"), controllers.RemoveProcessPermission)
	permissions.Post("/process-manager/add", middlewares.Authorize("resource:permissions-process-manager", "create"), controllers.AddProcessPermission)
	permissions.Delete("/process-manager/remove", middlewares.Authorize("resource:permissions-process-manager", "delete"), controllers.RemoveProcessPermission)

	// Group untuk Container
	container := permissions.Group("/containers")

	container.Get("/", middlewares.JWTMiddleware, middlewares.Authorize("resource:permissions-containers", "read"), controllers.GetUserPermissions)
	container.Get("/:id/user", middlewares.Authorize("resource:permissions-containers", "read-user"), controllers.GetUserPermissionsByUserID)
	container.Post("/add", middlewares.Authorize("resource:permissions-containers", "create"), controllers.AddContainerPermission)
	container.Post("/remove", middlewares.Authorize("resource:permissions-containers", "delete"), controllers.RemoveContainerPermission)

	// Group untuk Docker
	docker := api.Group("/docker", middlewares.JWTMiddleware)
	docker.Get("/containers", middlewares.Authorize("resource:docker-containers", "read"), middlewares.JWTMiddleware, controllers.GetRunningConainters)
}
