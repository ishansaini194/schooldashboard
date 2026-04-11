package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/models"
)

func GetClasses(c *fiber.Ctx) error {
	var classes []models.Class
	result := database.DB.Find(&classes)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	database.DB.Order("class asc").Find(&classes)
	return c.JSON(classes)
}

func GetClass(c *fiber.Ctx) error {
	id := c.Params("id")

	var class models.Class
	result := database.DB.First(&class, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "class not found"})
	}
	return c.JSON(class)
}

func CreateClass(c *fiber.Ctx) error {
	var class models.Class
	if err := c.BodyParser(&class); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result := database.DB.Create(&class)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "class created"})
}

func UpdateClass(c *fiber.Ctx) error {
	id := c.Params("id")

	var class models.Class
	if err := c.BodyParser(&class); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result := database.DB.Model(&models.Class{}).Where("id = ?", id).Updates(&class)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.JSON(fiber.Map{"message": "class updated"})
}

func DeleteClass(c *fiber.Ctx) error {
	id := c.Params("id")

	result := database.DB.Delete(&models.Class{}, id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.JSON(fiber.Map{"message": "class deleted"})
}
