package seeders

import (
	"docker-control-go/src/database/models"
	"docker-control-go/src/helpers"
	"log"

	"xorm.io/xorm"
)

// UserSeeder menambahkan data awal ke tabel users
func UserSeeder(db *xorm.Engine) {
	users := []models.User{
		{Username: "superadmin", Password: "superadmin@2025", Role: "superadmin"},
		{Username: "user", Password: "password", Role: "user"},
	}

	// hash password
	for i := range users {
		hashedPassword, err := helpers.HashPassword(users[i].Password)
		if err != nil {
			log.Println("Failed to hash password:", err)
		} else {
			users[i].Password = hashedPassword
		}
	}

	for _, user := range users {
		exists, _ := db.Exist(&models.User{Username: user.Username})
		if !exists {
			_, err := db.Insert(&user)
			if err != nil {
				log.Println("Failed to insert user:", err)
			} else {
				log.Println("User inserted:", user.Username)
			}
		}
	}
}
