package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/middleware"
	"github.com/ishansaini194/dashboard/models/academic"
)

// POST /api/notices
func CreateNotice(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	var n academic.Notice
	if err := c.BodyParser(&n); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := database.DB.Create(&n).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "notice created"})
}

// GET /api/notices?target=all&target=8-A
func GetNotices(c *fiber.Ctx) error {
	target := c.Query("target", "")

	query := database.DB.Order("created_at desc")

	if target != "" {
		query = query.Where("target = ? OR target = ?", "all", target)
	}

	var notices []academic.Notice
	if err := query.Limit(50).Find(&notices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(notices)
}

// GET /api/notices/:id
func GetNoticeByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var n academic.Notice
	if err := database.DB.First(&n, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "notice not found"})
	}

	return c.JSON(n)
}

// PUT /api/notices/:id
func UpdateNotice(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	id := c.Params("id")

	var n academic.Notice
	if err := database.DB.First(&n, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "notice not found"})
	}

	if err := c.BodyParser(&n); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	database.DB.Save(&n)
	return c.JSON(fiber.Map{"message": "notice updated"})
}

// DELETE /api/notices/:id
func DeleteNotice(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	id := c.Params("id")

	if err := database.DB.Delete(&academic.Notice{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "notice deleted"})
}
