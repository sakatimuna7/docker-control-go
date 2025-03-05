package controllers

import (
	"docker-control-go/src/configs"
	"docker-control-go/src/helpers"
	logger "docker-control-go/src/log"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type PermissionProcessPayload struct {
	UserID       string `json:"user_id"`
	IdentityName string `json:"identity_name"`
	Action       string `json:"action"` // Misal: "read", "start", "stop", dll.
}

func AddProcessPermission(c *fiber.Ctx) error {
	req := new(PermissionProcessPayload)
	if err := c.BodyParser(req); err != nil {
		logger.Log.Info("Invalid request ", err)
		return helpers.ErrorResponse(c, 400, "Invalid request", err)
	}
	pm2IdentityName := fmt.Sprintf("pm2:%v", req.IdentityName)

	// Tambahkan aturan Casbin (user, container, action)
	success, err := configs.Enforcer.AddPolicy(req.UserID, pm2IdentityName, req.Action)
	if err != nil {
		logger.Log.Error("Failed to add permission ", err)
		return helpers.ErrorResponse(c, 500, "Failed to add permission", err)
	}

	if !success {
		logger.Log.Info("Permission already exists ", err)
		return helpers.ErrorResponse(c, 400, "Permission already exists", nil)
	}

	logger.Log.Info("Permission added successfully to Casbin")
	return helpers.SuccessResponse(c, 201, "Permission added successfully", nil)
}

// Hapus izin user dari process
func RemoveProcessPermission(c *fiber.Ctx) error {
	req := new(PermissionProcessPayload)
	if err := c.BodyParser(req); err != nil {
		logger.Log.Error("Invalid request ", err)
		return helpers.ErrorResponse(c, 400, "Invalid request", err)
	}

	pm2IdentityName := fmt.Sprintf("pm2:%v", req.IdentityName)

	// Hapus aturan Casbin (user, container, action)
	success, err := configs.Enforcer.RemovePolicy(req.UserID, pm2IdentityName, req.Action)
	if err != nil {
		logger.Log.Error("Failed to remove permission ", err)
		return helpers.ErrorResponse(c, 500, "Failed to remove permission", err)
	}

	if !success {
		logger.Log.Info("Permission already exists ", err)
		return helpers.ErrorResponse(c, 404, "Permission not found", nil)
	}
	logger.Log.Info("Permission removed successfully ")
	return helpers.SuccessResponse(c, 200, "Permission removed successfully", nil)
}
