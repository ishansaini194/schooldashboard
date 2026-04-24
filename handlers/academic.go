package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/middleware"
	"github.com/ishansaini194/dashboard/models"
)

// ── Homework ──────────────────────────────────────────────────

func CreateHomework(c *fiber.Ctx) error {
	role := middleware.GetRole(c)
	if role == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "students cannot post homework"})
	}
	var body struct {
		ClassID  uint   `json:"class_id"`
		ClassNum string `json:"class"`
		Section  string `json:"section"`
		Subject  string `json:"subject"`
		Content  string `json:"content"`
		DueDate  string `json:"due_date"`
	}
	c.BodyParser(&body)

	// resolve class_id if not provided directly
	classID := body.ClassID
	if classID == 0 {
		cls := findClass(body.ClassNum, body.Section)
		classID = cls.ID
	}

	hw := models.Homework{
		ClassID:   classID,
		TeacherID: middleware.GetTeacherID(c),
		Subject:   body.Subject,
		Content:   body.Content,
		DueDate:   body.DueDate,
	}
	database.DB.Create(&hw)
	return c.Status(fiber.StatusCreated).JSON(hw)
}

func UpdateHomework(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	var hw models.Homework
	database.DB.First(&hw, c.Params("id"))
	c.BodyParser(&hw)
	database.DB.Save(&hw)
	return c.JSON(hw)
}

func DeleteHomework(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	database.DB.Delete(&models.Homework{}, c.Params("id"))
	return c.JSON(fiber.Map{"message": "deleted"})
}

// GET /api/homework/class/:class/section/:section
func GetHomework(c *fiber.Ctx) error {
	cls := findClass(c.Params("class"), c.Params("section"))

	type HomeworkResponse struct {
		models.Homework
		Class   string `json:"class"`
		Section string `json:"section"`
	}

	var hws []models.Homework
	database.DB.Where("class_id = ?", cls.ID).Order("created_at desc").Find(&hws)

	result := make([]HomeworkResponse, 0, len(hws))
	for _, hw := range hws {
		result = append(result, HomeworkResponse{
			Homework: hw,
			Class:    c.Params("class"),
			Section:  c.Params("section"),
		})
	}
	return c.JSON(result)
}

func GetHomeworkByID(c *fiber.Ctx) error {
	var hw models.Homework
	database.DB.First(&hw, c.Params("id"))
	return c.JSON(hw)
}

// ── Notices ───────────────────────────────────────────────────

func CreateNotice(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	var body struct {
		Title         string `json:"title"`
		Body          string `json:"body"`
		Target        string `json:"target"` // "all" or "6-A"
		TargetClassID *uint  `json:"target_class_id"`
	}
	c.BodyParser(&body)

	targetType := "all"
	var targetClassID *uint

	if body.Target != "all" && body.Target != "" {
		targetType = "class"
		targetClassID = body.TargetClassID
	}

	n := models.Notice{
		SchoolID:      schoolID(),
		PostedBy:      middleware.GetUserID(c),
		Title:         body.Title,
		Body:          body.Body,
		TargetType:    targetType,
		TargetClassID: targetClassID,
	}
	database.DB.Create(&n)
	return c.Status(fiber.StatusCreated).JSON(n)
}

func UpdateNotice(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	var n models.Notice
	database.DB.First(&n, c.Params("id"))
	c.BodyParser(&n)
	database.DB.Save(&n)
	return c.JSON(n)
}

func DeleteNotice(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	database.DB.Delete(&models.Notice{}, c.Params("id"))
	return c.JSON(fiber.Map{"message": "deleted"})
}

// GET /api/notices?target=6-A
func GetNotices(c *fiber.Ctx) error {
	target := c.Query("target")

	type NoticeResponse struct {
		models.Notice
		Target   string `json:"target"`
		PostedBy string `json:"posted_by"` // username string for backward compat
	}

	var notices []models.Notice
	q := database.DB.Where("school_id = ?", schoolID())
	if target != "" && target != "all" {
		q = q.Where("target_type = 'all' OR (target_type = 'class' AND target_class_id IN (SELECT id FROM classes WHERE CONCAT(number, '-', section) = ?))", target)
	}
	q.Order("created_at desc").Find(&notices)

	result := make([]NoticeResponse, 0, len(notices))
	for _, n := range notices {
		var u models.User
		database.DB.First(&u, n.PostedBy)
		targetStr := "all"
		if n.TargetClassID != nil {
			var cls models.Class
			database.DB.First(&cls, n.TargetClassID)
			targetStr = fmt.Sprintf("%d-%s", cls.Number, cls.Section)
		}
		result = append(result, NoticeResponse{
			Notice: n, Target: targetStr, PostedBy: u.Username,
		})
	}
	return c.JSON(result)
}

