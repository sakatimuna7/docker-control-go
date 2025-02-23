package services

import (
	"docker-control-go/src/database/models"
	"docker-control-go/src/database/repositories"
	"docker-control-go/src/helpers"
	"errors"
)

func FetchUsers() ([]models.User, error) {
	return repositories.GetAllUsers()
}

func RegisterUser(user *models.User) error {
	return repositories.CreateUser(user)
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
