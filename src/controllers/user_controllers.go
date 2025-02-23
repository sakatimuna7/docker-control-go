package controllers

import (
	"docker-control-go/src/database/models"
	"docker-control-go/src/database/validations"
	"docker-control-go/src/helpers"
	middleware "docker-control-go/src/middlewares"
	"docker-control-go/src/services"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	users, err := services.FetchUsers()
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to fetch users", err)
	}
	mask := c.Query("mask") == "true"

	// Pastikan data ter-masking
	maskedUsers := helpers.MaskPrivateFields(users, mask)
	return helpers.SuccessResponse(c, 200, "Users fetched successfully", maskedUsers)
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

	hashPassword, err := helpers.HashPassword(payload.Password)
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to hash password", err)
	}

	payload.Password = hashPassword

	// Buat user baru
	user := models.User{
		Username: payload.Username,
		Password: payload.Password,
		Role:     "user",
	}

	if err := services.RegisterUser(&user); err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to create user", err)
	}
	// Pastikan data ter-masking
	maskedUser := helpers.MaskPrivateFields(user, true)

	return helpers.SuccessResponse(c, 200, "User created successfully", maskedUser)
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

	token, _ := middleware.GenerateToken(uint(user.ID), user.Role)
	services.LogActivity(user.ID, "User login")
	return helpers.SuccessResponse(c, 200, "User logged in successfully", map[string]interface{}{
		"token": token,
	})

}