func GetNoticeByID(c *fiber.Ctx) error {
	var n models.Notice
	database.DB.First(&n, c.Params("id"))
	return c.JSON(n)
}

// ── Papers ────────────────────────────────────────────────────

func CreatePaper(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	var body struct {
		ClassID   uint   `json:"class_id"`
		ClassNum  string `json:"class"`
		Section   string `json:"section"`
		Subject   string `json:"subject"`
		ExamType  string `json:"exam_type"`
		Year      int    `json:"year"`
		DriveLink string `json:"drive_link"`
	}
	c.BodyParser(&body)

	classID := body.ClassID
	if classID == 0 {
		cls := findClass(body.ClassNum, body.Section)
		classID = cls.ID
	}

	p := models.Paper{
		ClassID:   classID,
		TeacherID: middleware.GetTeacherID(c),
		Subject:   body.Subject,
		ExamType:  body.ExamType,
		Year:      body.Year,
		DriveLink: body.DriveLink,
	}
	database.DB.Create(&p)
	return c.Status(fiber.StatusCreated).JSON(p)
}

func UpdatePaper(c *fiber.Ctx) error {
	var p models.Paper
	database.DB.First(&p, c.Params("id"))
	c.BodyParser(&p)
	database.DB.Save(&p)
	return c.JSON(p)
}

func DeletePaper(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	database.DB.Delete(&models.Paper{}, c.Params("id"))
	return c.JSON(fiber.Map{"message": "deleted"})
}

func GetPapers(c *fiber.Ctx) error {
	type PaperResponse struct {
		models.Paper
		Class      string `json:"class"`
		Section    string `json:"section"`
		UploadedBy string `json:"uploaded_by"`
	}

	var papers []models.Paper
	database.DB.Order("created_at desc").Find(&papers)

	result := make([]PaperResponse, 0, len(papers))
	for _, p := range papers {
		var cls models.Class
		database.DB.First(&cls, p.ClassID)
		var t models.Teacher
		database.DB.First(&t, p.TeacherID)
		result = append(result, PaperResponse{
			Paper:      p,
			Class:      fmt.Sprintf("%d", cls.Number),
			Section:    cls.Section,
			UploadedBy: t.Name,
		})
	}
	return c.JSON(result)
}

func GetPaperByID(c *fiber.Ctx) error {
	var p models.Paper
	database.DB.First(&p, c.Params("id"))
	return c.JSON(p)
}

// ── Results ───────────────────────────────────────────────────

