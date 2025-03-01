package controllers

import (
	"docker-control-go/src/configs"
	"docker-control-go/src/database/models"
	"docker-control-go/src/helpers"
	logger "docker-control-go/src/log"

	"github.com/gofiber/fiber/v2"
)

func AddRolePermission(c *fiber.Ctx) error {

	var req models.Permission
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to add permission to Casbin: ", err)
		return helpers.ErrorResponse(c, 400, "Invalid request", err)
	}

	_, err := configs.Enforcer.AddPolicy(req.Role, req.Obj, req.Act)
	if err != nil {
		logger.Log.Error("Failed to add permission to Casbin: ", err)
		return helpers.ErrorResponse(c, 500, "Failed to add permission", err)
	}

	logger.Log.Info("Permission added successfully to Casbin")
	return helpers.SuccessResponse(c, 201, "Permission added successfully", nil)
}
