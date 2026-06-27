package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/models"
	"lumora/backend/utils"
)

// Protected returns a Fiber middleware that requires a valid Bearer token and
// attaches the resolved *models.User to the request locals under "user".
func Protected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")

		userID, err := utils.ParseToken(tokenString, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		}

		c.Locals("user", &user)
		return c.Next()
	}
}

// CurrentUser is a small helper for controllers to fetch the authenticated user.
func CurrentUser(c *fiber.Ctx) *models.User {
	if u, ok := c.Locals("user").(*models.User); ok {
		return u
	}
	return nil
}
