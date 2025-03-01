package services

import (
	"docker-control-go/src/database/models"
	"docker-control-go/src/database/repositories"
	"docker-control-go/src/database/validations"
	"docker-control-go/src/helpers"
	"errors"

	"github.com/google/uuid"
)

func FetchUsers(mask bool) (interface{}, error) {
	// Pastikan data ter-masking
	users, err := repositories.GetAllUsers()
	if err != nil {
		return nil, err
	}
	maskedUsers := helpers.MaskPrivateFields(users, mask)

	return maskedUsers, nil
}

func RegisterUser(payload *validations.UserCreatePayload) (*interface{}, error) {
	var err error

	hashPassword, err := helpers.HashPassword(payload.Password)
	if err != nil {
		return nil, errors.New("Failed to hash password : " + err.Error())
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
		return nil, errors.New("Failed to create user : " + err.Error())
	}
	maskedUser := helpers.MaskPrivateFields(user, true)

	return &maskedUser, nil
}

func LoginUser(username, password string) (*models.User, error) {
	user, err := repositories.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if !helpers.ComparePassword(password, user.Password) {
		return nil, errors.New("invalid password")
	}
	return user, nil
}
