package routes

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing Authorization header",
		})
	}

	// Extract the token from the Authorization header
	secretKey := os.Getenv("JWTKEY")
	tokenString = tokenString[7:] // Remove "Bearer " prefix
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Provide the key to validate the token
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
	}

	// Token is valid, set user information in context
	claims, _ := token.Claims.(jwt.MapClaims)
	c.Locals("user", claims["user"])

	return c.Next()
}
