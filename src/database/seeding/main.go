package main

import (
	"log"

	database "docker-control-go/src/configs"
	"docker-control-go/src/database/seeders"
	logger "docker-control-go/src/log"

	"github.com/joho/godotenv"
	"xorm.io/xorm"
)

func SeedAll(db *xorm.Engine) {
	log.Println("Running all seeders...")

	seeders.UserSeeder(db)

	log.Println("All seeders completed!")
}

func main() {
	// Load file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using default settings from sedding")
		logger.Log.Warn("Warning: No .env file found, using default settings from sedding")
	}

	// Inisialisasi Logger
	logger.InitLogger()
	logger.Log.Info("Seeding Logger initialized!")

	// Inisialisasi database
	database.InitDB()
	if database.DB == nil {
		log.Fatal("❌ Database not initialized")
		logger.Log.Error("❌ Database not initialized")
	}

	defer func() {
		if database.DB != nil {
			logger.Log.Info("Database connection closed")
			database.DB.Close()
		}
	}()

	// Jalankan semua seeder
	SeedAll(database.DB)
	log.Println("All seeders executed successfully!")
	logger.Log.Info("All seeders executed successfully!")
}
