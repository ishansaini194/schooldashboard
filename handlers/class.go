package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/models"
)

// GET /api/classes
func GetClasses(c *fiber.Ctx) error {
	var classes []models.Class
	database.DB.Where("school_id = ?", schoolID()).Order("number asc, section asc").Find(&classes)

	// ClassResponse keeps backward-compat field name "class" (was models.Class.Class int)
	// now models.Class.Number int — we alias it back so frontend doesn't change
	type ClassResponse struct {
		ID             uint   `json:"id"`
		Class          int    `json:"class"` // alias of Number — keeps frontend working
		Number         int    `json:"number"`
		Section        string `json:"section"`
		TeacherName    string `json:"teacher_name"`
		TeacherContact string `json:"teacher_contact"`
		TuitionFee     int    `json:"tuition_fee"`
		TransportFee   int    `json:"transport_fee"`
	}

	result := make([]ClassResponse, 0, len(classes))
	for _, cls := range classes {
		cr := ClassResponse{
			ID: cls.ID, Class: cls.Number, Number: cls.Number,
			Section: cls.Section, TuitionFee: cls.TuitionFee, TransportFee: cls.TransportFee,
		}
		if cls.ClassTeacherID != nil {
			var t models.Teacher
			database.DB.First(&t, cls.ClassTeacherID)
			cr.TeacherName = t.Name
			cr.TeacherContact = t.Phone
		}
		result = append(result, cr)
	}
	return c.JSON(result)
}

// GET /api/classes/:id
func GetClass(c *fiber.Ctx) error {
	var cls models.Class
	database.DB.Where("id = ? AND school_id = ?", c.Params("id"), schoolID()).First(&cls)
	if cls.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "class not found"})
	}
	return c.JSON(cls)
}

// POST /api/classes
func CreateClass(c *fiber.Ctx) error {
	type Body struct {
		Class          int    `json:"class"`   // backward compat field name
		Number         int    `json:"number"`  // new field name
		Section        string `json:"section"`
		TeacherName    string `json:"teacher_name"`
		TeacherContact string `json:"teacher_contact"`
		TuitionFee     int    `json:"tuition_fee"`
		TransportFee   int    `json:"transport_fee"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// accept either field name
	num := body.Number
	if num == 0 {
		num = body.Class
	}

	// upsert teacher if name provided
	var teacherID *uint
	if body.TeacherName != "" {
		t := models.Teacher{
			SchoolID: schoolID(), Name: body.TeacherName, Phone: body.TeacherContact, IsActive: true,
		}
		database.DB.Where(models.Teacher{SchoolID: schoolID(), Name: body.TeacherName}).FirstOrCreate(&t)
		if body.TeacherContact != "" {
			t.Phone = body.TeacherContact
			database.DB.Save(&t)
		}
		teacherID = &t.ID
	}

	cls := models.Class{
		SchoolID: schoolID(), Number: num, Section: body.Section,
		ClassTeacherID: teacherID, TuitionFee: body.TuitionFee, TransportFee: body.TransportFee,
	}
	if err := database.DB.Create(&cls).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": cls.ID, "class": cls.Number, "section": cls.Section})
}

// PUT /api/classes/:id
func UpdateClass(c *fiber.Ctx) error {
	var cls models.Class
	database.DB.Where("id = ? AND school_id = ?", c.Params("id"), schoolID()).First(&cls)
	if cls.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "class not found"})
	}

	type Body struct {
		Class          int    `json:"class"`
		Number         int    `json:"number"`
		Section        string `json:"section"`
		TeacherName    string `json:"teacher_name"`
		TeacherContact string `json:"teacher_contact"`
		TuitionFee     int    `json:"tuition_fee"`
		TransportFee   int    `json:"transport_fee"`
	}
	var body Body
	c.BodyParser(&body)

	num := body.Number
	if num == 0 {
		num = body.Class
	}

	var teacherID *uint
	if body.TeacherName != "" {
		t := models.Teacher{SchoolID: schoolID(), Name: body.TeacherName, IsActive: true}
		database.DB.Where(models.Teacher{SchoolID: schoolID(), Name: body.TeacherName}).FirstOrCreate(&t)
		if body.TeacherContact != "" {
			t.Phone = body.TeacherContact
			database.DB.Save(&t)
		}
		teacherID = &t.ID
	}

	if num != 0 {
		cls.Number = num
	}
	cls.Section = body.Section
	cls.ClassTeacherID = teacherID
	cls.TuitionFee = body.TuitionFee
	cls.TransportFee = body.TransportFee
	database.DB.Save(&cls)

	return c.JSON(fiber.Map{"message": "class updated"})
}

// DELETE /api/classes/:id
func DeleteClass(c *fiber.Ctx) error {
	result := database.DB.Where("id = ? AND school_id = ?", c.Params("id"), schoolID()).Delete(&models.Class{})
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "class not found"})
	}
	return c.JSON(fiber.Map{"message": "class deleted"})
}
