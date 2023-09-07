package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/recover"
	// "github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	// "github.com/mreym/go-fiber-postgres/middleware"
	"github.com/mreym/go-fiber-postgres/models"
	"github.com/mreym/go-fiber-postgres/storage"
	"github.com/mreym/go-fiber-postgres/tokens"
)

var Validate = validator.New()

var (
	DB *gorm.DB
)

func SetupDatabase() {
	dsn := ("%s=localhost %s=5433 %s=shopping_end %s=postgres %p=ostgres %s=disable")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = DB.AutoMigrate(&models.Users{}, &models.Product{})
	if err != nil {
		log.Fatalf("Failed to auto migrate the database: %v", err)
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(usePassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(usePassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}

func Signup(c *fiber.Ctx) error {
	var user models.Users
	if err := c.BodyParser(&user); err != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return err
	}

	validationErr := Validate.Struct(user)
	if validationErr != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})
		return validationErr
	}

	var count int64
	DB.Model(&models.Users{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "user already exists"})
		return nil
	}

	count = 0
	DB.Model(&models.Users{}).Where("phone = ?", user.Phone).Count(&count)
	if count > 0 {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "this phone no. is already in use"})
		return nil
	}

	password := HashPassword(*&user.Password)
	user.Password = *&password

	user.Created_At = time.Now()
	user.Updated_At = time.Now()
	user.UserCart = make([]models.ProductUser, 0)
	user.Address_Details = make([]models.Address, 0)
	user.Order_Status = make([]models.Order, 0)
	user.ID = 0
	// user.User_ID = ""

	token, refreshToken, _ := tokens.GenerateTokens(user.Email, user.First_Name, user.Last_Name, uint(user.User_ID))
	user.Token = token
	user.Refresh_Token = refreshToken

	err := DB.Create(&user).Error
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "the user did not get created"})
		return err
	}

	c.Status(http.StatusCreated)
	return nil
}

func Login(c *fiber.Ctx) error {
	var user models.Users
	if err := c.BodyParser(&user); err != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return err
	}

	var foundUser models.Users
	err := DB.Model(&models.Users{}).Where("email = ?", user.Email).First(&foundUser).Error
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "login or password incorrect"})
		return err
	}

	passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
	if !passwordIsValid {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": msg})
		fmt.Println(msg)
		return nil
	}

	// token, refreshToken, _ := tokens.GenerateTokens(foundUser.Email, foundUser.First_Name, foundUser.Last_Name, uint(foundUser.User_ID))

	// err = tokens.UpdateAllTokens(token, refreshToken, foundUser.User_ID)
	// // tokErr := tokens.UpdateAllTokens()
	// if err != nil {
	// 	c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update tokens"})
	// 	return err
	// }

	c.Status(http.StatusFound).JSON(foundUser)
	return nil
}
func ProductViewerAdmin(c *fiber.Ctx) error {
	var products models.Product
	if err := c.BodyParser(&products); err != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return err
	}

	products.Product_ID = 0 //
	err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}},
		DoNothing: true,
	}).Create(&products).Error

	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "not inserted"})
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func SearchProduct(c *fiber.Ctx) error {
	var productList []models.Product

	err := DB.Find(&productList).Error
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON("something went wrong, please try after some time")
		return err
	}

	c.Status(http.StatusOK).JSON(productList)
	return nil
}

func SearchProductByQuery(c *fiber.Ctx) error {
	var searchProducts []models.Product
	queryParam := c.Query("name")

	if queryParam == "" {
		c.Status(http.StatusNotFound).JSON(fiber.Map{"Error": "Invalid search index"})
		return nil
	}

	err := DB.Where("product_name LIKE ?", "%"+queryParam+"%").Find(&searchProducts).Error
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "something went wrong while fetching the data"})
		return err
	}

	c.Status(http.StatusOK).JSON(searchProducts)
	return nil
}

// Example handler to get user data
func GetUserHandler(c *fiber.Ctx) error {
	userData := storage.UserData("users")
	var users []models.Users
	userData.Find(&users)

	// Return the users as JSON response
	return c.JSON(users)
}

// Example handler to get product data
func GetProductHandler(c *fiber.Ctx) error {
	productData := storage.ProductData("products")
	var products []models.Product
	productData.Find(&products)

	// Return the products as JSON response
	return c.JSON(products)
}

// func main() {

// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	SetupDatabase()

// 	app := fiber.New()

// 	app.Use(recover.New())

// 	app.Post("/signup", Signup)
// 	app.Post("/login", Login)
// 	app.Get("/search", SearchProductByQuery)
// 	app.Post("/admin/products", ProductViewerAdmin)

// 	err = app.Listen(":8080")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
