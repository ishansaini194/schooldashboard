package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/middleware"
	"github.com/ishansaini194/dashboard/models/academic"
)

// POST /api/homework
func CreateHomework(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	var h academic.Homework
	if err := c.BodyParser(&h); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := database.DB.Create(&h).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "homework added"})
}

// GET /api/homework/class/:class/section/:section
func GetHomework(c *fiber.Ctx) error {
	class := c.Params("class")
	section := c.Params("section")

	var homework []academic.Homework
	if err := database.DB.
		Where("class = ? AND section = ?", class, section).
		Order("created_at desc").
		Limit(30).
		Find(&homework).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(homework)
}

// GET /api/homework/:id
func GetHomeworkByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var h academic.Homework
	if err := database.DB.First(&h, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "homework not found"})
	}

	return c.JSON(h)
}

// PUT /api/homework/:id
func UpdateHomework(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	id := c.Params("id")

	var h academic.Homework
	if err := database.DB.First(&h, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "homework not found"})
	}

	if err := c.BodyParser(&h); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	database.DB.Save(&h)
	return c.JSON(fiber.Map{"message": "homework updated"})
}

// DELETE /api/homework/:id
func DeleteHomework(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	id := c.Params("id")

	if err := database.DB.Delete(&academic.Homework{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "homework deleted"})
}
