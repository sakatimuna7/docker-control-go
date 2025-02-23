package database

import (
	"docker-control-go/src/database/models"
	"fmt"
	"log"

	"xorm.io/xorm"

	_ "github.com/mattn/go-sqlite3" // Import driver SQLite
)

var DB *xorm.Engine

func InitDB() {
	var err error
	DB, err = xorm.NewEngine("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Aktifkan logging untuk debugging
	DB.ShowSQL(true)

	// Auto migrate table
	err = DB.Sync2(new(models.User), new(models.ActivityLog))
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database connected & migrated successfully")
}
