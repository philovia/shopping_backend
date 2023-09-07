package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mreym/go-fiber-postgres/controllers"
	"github.com/mreym/go-fiber-postgres/middleware"
	// "github.com/mreym/go-fiber-postgres/routes"

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
		log.Fatalf("Failed to connect to the database: %v", err)
	}
}

func main() {
	SetupDatabase()

	app := fiber.New()

	appController := controllers.NewApplication(DB)

	//  http://127.0.0.1:8080/api/addtocart
	app.Use(middleware.Authentication())

	api := app.Group("/api")

	api.Get("/addtocart", appController.AddToCart())
	api.Get("/removeitem", appController.RemoveItem())
	api.Get("/cartcheckout", appController.BuyFromCart())
	api.Get("/instantbuy", appController.InstantBuy())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Fatal(app.Listen(":" + port))
}
