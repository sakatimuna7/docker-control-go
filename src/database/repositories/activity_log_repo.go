package repositories

import (
	"docker-control-go/src/database"
	"docker-control-go/src/database/models"
)

func CreateLogActivity(activity *models.ActivityLog) error {
	_, err := database.DB.Insert(activity)
	return err
}
