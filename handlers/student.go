package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/models"
)

func GetStudents(c *fiber.Ctx) error {
	class := c.Params("class")

	var students []models.Student
	result := database.DB.Where("class = ?", class).Order("roll_no asc").Find(&students)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.JSON(students)
}

func GetStudent(c *fiber.Ctx) error {
	rollNo := c.Params("roll_no")

	var s models.Student
	result := database.DB.Where("roll_no = ?", rollNo).First(&s)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}

	return c.JSON(s)
}

func CreateStudent(c *fiber.Ctx) error {
	var s models.Student
	if err := c.BodyParser(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result := database.DB.Create(&s)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "student created"})
}

func UpdateStudent(c *fiber.Ctx) error {
	rollNo := c.Params("roll_no")

	// first find the existing student
	var s models.Student
	result := database.DB.Where("roll_no = ?", rollNo).First(&s)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}

	// parse body into existing student
	if err := c.BodyParser(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// save updates all fields
	result = database.DB.Save(&s)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"message": "student updated"})
}

func DeleteStudent(c *fiber.Ctx) error {
	rollNo := c.Params("roll_no")

	result := database.DB.Where("roll_no = ?", rollNo).Delete(&models.Student{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"message": "student deleted"})
}
