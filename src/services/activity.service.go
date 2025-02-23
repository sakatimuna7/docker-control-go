package services

import (
	"docker-control-go/src/database/models"
	"docker-control-go/src/database/repositories"
	"log"
	"time"
)

func LogActivity(userID int64, action string) error {
	log_activity := models.ActivityLog{
		UserID:    int64(userID), // Perbaikan konversi
		Action:    action,
		Timestamp: time.Now(),
	}

	err := repositories.CreateLogActivity(&log_activity)
	if err != nil {
		log.Printf("Failed to insert log: %v", err)
	}
	return err
}
