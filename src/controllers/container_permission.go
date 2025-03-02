package controllers

import (
	"docker-control-go/src/configs"
	"docker-control-go/src/database/models"
	"docker-control-go/src/helpers"
	logger "docker-control-go/src/log"

	"github.com/gofiber/fiber/v2"
)

// Tambahkan izin user ke container
func AddContainerPermission(c *fiber.Ctx) error {
	req := new(models.PermissionContainerRequest)
	if err := c.BodyParser(req); err != nil {
		logger.Log.Info("Invalid request ", err)
		return helpers.ErrorResponse(c, 400, "Invalid request", err)
	}

	// Tambahkan aturan Casbin (user, container, action)
	success, err := configs.Enforcer.AddPolicy(req.UserID, req.ContainerName, req.Action)
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

// Hapus izin user dari container
func RemoveContainerPermission(c *fiber.Ctx) error {
	req := new(models.PermissionContainerRequest)
	if err := c.BodyParser(req); err != nil {
		logger.Log.Error("Invalid request ", err)
		return helpers.ErrorResponse(c, 400, "Invalid request", err)
	}

	// Hapus aturan Casbin (user, container, action)
	success, err := configs.Enforcer.RemovePolicy(req.UserID, req.ContainerName, req.Action)
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

// Lihat semua izin user terhadap container
func GetUserPermissions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// Ambil semua izin yang diberikan ke user
	policies, _ := configs.Enforcer.GetFilteredPolicy(0, userID)
	var permissions []map[string]string

	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, map[string]string{
				"container_id": policy[1],
				"action":       policy[2],
			})
		}
	}

	if len(permissions) == 0 {
		logger.Log.Info("No permissions found for user")
		return helpers.SuccessResponse(c, 404, "No permissions found for user", []interface{}{})
	}

	return helpers.SuccessResponse(c, 200, "User permissions fetched successfully", permissions)
}

func GetUserPermissionsByUserID(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Ambil semua izin yang diberikan ke user
	policies, _ := configs.Enforcer.GetFilteredPolicy(0, userID)
	var permissions []map[string]string

	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, map[string]string{
				"container_id": policy[1],
				"action":       policy[2],
			})
		}
	}

	if len(permissions) == 0 {
		logger.Log.Info("No permissions found for user")
		return helpers.SuccessResponse(c, 404, "No permissions found for user", []interface{}{})
	}

	return helpers.SuccessResponse(c, 200, "User permissions fetched successfully", permissions)
}
