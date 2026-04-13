package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/middleware"
	"github.com/ishansaini194/dashboard/models/academic"
)

// POST /api/results
func CreateResult(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	var r academic.Result
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result := database.DB.Create(&r)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "result added"})
}

// GET /api/results/student/:student_id?exam_type=midterm&year=2026
func GetStudentResults(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	examType := c.Query("exam_type", "")
	year := c.Query("year", "")

	query := database.DB.Where("student_id = ?", studentID)

	if examType != "" {
		query = query.Where("exam_type = ?", examType)
	}
	if year != "" {
		query = query.Where("year = ?", year)
	}

	var results []academic.Result
	if err := query.Order("subject asc").Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// GET /api/results/class/:class/section/:section?exam_type=midterm&year=2026
func GetClassResults(c *fiber.Ctx) error {
	class := c.Params("class")
	section := c.Params("section")
	examType := c.Query("exam_type", "")
	year := c.Query("year", "")

	query := database.DB.Where("class = ? AND section = ?", class, section)

	if examType != "" {
		query = query.Where("exam_type = ?", examType)
	}
	if year != "" {
		query = query.Where("year = ?", year)
	}

	var results []academic.Result
	if err := query.Order("subject asc").Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// PUT /api/results/:id
func UpdateResult(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	id := c.Params("id")

	var r academic.Result
	if err := database.DB.First(&r, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "result not found"})
	}

	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	database.DB.Save(&r)
	return c.JSON(fiber.Map{"message": "result updated"})
}

// DELETE /api/results/:id
func DeleteResult(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not allowed"})
	}

	id := c.Params("id")

	if err := database.DB.Delete(&academic.Result{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "result deleted"})
}
