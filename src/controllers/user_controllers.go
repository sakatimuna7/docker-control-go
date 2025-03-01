package controllers

import (
	"docker-control-go/src/database/repositories"
	"docker-control-go/src/database/validations"
	"docker-control-go/src/helpers"
	middleware "docker-control-go/src/middlewares"
	"docker-control-go/src/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	maskStr := c.Query("mask", "true") // Default ke "true" jika tidak ada query
	mask, _ := strconv.ParseBool(maskStr)

	users, err := services.FetchUsers(mask)
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to fetch users", err)
	}
	return helpers.SuccessResponse(c, 200, "Users fetched successfully", users)
}

func CreateUser(c *fiber.Ctx) error {
	var payload validations.UserCreatePayload

	// Parse request body
	if err := c.BodyParser(&payload); err != nil {
		return helpers.ErrorResponse(c, 400, "Invalid request payload", err)
	}

	// Validasi input
	if errors := helpers.ValidateStruct(payload); len(errors) > 0 {
		return helpers.ErrorResponse(c, 400, "Invalid request payload", errors)
	}

	// create user
	user, err := services.RegisterUser(&payload)

	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to create user", err)
	}

	return helpers.SuccessResponse(c, 200, "User created successfully", user)
}

func GetUserByID(c *fiber.Ctx) error {
	// Ambil parameter "id" dari URL
	id := c.Params("id")
	// Konversi string ke int64
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
	// Gunakan userID dalam fungsi
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to fetch user", err)
	}
	mask := c.Query("mask") == "true"
	// Pastikan data ter-masking
	maskedUser := helpers.MaskPrivateFields(user, mask)
	return helpers.SuccessResponse(c, 200, "User fetched successfully", maskedUser)
}

func UpdateUser(c *fiber.Ctx) error {
	// Ambil parameter "id" dari URL
	id := c.Params("id")
	// Konversi string ke int64
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return helpers.ErrorResponse(c, 400, "Invalid user ID", err)
	}
	// Gunakan userID dalam fungsi
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to fetch user", err)
	}

	var payload validations.UserUpdatePayload

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

	user.Password = hashPassword

	if err := repositories.UpdateUser(user); err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to update user", err)
	}
	// Pastikan data ter-masking
	maskedUser := helpers.MaskPrivateFields(user, true)
	return helpers.SuccessResponse(c, 200, "User updated successfully", maskedUser)
}

func UserLogin(c *fiber.Ctx) error {
	var payload validations.UserCreatePayload

	// Parse request body
	if err := c.BodyParser(&payload); err != nil {
		return helpers.ErrorResponse(c, 400, "Invalid request payload", err)
	}

	// Validasi input
	if errors := helpers.ValidateStruct(payload); len(errors) > 0 {
		return helpers.ErrorResponse(c, 400, "Invalid request payload", errors)
	}

	user, err := services.LoginUser(payload.Username, payload.Password)
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to login user", err)
	}

	token, _ := middleware.GenerateToken(user)
	return helpers.SuccessResponse(c, 200, "User logged in successfully", map[string]interface{}{
		"token": token,
	})

}
