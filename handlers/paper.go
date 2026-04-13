package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/middleware"
	"github.com/ishansaini194/dashboard/models/academic"
)

// POST /api/papers
func CreatePaper(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	var p academic.Paper
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := database.DB.Create(&p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "paper added"})
}

// GET /api/papers?class=8&section=A&subject=Math
func GetPapers(c *fiber.Ctx) error {
	class := c.Query("class", "")
	section := c.Query("section", "")
	subject := c.Query("subject", "")

	query := database.DB.Order("year desc")

	if class != "" {
		query = query.Where("class = ?", class)
	}
	if section != "" {
		query = query.Where("section = ?", section)
	}
	if subject != "" {
		query = query.Where("subject = ?", subject)
	}

	var papers []academic.Paper
	if err := query.Find(&papers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(papers)
}

// GET /api/papers/:id
func GetPaperByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var p academic.Paper
	if err := database.DB.First(&p, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "paper not found"})
	}

	return c.JSON(p)
}

// PUT /api/papers/:id
func UpdatePaper(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	id := c.Params("id")

	var p academic.Paper
	if err := database.DB.First(&p, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "paper not found"})
	}

	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	database.DB.Save(&p)
	return c.JSON(fiber.Map{"message": "paper updated"})
}

// DELETE /api/papers/:id
func DeletePaper(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	id := c.Params("id")

	if err := database.DB.Delete(&academic.Paper{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "paper deleted"})
}
