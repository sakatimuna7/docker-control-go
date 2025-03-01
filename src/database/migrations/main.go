package main

import (
	"docker-control-go/src/configs"
	"docker-control-go/src/database/models"
	"docker-control-go/src/database/seeders"
	logger "docker-control-go/src/log"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // Import driver SQLite
)

func main() {
	var err error
	// Load file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using default settings from sedding")
		logger.Log.Warn("Warning: No .env file found, using default settings from sedding")
	}

	// Inisialisasi Logger
	logger.InitLogger()
	logger.Log.Info("Seeding Logger initialized!")

	// Inisialisasi database
	configs.InitDB()
	if configs.DB == nil {
		log.Fatal("❌ Database not initialized")
		logger.Log.Error("❌ Database not initialized")
	}

	defer func() {
		if configs.DB != nil {
			logger.Log.Info("Database connection closed")
			configs.DB.Close()
		}
	}()
	// Inisialisasi Casbin
	configs.InitCasbin(configs.DB)

	// Reset data Casbin
	seeders.ResetCasbinData()
	// Hapus semua tabel
	configs.DB.DropTables(new(models.User))
	// Sinkronisasi ulang tabel (migrasi ulang)
	err = configs.DB.Sync2(new(models.User))
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database berhasil di-reset dan dimigrasi ulang!")
}
