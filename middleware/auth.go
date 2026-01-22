package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthGuard(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization") // Next.js kirim 'Bearer <token>'
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Kamu belum login"})
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("user_id", uint(claims["user_id"].(float64)))

	return c.Next()
}