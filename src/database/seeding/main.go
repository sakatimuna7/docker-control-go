package main

import (
	"log"

	"docker-control-go/src/database"
	"docker-control-go/src/database/seeders" // ✅ Import seeders package

	"xorm.io/xorm"
)

func SeedAll(db *xorm.Engine) {
	log.Println("Running all seeders...")

	seeders.UserSeeder(db)

	log.Println("All seeders completed!")
}

func main() {
	// Inisialisasi database
	database.InitDB()
	defer database.DB.Close()

	// Jalankan semua seeder
	SeedAll(database.DB)
	log.Println("All seeders executed successfully!")
}