func CreateResult(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	var body struct {
		ExamID       uint `json:"exam_id"`
		StudentID    uint `json:"student_id"`
		EnrollmentID uint `json:"enrollment_id"`
		Marks        int  `json:"marks"`
		// these help create exam if exam_id not provided
		ClassNum string `json:"class"`
		Section  string `json:"section"`
		Subject  string `json:"subject"`
		ExamType string `json:"exam_type"`
		MaxMarks int    `json:"max_marks"`
		Year     int    `json:"year"`
	}
	c.BodyParser(&body)

	// resolve enrollment_id
	enrollmentID := body.EnrollmentID
	if enrollmentID == 0 && body.StudentID != 0 {
		var e models.Enrollment
		database.DB.Where("student_id = ? AND academic_year_id = ?", body.StudentID, currentAcademicYearID()).First(&e)
		enrollmentID = e.ID
	}

	// resolve or create exam
	examID := body.ExamID
	if examID == 0 {
		cls := findClass(body.ClassNum, body.Section)
		var exam models.Exam
		database.DB.Where("class_id = ? AND name = ? AND subject = ? AND academic_year_id = ?",
			cls.ID, body.ExamType, body.Subject, currentAcademicYearID()).First(&exam)
		if exam.ID == 0 {
			exam = models.Exam{
				ClassID:        cls.ID,
				AcademicYearID: currentAcademicYearID(),
				Name:           body.ExamType,
				Subject:        body.Subject,
				MaxMarks:       body.MaxMarks,
			}
			database.DB.Create(&exam)
		}
		examID = exam.ID
	}

	// upsert result
	var result models.Result
	database.DB.Where("exam_id = ? AND enrollment_id = ?", examID, enrollmentID).First(&result)
	if result.ID == 0 {
		result = models.Result{
			ExamID:       examID,
			EnrollmentID: enrollmentID,
			Marks:        body.Marks,
			EnteredBy:    middleware.GetUserID(c),
		}
		database.DB.Create(&result)
	} else {
		result.Marks = body.Marks
		result.EnteredBy = middleware.GetUserID(c)
		database.DB.Save(&result)
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

func UpdateResult(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	var r models.Result
	database.DB.First(&r, c.Params("id"))
	c.BodyParser(&r)
	database.DB.Save(&r)
	return c.JSON(r)
}

func DeleteResult(c *fiber.Ctx) error {
	if middleware.GetRole(c) == "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	database.DB.Delete(&models.Result{}, c.Params("id"))
	return c.JSON(fiber.Map{"message": "deleted"})
}

func GetStudentResults(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	var e models.Enrollment
	database.DB.Where("student_id = ? AND academic_year_id = ?", studentID, currentAcademicYearID()).First(&e)

	type ResultResponse struct {
		ID       uint   `json:"id"`
		Subject  string `json:"subject"`
		ExamType string `json:"exam_type"`
		Marks    int    `json:"marks"`
		MaxMarks int    `json:"max_marks"`
		Year     int    `json:"year"`
		Class    string `json:"class"`
		Section  string `json:"section"`
	}

	var results []models.Result
	database.DB.Where("enrollment_id = ?", e.ID).Find(&results)

	response := make([]ResultResponse, 0, len(results))
	for _, r := range results {
		var exam models.Exam
		database.DB.First(&exam, r.ExamID)
		var ay models.AcademicYear
		database.DB.First(&ay, exam.AcademicYearID)
		var cls models.Class
		database.DB.First(&cls, exam.ClassID)
		// Parse year from academic year name (e.g., "2025-26" -> 2025)
		year := 0
		if len(ay.Name) >= 4 {
			fmt.Sscanf(ay.Name[:4], "%d", &year)
		}
		response = append(response, ResultResponse{
			ID: r.ID, Subject: exam.Subject, ExamType: exam.Name,
			Marks: r.Marks, MaxMarks: exam.MaxMarks, Year: year,
			Class: fmt.Sprintf("%d", cls.Number), Section: cls.Section,
		})
	}
	return c.JSON(response)
}

func GetClassResults(c *fiber.Ctx) error {
	cls := findClass(c.Params("class"), c.Params("section"))
	ayID := currentAcademicYearID()

	// get all exams for this class
	var exams []models.Exam
	database.DB.Where("class_id = ? AND academic_year_id = ?", cls.ID, ayID).Find(&exams)

	type StudentResult struct {
		StudentID   uint   `json:"student_id"`
		StudentName string `json:"student_name"`
		RollNo      string `json:"roll_no"`
		Subject     string `json:"subject"`
		ExamType    string `json:"exam_type"`
		Marks       int    `json:"marks"`
		MaxMarks    int    `json:"max_marks"`
	}

	var response []StudentResult
	for _, exam := range exams {
		var results []models.Result
		database.DB.Where("exam_id = ?", exam.ID).Find(&results)
		for _, r := range results {
			var e models.Enrollment
			database.DB.First(&e, r.EnrollmentID)
			var s models.Student
			database.DB.First(&s, e.StudentID)
			response = append(response, StudentResult{
				StudentID: s.ID, StudentName: s.Name,
				RollNo:  fmt.Sprintf("%d", e.RollNo),
				Subject: exam.Subject, ExamType: exam.Name,
				Marks: r.Marks, MaxMarks: exam.MaxMarks,
			})
		}
	}
	return c.JSON(response)
}

// GET /api/results/mine — teacher's own results
func GetMyResults(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var results []models.Result
	database.DB.Where("entered_by = ?", userID).Find(&results)
	return c.JSON(results)
}
