package config

import (
	"fmt"
	"investment-tracker-backend/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Get PostgreSQL connection details from environment
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatal("DB_PASSWORD environment variable not set")
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "investmentdb"
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode)

	// Connect to PostgreSQL
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	log.Println("PostgreSQL connected successfully")

	// Auto-migrate models - this will create tables if they don't exist
	log.Println("Running auto-migration...")
	err = DB.AutoMigrate(
		&models.User{},
		&models.Budget{},
		&models.Expense{},
		&models.Goal{},
		&models.Investment{},
	)
	if err != nil {
		log.Fatal("Failed to auto-migrate models:", err)
	}
	log.Println("Auto-migration completed successfully")
}

func DisconnectDatabase() {
	if DB == nil {
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatal("Failed to disconnect PostgreSQL:", err)
	}

	log.Println("PostgreSQL disconnected")
}

// Made with Bob
