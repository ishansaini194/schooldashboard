package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		secret := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("user", token.Claims)
		return c.Next()
	}
}

func GetRole(c *fiber.Ctx) string {
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return ""
	}
	role, _ := claims["role"].(string)
	return role
}

func GetUserID(c *fiber.Ctx) uint {
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return 0
	}
	id, _ := claims["user_id"].(float64)
	return uint(id)
}

func GetTeacherID(c *fiber.Ctx) uint {
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return 0
	}
	id, _ := claims["teacher_id"].(float64)
	return uint(id)
}

func GetStudentID(c *fiber.Ctx) uint {
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return 0
	}
	id, _ := claims["student_id"].(float64)
	return uint(id)
}
