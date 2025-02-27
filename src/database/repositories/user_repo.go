package repositories

import (
	database "docker-control-go/src/configs"
	"docker-control-go/src/database/models"
	"errors"
)

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := database.DB.Find(&users)
	return users, err
}

func CreateUser(user *models.User) error {
	_, err := database.DB.Insert(user)
	return err
}

func GetUserByID(id int64) (*models.User, error) {
	var user models.User
	has, err := database.DB.ID(id).Get(&user)
	if !has {
		return nil, err
	}
	return &user, err
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	has, err := database.DB.Where("username = ?", username).Get(&user)
	if err != nil {
		return nil, err // Jika query error (misalnya masalah koneksi DB)
	}
	if !has {
		return nil, errors.New("user not found") // Tambahkan error jika user tidak ditemukan
	}
	return &user, nil
}

func UpdateUser(user *models.User) error {
	_, err := database.DB.ID(user.ID).Update(user)
	return err
}
