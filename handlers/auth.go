package handlers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/middleware"
	"github.com/ishansaini194/dashboard/models"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func getSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		panic("JWT_SECRET not set")
	}
	return []byte(s)
}

// POST /api/auth/register
func Register(c *fiber.Ctx) error {
	var body struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Role      string `json:"role"`
		EpunjabID string `json:"epunjab_id"`
		StudentID *uint  `json:"student_id"`
		TeacherID *uint  `json:"teacher_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}

	user := models.User{
		SchoolID:  schoolID(),
		Username:  body.Username,
		Password:  string(hash),
		Role:      body.Role,
		EpunjabID: body.EpunjabID,
		StudentID: body.StudentID,
		TeacherID: body.TeacherID,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "username already exists"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "registered"})
}

// POST /api/auth/login
func Login(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	database.DB.Where("username = ? OR epunjab_id = ?", body.Username, body.Username).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// update last login
	database.DB.Model(&user).Update("last_login", time.Now().Format("2006-01-02 15:04:05"))

	studentID := uint(0)
	if user.StudentID != nil {
		studentID = *user.StudentID
	}
	teacherID := uint(0)
	if user.TeacherID != nil {
		teacherID = *user.TeacherID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"username":   user.Username,
		"role":       user.Role,
		"student_id": studentID,
		"teacher_id": teacherID,
		"epunjab_id": user.EpunjabID,
		"school_id":  user.SchoolID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString(getSecret())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"token":      tokenStr,
		"username":   user.Username,
		"role":       user.Role,
		"student_id": studentID,
		"teacher_id": teacherID,
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
	c.BodyParser(&body)

	var user models.User
	database.DB.First(&user, userID)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "incorrect current password"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	database.DB.Model(&user).Update("password", string(hash))
	return c.JSON(fiber.Map{"message": "password changed"})
}

// PUT /api/auth/reset-password/:user_id
func ResetPassword(c *fiber.Ctx) error {
	if middleware.GetRole(c) != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}

	var body struct {
		NewPassword string `json:"new_password"`
	}
	c.BodyParser(&body)

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	database.DB.Model(&models.User{}).Where("id = ?", c.Params("user_id")).Update("password", string(hash))
	return c.JSON(fiber.Map{"message": "password reset"})
}

// GET /api/users/epunjab/:epunjab_id
func GetUserByEpunjabID(c *fiber.Ctx) error {
	var user models.User
	database.DB.Where("epunjab_id = ?", c.Params("epunjab_id")).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"role":       user.Role,
		"epunjab_id": user.EpunjabID,
	})
}

// GET /api/teachers — list all teachers (for admin teacher creation page)
func GetTeachers(c *fiber.Ctx) error {
	var teachers []models.Teacher
	database.DB.Where("school_id = ? AND is_active = true", schoolID()).Order("name asc").Find(&teachers)
	return c.JSON(teachers)
}

// POST /api/teachers — create teacher + user account
func CreateTeacher(c *fiber.Ctx) error {
	if middleware.GetRole(c) != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}

	var body struct {
		Name          string `json:"name"`
		Phone         string `json:"phone"`
		Email         string `json:"email"`
		EmployeeID    string `json:"employee_id"`
		Subject       string `json:"subject"`
		Qualification string `json:"qualification"`
		Username      string `json:"username"`
		Password      string `json:"password"`
		EpunjabID     string `json:"epunjab_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if body.Name == "" || body.Username == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name, username and password are required"})
	}

	// create teacher record
	t := models.Teacher{
		SchoolID:      schoolID(),
		Name:          body.Name,
		Phone:         body.Phone,
		Email:         body.Email,
		EmployeeID:    body.EmployeeID,
		Subject:       body.Subject,
		Qualification: body.Qualification,
		IsActive:      true,
	}
	if err := database.DB.Create(&t).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// create user account
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	user := models.User{
		SchoolID:  schoolID(),
		Username:  body.Username,
		Password:  string(hash),
		Role:      "teacher",
		EpunjabID: body.EpunjabID,
		TeacherID: &t.ID,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		// rollback teacher creation
		database.DB.Delete(&t)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "username already exists"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "teacher created",
		"teacher_id": t.ID,
		"user_id":    user.ID,
	})
}

// DELETE /api/teachers/:id
func DeleteTeacher(c *fiber.Ctx) error {
	if middleware.GetRole(c) != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}
	database.DB.Model(&models.Teacher{}).Where("id = ? AND school_id = ?", c.Params("id"), schoolID()).
		Update("is_active", false)
	// also deactivate user
	database.DB.Model(&models.User{}).Where("teacher_id = ?", c.Params("id")).Delete(&models.User{})
	return c.JSON(fiber.Map{"message": "teacher deactivated"})
}
