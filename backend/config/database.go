package main

import (
	"investment-tracker-backend/config"
	"investment-tracker-backend/controllers"
	"investment-tracker-backend/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	// Initialize database
	config.ConnectDatabase()

	// Initialize OAuth configuration
	controllers.InitOAuth()

	// Create Gin router
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Get port from environment (Render provides PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Update CORS for production
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}

	// CORS middleware - allow credentials for cookies
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	log.Printf("✅ Server starting on port %s", port)
	log.Println("✅ Google OAuth configured")
	log.Println("✅ MongoDB connected")
	log.Printf("✅ CORS configured for: %s", allowedOrigins)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
