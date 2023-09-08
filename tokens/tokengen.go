package tokens

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "github.com/mreym/go-fiber-postgres/models"
)

type NewUser struct {
	ID         uint   `gorm:"primaryKey"`
	Email      string `gorm:"unique"`
	Password   string
	Created_At time.Time
	Updated_At time.Time
}
type User struct {
	ID         uint
	Email      string
	Password   string
	First_Name string
	Last_Name  string
}

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var db *gorm.DB
var SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASS") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	var errDB error
	db, errDB = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDB != nil {
		log.Fatalf("Failed to connect to the database: %v", errDB)
	}

	db.AutoMigrate(&User{}) // Create the User table if it doesn't exist

	app := fiber.New()

	app.Use(recover.New())

	// Routes
	app.Post("/register", Register)
	app.Post("/login", Login)

	// Example protected route
	app.Use(AuthMiddleware)

	app.Get("/protected", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*SignedDetails)
		return c.JSON(fiber.Map{"message": "You have access to this protected route", "user": user})
	})

	err = app.Listen(":8081")
	if err != nil {
		log.Fatal(err)
	}
}

func Register(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Check if the email is already registered
	var existingUser User
	result := db.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "Email already exists"})
	}

	// Hash the password before storing it in the database (you should use a password hashing library)
	user.Password = HashPassword(user.Password)

	// Create the user
	db.Create(&user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Registration successful"})
}

func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Find the user by email
	var user User
	result := db.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	// Verify the password (you should use a password hashing library)
	if !CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Incorrect password"})
	}

	// Generate JWT tokens
	signedToken, signedRefreshToken, err := GenerateTokens(user.Email, user.First_Name, user.Last_Name, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Token generation failed"})
	}

	return c.JSON(fiber.Map{"message": "Login successful", "token": signedToken, "refresh_token": signedRefreshToken})
}

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	claims, msg := ValidateToken(token)
	if msg != "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": msg})
	}

	c.Locals("user", claims)
	return c.Next()
}

// Implement your own password hashing and verification functions here
func HashPassword(password string) string {
	// Replace this with a secure password hashing library (e.g., bcrypt)
	return password
}

func CheckPasswordHash(password, hash string) bool {
	// Replace this with the corresponding password verification logic
	// using the same secure password hashing library
	return password == hash
}

func GenerateTokens(email string, firstname string, lastname string, uid uint) (signedToken string, signedRefreshToken string, err error) {
	// Token generation logic
	// Use the provided user information to create JWT tokens
	// Set appropriate claims and sign the tokens with SECRET_KEY
	// Return the signed tokens and any errors
	claims := &SignedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		Uid:        string(rune(uid)),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 168).Unix(), // 1 week
		},
	}
	signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	// Token validation logic
	// Parse the provided token using SECRET_KEY
	// Verify the token's expiration and integrity
	// Return the claims and an error message if applicable
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "Invalid token claims"
		return
	}

	if claims.ExpiresAt < time.Now().Unix() {
		msg = "Token has expired"
		return
	}

	return claims, msg
}
