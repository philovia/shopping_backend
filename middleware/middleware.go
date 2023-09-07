package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	// "github.com/mreym/go-fiber-postgres/tokens"
	// "github.com/mreym/go-fiber-postgres/tokens"

)

func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ClientToken := c.Get("tokens")
		if ClientToken == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "No authorization header provided"})
		}

		// claims, msg := tokens.ValidateToken(ClientToken)
		// if msg != "" {
		// 	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": msg})
		// }

		// c.Locals("emails", claims.Email)
		// c.Locals("uid", claims.Uid)
		return c.Next()
	}
}
