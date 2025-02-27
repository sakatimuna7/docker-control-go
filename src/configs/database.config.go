package database

import (
	"docker-control-go/src/database/models"
	logger "docker-control-go/src/log"
	"fmt"
	"log"
	"os"

	"xorm.io/xorm"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *xorm.Engine

func InitDB() {
	dbDriver := os.Getenv("DB_DRIVER") // Tidak perlu load env lagi
	if dbDriver == "" {
		dbDriver = "sqlite3"
	}

	var connectionString string
	switch dbDriver {
	case "postgres":
		connectionString = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	case "sqlite3":
		connectionString = os.Getenv("DB_SQLITE_PATH")
	default:
		log.Fatalf("Unsupported database driver: %s", dbDriver)
	}
	var err error
	DB, err = xorm.NewEngine(dbDriver, connectionString)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
		logger.Log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	if DB == nil {
		log.Fatal("❌ Database connection is nil")
		logger.Log.Fatal("❌ Database connection is nil")
	}
	devMode := os.Getenv("DEV_MODE") == "true"
	DB.ShowSQL(devMode)

	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ Database not responding: %v", err)
		logger.Log.Fatalf("❌ Database not responding: %v", err)
	}

	err = DB.Sync2(new(models.User))
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		logger.Log.Fatalf("Failed to migrate database: %v", err)
	}

	// fmt.Printf("Database connected successfully using %s\n", dbDriver)
	// logger.Log.Infof("Database connected successfully using %s", dbDriver)
}
