package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/middleware"
	"github.com/ishansaini194/dashboard/models"
)

// generate receipt number e.g. RCPT-2026-001
func generateReceiptNo() string {
	var count int64
	database.DB.Model(&models.Fee{}).Count(&count)
	year := time.Now().Year()
	return fmt.Sprintf("RCPT-%d-%03d", year, count+1)
}

// POST /api/fees/pay
func PayFee(c *fiber.Ctx) error {
	if middleware.GetRole(c) != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}

	var body struct {
		StudentID      uint   `json:"student_id"`
		EpunjabID      string `json:"epunjab_id"`
		StudentName    string `json:"student_name"`
		RollNo         string `json:"roll_no"` // add this
		Class          string `json:"class"`
		Section        string `json:"section"`
		Month          string `json:"month"`
		Year           int    `json:"year"`
		FeeType        string `json:"fee_type"`
		BaseAmount     int    `json:"base_amount"`
		Discount       int    `json:"discount"`
		DiscountReason string `json:"discount_reason"`
		PaidAmount     int    `json:"paid_amount"`
		DueDate        string `json:"due_date"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// validate first
	if body.PaidAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "paid amount must be greater than 0",
		})
	}

	finalAmount := body.BaseAmount - body.Discount
	remaining := finalAmount - body.PaidAmount

	status := "unpaid"
	if body.PaidAmount >= finalAmount {
		status = "paid"
	} else if body.PaidAmount > 0 {
		status = "partial"
	}

	fee := models.Fee{
		StudentID:      body.StudentID,
		EpunjabID:      body.EpunjabID,
		StudentName:    body.StudentName,
		RollNo:         body.RollNo, // add this
		Class:          body.Class,
		Section:        body.Section,
		Month:          body.Month,
		Year:           body.Year,
		FeeType:        body.FeeType,
		BaseAmount:     body.BaseAmount,
		Discount:       body.Discount,
		DiscountReason: body.DiscountReason,
		FinalAmount:    finalAmount,
		PaidAmount:     body.PaidAmount,
		Remaining:      remaining,
		Status:         status,
		DueDate:        body.DueDate,
		ReceiptNo:      generateReceiptNo(),
		PaidAt:         time.Now().Format("2006-01-02 15:04:05"),
	}

	result := database.DB.Create(&fee)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fee)
}

// GET /api/fees/student/:student_id
func GetStudentFees(c *fiber.Ctx) error {
	studentID := c.Params("student_id")

	var fees []models.Fee
	result := database.DB.Where("student_id = ?", studentID).
		Order("year desc, created_at desc").
		Find(&fees)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.JSON(fees)
}

// GET /api/fees/class/:class/month/:month/year/:year?section=A
func GetClassFeeStatus(c *fiber.Ctx) error {
	class := c.Params("class")
	month := c.Params("month")
	year := c.Params("year")
	section := c.Query("section")

	studentQuery := database.DB.Where("class = ?", class)
	if section != "" {
		studentQuery = studentQuery.Where("section = ?", section)
	}

	var students []models.Student
	studentQuery.Order("CAST(roll_no AS INTEGER) asc").Find(&students)

	// build student id set so we only include fees for these students
	studentIDs := make([]uint, 0, len(students))
	for _, s := range students {
		studentIDs = append(studentIDs, s.ID)
	}

	var fees []models.Fee
	if len(studentIDs) > 0 {
		database.DB.Where("class = ? AND month = ? AND year = ? AND student_id IN ?",
			class, month, year, studentIDs).Find(&fees)
	}

	feeMap := map[uint][]models.Fee{}
	for _, f := range fees {
		feeMap[f.StudentID] = append(feeMap[f.StudentID], f)
	}

	type StudentFeeStatus struct {
		StudentID   uint         `json:"student_id"`
		StudentName string       `json:"student_name"`
		RollNo      string       `json:"roll_no"`
		Fees        []models.Fee `json:"fees"`
		TotalPaid   int          `json:"total_paid"`
		HasPaid     bool         `json:"has_paid"`
	}

	var result []StudentFeeStatus
	for _, s := range students {
		studentFees := feeMap[s.ID]
		totalPaid := 0
		hasPaid := false

		for _, f := range studentFees {
			totalPaid += f.PaidAmount
			if f.Status == "paid" {
				hasPaid = true
			}
		}

		result = append(result, StudentFeeStatus{
			StudentID:   s.ID,
			StudentName: s.Name,
			RollNo:      s.RollNo,
			Fees:        studentFees,
			TotalPaid:   totalPaid,
			HasPaid:     hasPaid,
		})
	}

	return c.JSON(result)
}

// GET /api/fees/pending/:class/:month/:year?section=A
func GetPendingFees(c *fiber.Ctx) error {
	class := c.Params("class")
	month := c.Params("month")
	year := c.Params("year")
	section := c.Query("section")

	studentQuery := database.DB.Where("class = ?", class)
	if section != "" {
		studentQuery = studentQuery.Where("section = ?", section)
	}

	var students []models.Student
	studentQuery.Find(&students)

	studentIDs := make([]uint, 0, len(students))
	for _, s := range students {
		studentIDs = append(studentIDs, s.ID)
	}

	var fees []models.Fee
	if len(studentIDs) > 0 {
		database.DB.Where("class = ? AND month = ? AND year = ? AND student_id IN ?",
			class, month, year, studentIDs).Find(&fees)
	}

	paidIDs := map[uint]bool{}
	for _, f := range fees {
		if f.Status == "paid" {
			paidIDs[f.StudentID] = true
		}
	}

	type Pending struct {
		StudentID   uint   `json:"student_id"`
		StudentName string `json:"student_name"`
		RollNo      string `json:"roll_no"`
		Phone       string `json:"phone"`
	}

	var pending []Pending
	for _, s := range students {
		if !paidIDs[s.ID] {
			pending = append(pending, Pending{
				StudentID:   s.ID,
				StudentName: s.Name,
				RollNo:      s.RollNo,
				Phone:       s.Phone,
			})
		}
	}

	return c.JSON(pending)
}

// GET /api/fees/receipt/:receipt_no
func GetReceipt(c *fiber.Ctx) error {
	receiptNo := c.Params("receipt_no")

	var fee models.Fee
	result := database.DB.Where("receipt_no = ?", receiptNo).First(&fee)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "receipt not found"})
	}

	return c.JSON(fee)
}

// GET /api/fees/student/:student_id/yearly?year=2026
func GetStudentYearlySummary(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	year := c.Query("year", fmt.Sprintf("%d", time.Now().Year()))

	// get all fees for this student this year
	var fees []models.Fee
	database.DB.Where("student_id = ? AND year = ?", studentID, year).Find(&fees)

	// map by month
	feeByMonth := map[string]models.Fee{}
	for _, f := range fees {
		feeByMonth[f.Month] = f
	}

	// only show months up to current month
	allMonths := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	currentMonth := int(time.Now().Month())

	type MonthSummary struct {
		Month      string `json:"month"`
		HasRecord  bool   `json:"has_record"`
		Status     string `json:"status"`
		PaidAmount int    `json:"paid_amount"`
		Remaining  int    `json:"remaining"`
		FeeType    string `json:"fee_type"`
		ReceiptNo  string `json:"receipt_no"`
		DueDate    string `json:"due_date"`
	}

	var summary []MonthSummary
	paidCount := 0

	for i, month := range allMonths {
		if i+1 > currentMonth {
			break // don't show future months
		}
		if f, exists := feeByMonth[month]; exists {
			if f.Status == "paid" {
				paidCount++
			}
			summary = append(summary, MonthSummary{
				Month:      month,
				HasRecord:  true,
				Status:     f.Status,
				PaidAmount: f.PaidAmount,
				Remaining:  f.Remaining,
				FeeType:    f.FeeType,
				ReceiptNo:  f.ReceiptNo,
				DueDate:    f.DueDate,
			})
		} else {
			summary = append(summary, MonthSummary{
				Month:     month,
				HasRecord: false,
				Status:    "unpaid",
			})
		}
	}

	return c.JSON(fiber.Map{
		"student_id":   studentID,
		"year":         year,
		"months":       summary,
		"paid_count":   paidCount,
		"total_months": len(summary),
	})
}

// PUT /api/fees/:id/complete
func CompleteFee(c *fiber.Ctx) error {
	if middleware.GetRole(c) != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}

	id := c.Params("id")

	var body struct {
		PaidAmount int    `json:"paid_amount"`
		DueDate    string `json:"due_date"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var fee models.Fee
	result := database.DB.First(&fee, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "fee not found"})
	}

	if fee.Status == "paid" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fee already fully paid"})
	}

	fee.PaidAmount += body.PaidAmount
	fee.Remaining = fee.FinalAmount - fee.PaidAmount

	if fee.PaidAmount >= fee.FinalAmount {
		fee.Status = "paid"
		fee.Remaining = 0
	} else {
		fee.Status = "partial"
		fee.DueDate = body.DueDate
	}

	fee.PaidAt = time.Now().Format("2006-01-02 15:04:05")
	database.DB.Save(&fee)

	return c.JSON(fee)
}
