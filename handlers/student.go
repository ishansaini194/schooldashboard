package handlers

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/middleware"
	"github.com/ishansaini194/dashboard/models"
	"github.com/ishansaini194/dashboard/validators"
)

// ── Helpers ───────────────────────────────────────────────────

func schoolID() uint {
	id, _ := strconv.Atoi(os.Getenv("SCHOOL_ID"))
	if id == 0 {
		return 1
	}
	return uint(id)
}

func currentAcademicYearID() uint {
	var ay models.AcademicYear
	database.DB.Where("school_id = ? AND is_current = true", schoolID()).First(&ay)
	return ay.ID
}

// StudentResponse — same shape as before so frontend doesn't break
type StudentResponse struct {
	ID                    uint   `json:"id"`
	Name                  string `json:"name"`
	Class                 string `json:"class"`
	Section               string `json:"section"`
	RollNo                string `json:"roll_no"`
	Phone                 string `json:"phone"`
	Gender                string `json:"gender"`
	DOB                   string `json:"dob"`
	AadharNo              string `json:"aadhar_no"`
	EpunjabId             string `json:"epunjab_id"`
	FatherName            string `json:"father_name"`
	FatherContact         string `json:"father_contact"`
	FatherAadhar          string `json:"father_aadhar"`
	MotherName            string `json:"mother_name"`
	MotherContact         string `json:"mother_contact"`
	Address               string `json:"address"`
	Caste                 string `json:"caste"`
	PreviousSchoolDetails string `json:"previous_school_details"`
	EnrollmentID          uint   `json:"enrollment_id"`
	ClassID               uint   `json:"class_id"`
}

func toStudentResponse(s models.Student, e models.Enrollment, cls models.Class) StudentResponse {
	return StudentResponse{
		ID: s.ID, Name: s.Name,
		Class: strconv.Itoa(cls.Number), Section: cls.Section,
		RollNo: strconv.Itoa(e.RollNo),
		Phone:  s.Phone, Gender: s.Gender, DOB: s.DOB,
		AadharNo: s.AadharNo, EpunjabId: s.EpunjabId,
		FatherName: s.FatherName, FatherContact: s.FatherContact, FatherAadhar: s.FatherAadhar,
		MotherName: s.MotherName, MotherContact: s.MotherContact,
		Address: s.Address, Caste: s.Caste,
		PreviousSchoolDetails: s.PreviousSchoolDetails,
		EnrollmentID:          e.ID, ClassID: cls.ID,
	}
}

func findClass(classNum, section string) models.Class {
	var cls models.Class
	q := database.DB.Where("school_id = ? AND number = ?", schoolID(), classNum)
	if section != "" {
		q = q.Where("section = ?", section)
	}
	q.First(&cls)
	return cls
}

func findEnrollmentByRoll(rollInt int, classID uint) models.Enrollment {
	var e models.Enrollment
	ayID := currentAcademicYearID()
	q := database.DB.Where("roll_no = ? AND academic_year_id = ?", rollInt, ayID)
	if classID != 0 {
		q = q.Where("class_id = ?", classID)
	}
	q.First(&e)
	return e
}

// ── Handlers ──────────────────────────────────────────────────

// GET /api/students/class/:class?section=A
func GetStudents(c *fiber.Ctx) error {
	cls := findClass(c.Params("class"), c.Query("section"))
	if cls.ID == 0 {
		return c.JSON([]StudentResponse{})
	}

	ayID := currentAcademicYearID()
	var enrollments []models.Enrollment
	database.DB.Where("class_id = ? AND academic_year_id = ? AND status = 'active'", cls.ID, ayID).
		Order("roll_no asc").Find(&enrollments)

	result := make([]StudentResponse, 0, len(enrollments))
	for _, e := range enrollments {
		var s models.Student
		database.DB.First(&s, e.StudentID)
		result = append(result, toStudentResponse(s, e, cls))
	}
	return c.JSON(result)
}

// GET /api/students/:roll_no?class=6&section=A
func GetStudent(c *fiber.Ctx) error {
	rollInt, _ := strconv.Atoi(c.Params("roll_no"))
	cls := findClass(c.Query("class"), c.Query("section"))
	e := findEnrollmentByRoll(rollInt, cls.ID)
	if e.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}
	if cls.ID == 0 {
		database.DB.First(&cls, e.ClassID)
	}
	var s models.Student
	database.DB.First(&s, e.StudentID)
	return c.JSON(toStudentResponse(s, e, cls))
}

// GET /api/students/epunjab/:epunjab_id
func GetStudentByEpunjabID(c *fiber.Ctx) error {
	var s models.Student
	database.DB.Where("epunjab_id = ? AND school_id = ?", c.Params("epunjab_id"), schoolID()).First(&s)
	if s.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}
	var e models.Enrollment
	database.DB.Where("student_id = ? AND academic_year_id = ?", s.ID, currentAcademicYearID()).First(&e)
	var cls models.Class
	database.DB.First(&cls, e.ClassID)
	return c.JSON(toStudentResponse(s, e, cls))
}

