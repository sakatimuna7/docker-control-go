package seeders

import (
	"docker-control-go/src/configs"
	"docker-control-go/src/database/models"
	"docker-control-go/src/helpers"
	"log"

	"github.com/google/uuid"
	"xorm.io/xorm"
)

// UserSeeder menambahkan data awal ke tabel users
func UserSeeder(db *xorm.Engine) {
	users := []models.User{
		{ID: uuid.NewString(), Username: "superadmin", Password: "superadmin@2025", Role: "admin"},
		{ID: uuid.NewString(), Username: "user", Password: "password", Role: "user"},
	}

	// Hash password
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
				log.Println("✅ User inserted:", user.Username)

				// Tambahkan user ke role di Casbin
				_, err := configs.Enforcer.AddGroupingPolicy(user.Username, user.Role)
				if err != nil {
					log.Println("❌ Failed to assign Casbin role:", err)
				} else {
					log.Println("✅ Assigned Casbin role:", user.Username, "->", user.Role)
				}
			}
		}
	}

	// Simpan perubahan ke database Casbin
	err := configs.Enforcer.SavePolicy()
	if err != nil {
		log.Println("❌ Failed to save Casbin policies:", err)
	} else {
		log.Println("✅ Casbin policies saved successfully!")
	}
}
