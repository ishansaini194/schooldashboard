package handlers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/models"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getJwtSecret())

func getJwtSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not set")
	}
	return secret
}

func Register(c *fiber.Ctx) error {
	var body struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Role      string `json:"role"`
		EpunjabID string `json:"epunjab_id"`
		StudentID uint   `json:"student_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// ── added ──
	if body.EpunjabID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "epunjab_id is required"})
	}
	// ── end ──

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}

	user := models.User{
		Username:  body.Username,
		Password:  string(hash),
		Role:      body.Role,
		EpunjabID: body.EpunjabID, // always set now since it's required
		StudentID: body.StudentID,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "username already exists"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func Login(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	result := database.DB.
		Where("username = ? OR epunjab_id = ?", body.Username, body.Username).
		First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"username":   user.Username,
		"role":       user.Role,
		"student_id": user.StudentID,
		"epunjab_id": user.EpunjabID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"token":      tokenStr,
		"username":   user.Username,
		"role":       user.Role,
		"student_id": user.StudentID,
		"epunjab_id": user.EpunjabID,
	})
}

// PUT /api/auth/change-password
func ChangePassword(c *fiber.Ctx) error {
	claims, _ := c.Locals("user").(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "incorrect current password"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}

	database.DB.Model(&user).Update("password", string(hash))
	return c.JSON(fiber.Map{"message": "password changed"})
}

// PUT /api/auth/reset-password/:user_id
func ResetPassword(c *fiber.Ctx) error {
	claims, _ := c.Locals("user").(jwt.MapClaims)
	role, _ := claims["role"].(string)

	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}

	userID := c.Params("user_id")

	var body struct {
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", string(hash)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "password reset"})
}