// POST /api/students
func CreateStudent(c *fiber.Ctx) error {
	type Body struct {
		Name                  string `json:"name"`
		Class                 string `json:"class"`
		Section               string `json:"section"`
		RollNo                string `json:"roll_no"`
		Phone                 string `json:"phone"`
		Gender                string `json:"gender"`
		DOB                   string `json:"dob"`
		AadharNo              string `json:"aadhar_no"`
		EpunjabId             string `json:"epunjab_id"`
		FatherName            string `json:"father_name"`
		FatherContact         string `json:"father_contact"`
		FatherAadhar          string `json:"father_aadhar"`
		MotherName            string `json:"mother_name"`
		MotherContact         string `json:"mother_contact"`
		Address               string `json:"address"`
		Caste                 string `json:"caste"`
		PreviousSchoolDetails string `json:"previous_school_details"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	vi := validators.StudentInput{
		Name: body.Name, Class: body.Class, Section: body.Section, RollNo: body.RollNo,
		EpunjabId: body.EpunjabId, Phone: body.Phone, AadharNo: body.AadharNo,
		FatherName: body.FatherName, FatherContact: body.FatherContact, FatherAadhar: body.FatherAadhar,
		MotherName: body.MotherName, MotherContact: body.MotherContact,
		Caste: body.Caste, Gender: body.Gender, DOB: body.DOB, Address: body.Address,
	}
	if errs := validators.ValidateStudent(vi); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "fields": errs})
	}

	cls := findClass(body.Class, body.Section)
	if cls.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "class not found — create the class first"})
	}

	s := models.Student{
		SchoolID: schoolID(), Name: body.Name, Phone: body.Phone, Gender: body.Gender, DOB: body.DOB,
		AadharNo: body.AadharNo, EpunjabId: body.EpunjabId,
		FatherName: body.FatherName, FatherContact: body.FatherContact, FatherAadhar: body.FatherAadhar,
		MotherName: body.MotherName, MotherContact: body.MotherContact,
		Address: body.Address, Caste: body.Caste, PreviousSchoolDetails: body.PreviousSchoolDetails,
		IsActive: true,
	}
	if err := database.DB.Create(&s).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	rollInt, _ := strconv.Atoi(body.RollNo)
	e := models.Enrollment{
		StudentID: s.ID, ClassID: cls.ID,
		AcademicYearID: currentAcademicYearID(),
		RollNo:         rollInt, Status: "active",
	}
	database.DB.Create(&e)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "student created", "id": s.ID})
}

// PUT /api/students/:roll_no?class=6&section=A
func UpdateStudent(c *fiber.Ctx) error {
	rollInt, _ := strconv.Atoi(c.Params("roll_no"))
	cls := findClass(c.Query("class"), c.Query("section"))
	e := findEnrollmentByRoll(rollInt, cls.ID)
	if e.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}

	type Body struct {
		Name                  string `json:"name"`
		Class                 string `json:"class"`
		Section               string `json:"section"`
		RollNo                string `json:"roll_no"`
		Phone                 string `json:"phone"`
		Gender                string `json:"gender"`
		DOB                   string `json:"dob"`
		AadharNo              string `json:"aadhar_no"`
		EpunjabId             string `json:"epunjab_id"`
		FatherName            string `json:"father_name"`
		FatherContact         string `json:"father_contact"`
		FatherAadhar          string `json:"father_aadhar"`
		MotherName            string `json:"mother_name"`
		MotherContact         string `json:"mother_contact"`
		Address               string `json:"address"`
		Caste                 string `json:"caste"`
		PreviousSchoolDetails string `json:"previous_school_details"`
	}
	var body Body
	c.BodyParser(&body)

	vi := validators.StudentInput{
		Name: body.Name, Class: body.Class, Section: body.Section, RollNo: body.RollNo,
		EpunjabId: body.EpunjabId, Phone: body.Phone, AadharNo: body.AadharNo,
		FatherName: body.FatherName, FatherContact: body.FatherContact, FatherAadhar: body.FatherAadhar,
		MotherName: body.MotherName, MotherContact: body.MotherContact,
		Caste: body.Caste, Gender: body.Gender, DOB: body.DOB, Address: body.Address,
	}
	if errs := validators.ValidateStudent(vi); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "fields": errs})
	}

	var s models.Student
	database.DB.First(&s, e.StudentID)
	s.Name = body.Name
	s.Phone = body.Phone
	s.Gender = body.Gender
	s.DOB = body.DOB
	s.AadharNo = body.AadharNo
	s.EpunjabId = body.EpunjabId
	s.FatherName = body.FatherName
	s.FatherContact = body.FatherContact
	s.FatherAadhar = body.FatherAadhar
	s.MotherName = body.MotherName
	s.MotherContact = body.MotherContact
	s.Address = body.Address
	s.Caste = body.Caste
	s.PreviousSchoolDetails = body.PreviousSchoolDetails
	database.DB.Save(&s)

	if body.Class != "" && body.Section != "" {
		newCls := findClass(body.Class, body.Section)
		if newCls.ID != 0 {
			e.ClassID = newCls.ID
		}
	}
	if body.RollNo != "" {
		newRoll, _ := strconv.Atoi(body.RollNo)
		e.RollNo = newRoll
	}
	database.DB.Save(&e)

	return c.JSON(fiber.Map{"message": "student updated"})
}

// DELETE /api/students/:roll_no?class=6&section=A
func DeleteStudent(c *fiber.Ctx) error {
	if middleware.GetRole(c) != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}
	rollInt, _ := strconv.Atoi(c.Params("roll_no"))
	cls := findClass(c.Query("class"), c.Query("section"))
	e := findEnrollmentByRoll(rollInt, cls.ID)
	if e.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}
	database.DB.Delete(&e)
	database.DB.Delete(&models.Student{}, e.StudentID)
	return c.JSON(fiber.Map{"message": "student deleted"})
}
