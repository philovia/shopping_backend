package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "github.com/mreym/go-fiber-postgres/"
)

var DB *gorm.DB

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func SetupDatabase() {
	LoadEnv()
	dsn := fmt.Sprintf("%s=localhost %s=5433 %s=shopping_end %s=postgres %s=postgres %s=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Successfully connected to PostgreSQL database")
}

func UserData(collectionName string) *gorm.DB {
	return DB.Table(collectionName)
}

func ProductData(collectionName string) *gorm.DB {
	return DB.Table(collectionName)
}
