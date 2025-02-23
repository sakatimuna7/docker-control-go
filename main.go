package main

import (
	controllers "docker-control-go/src/controllers"
	"docker-control-go/src/database"
	"docker-control-go/src/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load file .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Inisialisasi Docker Client
	controllers.InitDockerClient()

	// Ambil port dari .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Init database
	database.InitDB()

	// Inisialisasi Fiber
	app := fiber.New()

	// Setup Routes
	routes.SetupRoutesWS(app)
	routes.SetupRoutes(app)

	// Jalankan server
	log.Fatal(app.Listen(":" + port))
}
