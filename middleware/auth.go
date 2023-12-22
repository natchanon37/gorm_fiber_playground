package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	fmt.Println(token)

	return c.Next()
}
