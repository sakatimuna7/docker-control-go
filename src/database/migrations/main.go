package main

import (
	"docker-control-go/src/database/models"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import driver SQLite
	"xorm.io/xorm"
)

var DB *xorm.Engine

func main() {
	var err error
	DB, err = xorm.NewEngine("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Hapus semua tabel
	DB.DropTables(new(models.User))
	// Sinkronisasi ulang tabel (migrasi ulang)
	err = DB.Sync2(new(models.User))
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database berhasil di-reset dan dimigrasi ulang!")
}
