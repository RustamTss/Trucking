package middleware

import (
	"business-schedule-backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Отсутствует токен авторизации"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateJWT(tokenString, secret)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Неверный токен"})
		}

		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}

func GetUserIDFromToken(c *fiber.Ctx) (string, error) {
	userID := c.Locals("userID").(string)
	return userID, nil
}
