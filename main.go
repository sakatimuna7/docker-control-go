package main

import (
	"docker-control-go/src/configs"
	controllers "docker-control-go/src/controllers"
	logger "docker-control-go/src/log"
	"docker-control-go/src/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using default settings")
		logger.Log.Warn("Warning: No .env file found, using default settings")
	}

	// Inisialisasi Logger
	logger.InitLogger()
	logger.Log.Info("Logger initialized!")

	// Inisialisasi Docker Client
	controllers.InitDockerClient()

	// Ambil port dari .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Init database
	configs.InitDB()
	configs.InitRedis()

	// Inisialisasi Casbin
	configs.InitCasbin(configs.DB)

	// Inisialisasi Fiber
	app := fiber.New()

	// Setup Routes
	routes.SetupRoutesWS(app)
	routes.SetupRoutes(app)

	// Jalankan server
	logger.Log.Info("App running on port " + port)
	log.Fatal(app.Listen(":" + port))
}
