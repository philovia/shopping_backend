package controllers

import (
	"net/http"
	// Import your PostgreSQL driver here

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/mreym/go-fiber-postgres/models"

)

var app = fiber.New()
var db *gorm.DB

func Initialize(databaseURL string) {
	var err error
	db, err = gorm.Open("%s=localhost %s=5433 %s=shopping_end %s=postgres %p=ostgres %s=disable")
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&models.Users{})
	db.AutoMigrate(&models.Address{})

	app.Get("/address", GetAddress)
	app.Post("/address", AddAddress)
	app.Put("/address/:id", EditAddress)
	app.Delete("/address/:id", DeleteAddress)
}

func AddAddress(c *fiber.Ctx) error {
	userID := c.Query("id")
	if userID == "" {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Invalid code"})
	}

	var address models.Address
	if err := c.BodyParser(&address); err != nil {
		return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.Users
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	user.Address_Details = append(user.Address_Details, address)

	if err := db.Save(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully added the address"})
}

func EditAddress(c *fiber.Ctx) error {

	addressID := c.Params("id")

	userID := c.Query("id")

	if userID == "" || addressID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var editAddress models.Address
	if err := c.BodyParser(&editAddress); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.Users
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	for i, address_id := range user.Address_Details {
		if address_id == address_id {
			user.Address_Details[i] = editAddress
		}
	}	

	if err := db.Save(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully updated the address"})
}

func DeleteAddress(c *fiber.Ctx) error {

	addressID := c.Params("id")

	userID := c.Query("id")

	if userID == "" || addressID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.Users
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	for i, address_id := range user.Address_Details {
		if address_id == address_id {
			user.Address_Details = append(user.Address_Details[:i], user.Address_Details[i+1:]...)
		}
	}

	if err := db.Save(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully deleted the address"})
}

func GetAddress(c *fiber.Ctx) error {

	userID := c.Query("id")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var addresses []models.Address

	var user models.Users
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	if err := db.Model(&user).Association("Addresses").Find(&addresses); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(addresses)
}
