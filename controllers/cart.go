package controllers

import (
	// "context"
	"net/http"
	"strconv"
	// "time"
	// "go-fiber-postgres/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/mreym/go-fiber-postgres/models"

)

type Application struct {
	db *gorm.DB
}

func NewApplication(db *gorm.DB) *Application {
	return &Application{
		db: db,
	}
}

func (app *Application) AddToCart() fiber.Handler {
	return func(c *fiber.Ctx) error {
		productQueryID := c.Query("id")
		userQueryID := c.Query("userID")

		_, err := strconv.Atoi(productQueryID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
		}

		_, err = strconv.Atoi(userQueryID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully added to the cart"})
	}
}
func (app *Application) RemoveItem() fiber.Handler {
	return func(c *fiber.Ctx) error {
		productQueryID := c.Query("id")
		userQueryID := c.Query("userID")

		_, err := strconv.Atoi(productQueryID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
		}

		_, err = strconv.Atoi(userQueryID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully removed item from cart"})
	}
}
func (app *Application) GetItemFromCart() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userQueryID := c.Query("id")

		userID, err := strconv.Atoi(userQueryID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		var usercart []models.ProductUser
		app.db.Where("user_id = ?", userID).Find(&usercart)

		totalItems := len(usercart)

		return c.Status(http.StatusOK).JSON(fiber.Map{"totalItems": totalItems, "cartItems": usercart})
	}
}

func (app *Application) BuyFromCart() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userQueryID := c.Query("id")

		_, err := strconv.Atoi(userQueryID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully placed the order"})
	}
}

func (app *Application) InstantBuy() fiber.Handler {
	return func(c *fiber.Ctx) error {
		productQueryID := c.Query("id")
		userQueryID := c.Query("userID")

		_, err := strconv.Atoi(productQueryID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
		}

		_, err = strconv.Atoi(userQueryID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Successfully placed the order"})
	}
}
