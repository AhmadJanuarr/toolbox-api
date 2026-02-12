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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Warning: .env file not found, using system defaults")
	}

	// Mengambil router dari routes/router.go
	router := routes.Route()
	router.MaxMultipartMemory = 20 << 20 // 20MB
	env := os.Getenv("PORT")
	if env == "" {
		log.Fatal("PORT is not set in the .env file")
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
	router.Run(":" + env)

}
