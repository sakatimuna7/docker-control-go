package controllers

import (
	"docker-control-go/src/database/models"
	"docker-control-go/src/database/repositories"
	"docker-control-go/src/database/validations"
	"docker-control-go/src/helpers"
	middleware "docker-control-go/src/middlewares"
	"docker-control-go/src/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsers(c *fiber.Ctx) error {
	maskStr := c.Query("mask", "true") // Default ke "true" jika tidak ada query
	mask, _ := strconv.ParseBool(maskStr)

	users, err := repositories.GetAllUsers()
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to fetch users", err)
	}
	maskedUsers := helpers.MaskPrivateFields(users, mask)

	return helpers.SuccessResponse(c, 200, "Users fetched successfully", maskedUsers)
}

func CreateUser(c *fiber.Ctx) error {
	var payload validations.UserCreatePayload
	var err error

	// Parse request body
	if err := c.BodyParser(&payload); err != nil {
		return helpers.ErrorResponse(c, 400, "Invalid request payload", err)
	}

	// Validasi input
	if errors := helpers.ValidateStruct(payload); len(errors) > 0 {
		return helpers.ErrorResponse(c, 400, "Invalid request payload", errors)
	}

	hashPassword, err := helpers.HashPassword(payload.Password)
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to hash password", err)
	}
	payload.Password = hashPassword

	// Buat user baru
	user := models.User{
		ID:       uuid.NewString(),
		Username: payload.Username,
		Password: payload.Password,
		Role:     "user",
	}
	err = repositories.CreateUser(&user)
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to create user", err)
	}
	maskedUser := helpers.MaskPrivateFields(user, true)

	return helpers.SuccessResponse(c, 200, "User created successfully", maskedUser)
}

func GetUserByID(c *fiber.Ctx) error {
	// Ambil parameter "id" dari URL
	userID := c.Params("id")
	maskStr := c.Query("mask", "true") // Default ke "true" jika tidak ada query
	mask, _ := strconv.ParseBool(maskStr)

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return helpers.ErrorResponse(c, 404, "User Not Found", err)
		}
		return helpers.ErrorResponse(c, 500, "Failed to fetch user", err)
	}
	maskedUser := helpers.MaskPrivateFields(user, mask)
	return helpers.SuccessResponse(c, 200, "User fetched successfully", maskedUser)
}
func GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	maskStr := c.Query("mask", "true") // Default ke "true" jika tidak ada query
	mask, _ := strconv.ParseBool(maskStr)

	user, err := repositories.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "user not found" {
			// return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
			return helpers.ErrorResponse(c, 404, "User not found", err)
		}
		// return nil, err // Jika error lain, tetap dikembalikan
		return helpers.ErrorResponse(c, 500, "Failed to fetch user", err)

	}
	maskedUser := helpers.MaskPrivateFields(user, mask)

	return helpers.SuccessResponse(c, 200, "User fetched successfully", maskedUser)
}

func UpdateUserPassword(c *fiber.Ctx) error {
	var payload validations.UserUpdatePasswordPayload
	var err error

	userID := c.Locals("userID").(string) // Ambil role dari context
	// Parse request body
	if err := c.BodyParser(&payload); err != nil {
		return helpers.ErrorResponse(c, 400, "Invalid request payload", err)
	}
	// Validasi input
	if errors := helpers.ValidateStruct(payload); len(errors) > 0 {
		return helpers.ErrorResponse(c, 400, "Invalid request payload", errors)
	}
	// check new password and confirm password
	if payload.Password != payload.ConfirmPassword {
		return helpers.ErrorResponse(c, 400, "password and confirm password not match", "")
	}

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return helpers.ErrorResponse(c, 400, "Failed to update user", err)
	}
	// check old password
	if !helpers.ComparePassword(payload.OldPassword, user.Password) {
		return helpers.ErrorResponse(c, 400, "invalid old password", nil)
	}
	// hash new password
	hashPassword, err := helpers.HashPassword(payload.Password)
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to hash password", err)
	}
	payload.Password = hashPassword
	user.Password = payload.Password
	err = repositories.UpdateUser(user)
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to update user", err)
	}
	maskedUser := helpers.MaskPrivateFields(user, true)

	return helpers.SuccessResponse(c, 201, "User updated successfully", maskedUser)
}

func UserLogin(c *fiber.Ctx) error {
	var payload validations.UserCreatePayload
	var err error
	// Parse request body
	if err := c.BodyParser(&payload); err != nil {
		return helpers.ErrorResponse(c, 400, "Invalid request payload", err)
	}

	// Validasi input
	if errors := helpers.ValidateStruct(payload); len(errors) > 0 {
		return helpers.ErrorResponse(c, 400, "Invalid request payload", errors)
	}

	user, err := repositories.GetUserByUsername(payload.Username)
	if err != nil {
		return helpers.ErrorResponse(c, 404, "Failed to login user", err)
	}
	if !helpers.ComparePassword(payload.Password, user.Password) {
		return helpers.ErrorResponse(c, 400, "invalid password", err)
	}

	token, _ := middleware.GenerateToken(user)
	return helpers.SuccessResponse(c, 200, "User logged in successfully", map[string]interface{}{
		"token": token,
	})
}

func UserLogout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if err := services.DeleteSessionFromRedis(token); err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to logout user", err)
	}
	return helpers.SuccessResponse(c, 200, "User logged out successfully", nil)
}
