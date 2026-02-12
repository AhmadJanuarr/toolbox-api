package main

import (
	"log"
	"os"
	"time"
	"toolkits/internal/jobs"
	"toolkits/internal/routes"
	"toolkits/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {

	cleaner := jobs.StorageCleanup(
		[]string{
			services.TempProcessed,
			services.TempUploads,
		},
		24*time.Hour,
	)

	cleaner.Start()
	// Memuat file .env jika ada, jika tidak ada maka kembalikan warning
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Mengambil router dari routes/router.go
	router := routes.Route()
	router.MaxMultipartMemory = 20 << 20 // 20MB
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// configuration cors untuk mengizinkan request dari domain tertentu
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{allowedOrigins}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}

	// setup middleware cors
	router.Use(cors.New(config))
	router.Run(":" + port)

}
