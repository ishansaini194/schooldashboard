package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/middleware"
	"github.com/ishansaini194/dashboard/models"
	"github.com/ishansaini194/dashboard/validators"
)

// build a StudentInput view of a Student for validation
func toStudentInput(s models.Student) validators.StudentInput {
	return validators.StudentInput{
		Name:          s.Name,
		Class:         s.Class,
		Section:       s.Section,
		RollNo:        s.RollNo,
		EpunjabId:     s.EpunjabId,
		Phone:         s.Phone,
		AadharNo:      s.AadharNo,
		FatherName:    s.FatherName,
		FatherContact: s.FatherContact,
		FatherAadhar:  s.FatherAadhar,
		MotherName:    s.MotherName,
		MotherContact: s.MotherContact,
		Caste:         s.Caste,
		Gender:        s.Gender,
		DOB:           s.DOB,
		Address:       s.Address,
	}
}

// GET /api/students/class/:class?section=A
func GetStudents(c *fiber.Ctx) error {
	class := c.Params("class")
	section := c.Query("section")

	q := database.DB.Where("class = ?", class)
	if section != "" {
		q = q.Where("section = ?", section)
	}

	var students []models.Student
	result := q.Order("CAST(roll_no AS INTEGER) asc").Find(&students)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.JSON(students)
}

// GET /api/students/:roll_no?class=6&section=A
// roll_no is not unique across classes — pass class+section to disambiguate.
// If only roll_no is given, returns the first match (legacy behaviour with a warning).
func GetStudent(c *fiber.Ctx) error {
	rollNo := c.Params("roll_no")
	class := c.Query("class")
	section := c.Query("section")

	q := database.DB.Where("roll_no = ?", rollNo)
	if class != "" {
		q = q.Where("class = ?", class)
	}
	if section != "" {
		q = q.Where("section = ?", section)
	}

	var s models.Student
	result := q.First(&s)
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

	// validate
	if errs := validators.ValidateStudent(toStudentInput(s)); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "validation failed",
			"fields": errs,
		})
	}

	result := database.DB.Create(&s)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "student created"})
}

func UpdateStudent(c *fiber.Ctx) error {
	rollNo := c.Params("roll_no")
	class := c.Query("class")
	section := c.Query("section")

	// first find the existing student
	q := database.DB.Where("roll_no = ?", rollNo)
	if class != "" {
		q = q.Where("class = ?", class)
	}
	if section != "" {
		q = q.Where("section = ?", section)
	}

	var s models.Student
	result := q.First(&s)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}

	// parse body into existing student
	if err := c.BodyParser(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// validate
	if errs := validators.ValidateStudent(toStudentInput(s)); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "validation failed",
			"fields": errs,
		})
	}

	// save updates all fields
	result = database.DB.Save(&s)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"message": "student updated"})
}

func DeleteStudent(c *fiber.Ctx) error {
	if middleware.GetRole(c) != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}

	rollNo := c.Params("roll_no")
	class := c.Query("class")
	section := c.Query("section")

	q := database.DB.Where("roll_no = ?", rollNo)
	if class != "" {
		q = q.Where("class = ?", class)
	}
	if section != "" {
		q = q.Where("section = ?", section)
	}

	result := q.Delete(&models.Student{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}

	return c.JSON(fiber.Map{"message": "student deleted"})
}

// GET /api/students/epunjab/:epunjab_id
func GetStudentByEpunjabID(c *fiber.Ctx) error {
	epunjabID := c.Params("epunjab_id")

	var s models.Student
	result := database.DB.Where("epunjab_id = ?", epunjabID).First(&s)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}

	return c.JSON(s)
}
