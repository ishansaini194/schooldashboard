package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/middleware"
	"github.com/ishansaini194/dashboard/models"
)

// ── Receipt number ────────────────────────────────────────────

func generateReceiptNo() string {
	var count int64
	database.DB.Model(&models.Payment{}).Count(&count)
	return fmt.Sprintf("RCPT-%d-%03d", time.Now().Year(), count+1)
}

// ── Shared response types (same shape as before) ──────────────

type FeeWithPayment struct {
	ID             uint    `json:"id"`
	StudentID      uint    `json:"student_id"`
	EnrollmentID   uint    `json:"enrollment_id"`
	StudentName    string  `json:"student_name"`
	RollNo         string  `json:"roll_no"`
	EpunjabID      string  `json:"epunjab_id"`
	Class          string  `json:"class"`
	Section        string  `json:"section"`
	Month          string  `json:"month"` // returned as name e.g. "April"
	Year           int     `json:"year"`
	FeeType        string  `json:"fee_type"`
	BaseAmount     int     `json:"base_amount"`
	Discount       int     `json:"discount"`
	DiscountReason string  `json:"discount_reason"`
	FinalAmount    int     `json:"final_amount"`
	PaidAmount     int     `json:"paid_amount"`
	Remaining      int     `json:"remaining"`
	Status         string  `json:"status"`
	ReceiptNo      string  `json:"receipt_no"`
	PaidAt         string  `json:"paid_at"`
}

var monthNames = [13]string{"", "January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December"}

var monthNums = map[string]int{
	"January": 1, "February": 2, "March": 3, "April": 4,
	"May": 5, "June": 6, "July": 7, "August": 8,
	"September": 9, "October": 10, "November": 11, "December": 12,
}

// enrichFee loads related data and builds FeeWithPayment
func enrichFee(fee models.Fee) FeeWithPayment {
	var e models.Enrollment
	database.DB.First(&e, fee.EnrollmentID)
	var s models.Student
	database.DB.First(&s, e.StudentID)
	var cls models.Class
	database.DB.First(&cls, e.ClassID)

	// latest payment for this fee
	var payment models.Payment
	database.DB.Where("fee_id = ?", fee.ID).Order("paid_at desc").First(&payment)

	// sum all payments
	var totalPaid int
	database.DB.Model(&models.Payment{}).
		Where("fee_id = ?", fee.ID).
		Select("COALESCE(SUM(amount),0)").Scan(&totalPaid)

	remaining := fee.NetAmount - totalPaid
	if remaining < 0 {
		remaining = 0
	}

	monthName := ""
	if fee.Month >= 1 && fee.Month <= 12 {
		monthName = monthNames[fee.Month]
	}

	paidAt := ""
	if !payment.PaidAt.IsZero() {
		paidAt = payment.PaidAt.Format("2006-01-02 15:04:05")
	}

	return FeeWithPayment{
		ID: fee.ID, EnrollmentID: fee.EnrollmentID,
		StudentID: s.ID, StudentName: s.Name,
		RollNo: fmt.Sprintf("%d", e.RollNo), EpunjabID: s.EpunjabId,
		Class: fmt.Sprintf("%d", cls.Number), Section: cls.Section,
		Month: monthName, Year: fee.Year,
		FeeType: fee.FeeType,
		BaseAmount: fee.Amount, Discount: fee.Discount,
		DiscountReason: fee.DiscountReason, FinalAmount: fee.NetAmount,
		PaidAmount: totalPaid, Remaining: remaining,
		Status: fee.Status, ReceiptNo: payment.ReceiptNo, PaidAt: paidAt,
	}
}

// POST /api/fees/pay
func PayFee(c *fiber.Ctx) error {
	if middleware.GetRole(c) != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}

	var body struct {
		StudentID      uint   `json:"student_id"`
		EnrollmentID   uint   `json:"enrollment_id"`
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
	if body.PaidAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "paid amount must be > 0"})
	}

	// resolve enrollment_id — try student_id + current year if not provided
	enrollmentID := body.EnrollmentID
	if enrollmentID == 0 && body.StudentID != 0 {
		var e models.Enrollment
		database.DB.Where("student_id = ? AND academic_year_id = ?", body.StudentID, currentAcademicYearID()).First(&e)
		enrollmentID = e.ID
	}
	if enrollmentID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "enrollment not found"})
	}

	monthInt := monthNums[body.Month]
	if monthInt == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid month"})
	}

	netAmount := body.BaseAmount - body.Discount

	// find or create fee record
	var fee models.Fee
	database.DB.Where("enrollment_id = ? AND fee_type = ? AND month = ? AND year = ?",
		enrollmentID, body.FeeType, monthInt, body.Year).First(&fee)

	if fee.ID == 0 {
		fee = models.Fee{
			EnrollmentID:   enrollmentID,
			FeeType:        body.FeeType,
			Month:          monthInt,
			Year:           body.Year,
			Amount:         body.BaseAmount,
			Discount:       body.Discount,
			DiscountReason: body.DiscountReason,
			NetAmount:      netAmount,
			Status:         "unpaid",
		}
		database.DB.Create(&fee)
	}

	// sum existing payments
	var alreadyPaid int
	database.DB.Model(&models.Payment{}).Where("fee_id = ?", fee.ID).
		Select("COALESCE(SUM(amount),0)").Scan(&alreadyPaid)

	// create payment
	payment := models.Payment{
		FeeID:       fee.ID,
		CollectedBy: middleware.GetUserID(c),
		Amount:      body.PaidAmount,
		ReceiptNo:   generateReceiptNo(),
		PaymentMode: "cash",
		PaidAt:      time.Now(),
	}
	database.DB.Create(&payment)

	// update fee status
	totalPaid := alreadyPaid + body.PaidAmount
	if totalPaid >= netAmount {
		fee.Status = "paid"
	} else {
		fee.Status = "partial"
	}
	database.DB.Save(&fee)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"receipt_no": payment.ReceiptNo,
		"status":     fee.Status,
		"fee_id":     fee.ID,
	})
}

// PUT /api/fees/:id/complete
func CompleteFee(c *fiber.Ctx) error {
	if middleware.GetRole(c) != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}

	var fee models.Fee
	database.DB.First(&fee, c.Params("id"))
	if fee.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "fee not found"})
	}
	if fee.Status == "paid" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "already fully paid"})
	}

	var body struct {
		PaidAmount int `json:"paid_amount"`
	}
	c.BodyParser(&body)

	payment := models.Payment{
		FeeID:       fee.ID,
		CollectedBy: middleware.GetUserID(c),
		Amount:      body.PaidAmount,
		ReceiptNo:   generateReceiptNo(),
		PaymentMode: "cash",
		PaidAt:      time.Now(),
	}
	database.DB.Create(&payment)

	var totalPaid int
	database.DB.Model(&models.Payment{}).Where("fee_id = ?", fee.ID).
		Select("COALESCE(SUM(amount),0)").Scan(&totalPaid)

	if totalPaid >= fee.NetAmount {
		fee.Status = "paid"
	} else {
		fee.Status = "partial"
	}
	database.DB.Save(&fee)

	return c.JSON(fiber.Map{"receipt_no": payment.ReceiptNo, "status": fee.Status})
}

// GET /api/fees/class/:class/month/:month/year/:year?section=A
func GetClassFeeStatus(c *fiber.Ctx) error {
	cls := findClass(c.Params("class"), c.Query("section"))
	if cls.ID == 0 {
		return c.JSON([]fiber.Map{})
	}

	monthInt := monthNums[c.Params("month")]
	year := c.Params("year")
	ayID := currentAcademicYearID()

	var enrollments []models.Enrollment
	database.DB.Where("class_id = ? AND academic_year_id = ? AND status = 'active'", cls.ID, ayID).
		Order("roll_no asc").Find(&enrollments)

	type StudentFeeStatus struct {
		StudentID   uint           `json:"student_id"`
		StudentName string         `json:"student_name"`
		RollNo      string         `json:"roll_no"`
		HasPaid     bool           `json:"has_paid"`
		TotalPaid   int            `json:"total_paid"`
		Fees        []FeeWithPayment `json:"fees"`
	}

	result := make([]StudentFeeStatus, 0, len(enrollments))
	for _, e := range enrollments {
		var s models.Student
		database.DB.First(&s, e.StudentID)

		var fees []models.Fee
		database.DB.Where("enrollment_id = ? AND month = ? AND year = ?", e.ID, monthInt, year).Find(&fees)

		var totalPaid int
		hasPaid := false
		feeResponses := make([]FeeWithPayment, 0)

		for _, f := range fees {
			var paid int
			database.DB.Model(&models.Payment{}).Where("fee_id = ?", f.ID).
				Select("COALESCE(SUM(amount),0)").Scan(&paid)
			totalPaid += paid
			if f.Status == "paid" {
				hasPaid = true
			}
			feeResponses = append(feeResponses, enrichFee(f))
		}

		result = append(result, StudentFeeStatus{
			StudentID:   s.ID,
			StudentName: s.Name,
			RollNo:      fmt.Sprintf("%d", e.RollNo),
			HasPaid:     hasPaid,
			TotalPaid:   totalPaid,
			Fees:        feeResponses,
		})
	}
	return c.JSON(result)
}

// GET /api/fees/pending/all?month=April&year=2026
func GetAllPendingFees(c *fiber.Ctx) error {
	month := c.Query("month")
	year := c.Query("year")
	monthInt := monthNums[month]
	ayID := currentAcademicYearID()

	var enrollments []models.Enrollment
	database.DB.Where("academic_year_id = ? AND status = 'active'", ayID).Find(&enrollments)

	type PendingStudent struct {
		StudentID   uint   `json:"student_id"`
		StudentName string `json:"student_name"`
		RollNo      string `json:"roll_no"`
		Class       string `json:"class"`
		Phone       string `json:"phone"`
		Status      string `json:"status"`
		PaidAmount  int    `json:"paid_amount"`
		Remaining   int    `json:"remaining"`
	}

	var pending []PendingStudent
	for _, e := range enrollments {
		var s models.Student
		database.DB.First(&s, e.StudentID)
		var cls models.Class
		database.DB.First(&cls, e.ClassID)

		var fee models.Fee
		database.DB.Where("enrollment_id = ? AND month = ? AND year = ?", e.ID, monthInt, year).First(&fee)

		if fee.ID == 0 {
			pending = append(pending, PendingStudent{
				StudentID: s.ID, StudentName: s.Name,
				RollNo: fmt.Sprintf("%d", e.RollNo),
				Class: fmt.Sprintf("%d", cls.Number), Phone: s.Phone,
				Status: "unpaid", PaidAmount: 0, Remaining: 0,
			})
		} else if fee.Status == "partial" {
			var totalPaid int
			database.DB.Model(&models.Payment{}).Where("fee_id = ?", fee.ID).
				Select("COALESCE(SUM(amount),0)").Scan(&totalPaid)
			pending = append(pending, PendingStudent{
				StudentID: s.ID, StudentName: s.Name,
				RollNo: fmt.Sprintf("%d", e.RollNo),
				Class: fmt.Sprintf("%d", cls.Number), Phone: s.Phone,
				Status: "partial", PaidAmount: totalPaid,
				Remaining: fee.NetAmount - totalPaid,
			})
		}
	}

	return c.JSON(fiber.Map{"month": month, "year": year, "count": len(pending), "pending": pending})
}

// GET /api/fees/recent
func GetRecentPayments(c *fiber.Ctx) error {
	var payments []models.Payment
	database.DB.Order("paid_at desc").Limit(10).Find(&payments)

	result := make([]fiber.Map, 0, len(payments))
	for _, p := range payments {
		var fee models.Fee
		database.DB.First(&fee, p.FeeID)
		fw := enrichFee(fee)
		result = append(result, fiber.Map{
			"student_name": fw.StudentName,
			"class":        fw.Class,
			"month":        fw.Month,
			"fee_type":     fw.FeeType,
			"paid_amount":  p.Amount,
			"receipt_no":   p.ReceiptNo,
			"paid_at":      p.PaidAt,
		})
	}
	return c.JSON(result)
}

// GET /api/fees/overdue
func GetOverdueFees(c *fiber.Ctx) error {
	var fees []models.Fee
	database.DB.Where("status = 'partial'").Order("created_at asc").Find(&fees)
	result := make([]FeeWithPayment, 0)
	for _, f := range fees {
		result = append(result, enrichFee(f))
	}
	return c.JSON(result)
}

// GET /api/fees/student/:student_id
func GetStudentFees(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	var e models.Enrollment
	database.DB.Where("student_id = ? AND academic_year_id = ?", studentID, currentAcademicYearID()).First(&e)

	var fees []models.Fee
	database.DB.Where("enrollment_id = ?", e.ID).Order("year desc, month desc").Find(&fees)

	result := make([]FeeWithPayment, 0)
	for _, f := range fees {
		result = append(result, enrichFee(f))
	}
	return c.JSON(result)
}

// GET /api/fees/student/:student_id/yearly?year=2026
func GetStudentYearlySummary(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	yearStr := c.Query("year", fmt.Sprintf("%d", time.Now().Year()))

	var e models.Enrollment
	database.DB.Where("student_id = ? AND academic_year_id = ?", studentID, currentAcademicYearID()).First(&e)

	var fees []models.Fee
	database.DB.Where("enrollment_id = ? AND year = ?", e.ID, yearStr).Find(&fees)

	feeByMonth := map[int]models.Fee{}
	for _, f := range fees {
		feeByMonth[f.Month] = f
	}

	currentMonth := int(time.Now().Month())
	paidCount := 0

	type MonthSummary struct {
		Month      string `json:"month"`
		HasRecord  bool   `json:"has_record"`
		Status     string `json:"status"`
		PaidAmount int    `json:"paid_amount"`
		Remaining  int    `json:"remaining"`
		FeeType    string `json:"fee_type"`
		ReceiptNo  string `json:"receipt_no"`
	}

	var summary []MonthSummary
	for i := 1; i <= currentMonth; i++ {
		if f, exists := feeByMonth[i]; exists {
			var totalPaid int
			database.DB.Model(&models.Payment{}).Where("fee_id = ?", f.ID).
				Select("COALESCE(SUM(amount),0)").Scan(&totalPaid)
			var lastPayment models.Payment
			database.DB.Where("fee_id = ?", f.ID).Order("paid_at desc").First(&lastPayment)
			if f.Status == "paid" {
				paidCount++
			}
			summary = append(summary, MonthSummary{
				Month: monthNames[i], HasRecord: true,
				Status: f.Status, PaidAmount: totalPaid,
				Remaining: f.NetAmount - totalPaid,
				FeeType: f.FeeType, ReceiptNo: lastPayment.ReceiptNo,
			})
		} else {
			summary = append(summary, MonthSummary{Month: monthNames[i], Status: "unpaid"})
		}
	}

	return c.JSON(fiber.Map{
		"student_id":   studentID,
		"year":         yearStr,
		"months":       summary,
		"paid_count":   paidCount,
		"total_months": len(summary),
	})
}

// GET /api/fees/receipt/:receipt_no
func GetReceipt(c *fiber.Ctx) error {
	var payment models.Payment
	database.DB.Where("receipt_no = ?", c.Params("receipt_no")).First(&payment)
	if payment.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "receipt not found"})
	}
	var fee models.Fee
	database.DB.First(&fee, payment.FeeID)
	fw := enrichFee(fee)
	return c.JSON(fw)
}

// GET /api/fees/pending/:class/:month/:year
func GetPendingFees(c *fiber.Ctx) error {
	cls := findClass(c.Params("class"), c.Query("section"))
	monthInt := monthNums[c.Params("month")]
	year := c.Params("year")
	ayID := currentAcademicYearID()

	var enrollments []models.Enrollment
	database.DB.Where("class_id = ? AND academic_year_id = ? AND status = 'active'", cls.ID, ayID).Find(&enrollments)

	type Pending struct {
		StudentID   uint   `json:"student_id"`
		StudentName string `json:"student_name"`
		RollNo      string `json:"roll_no"`
		Phone       string `json:"phone"`
	}

	var pending []Pending
	for _, e := range enrollments {
		var fee models.Fee
		database.DB.Where("enrollment_id = ? AND month = ? AND year = ? AND status != 'paid'",
			e.ID, monthInt, year).First(&fee)
		if fee.ID != 0 || true {
			var s models.Student
			database.DB.First(&s, e.StudentID)
			var existingFee models.Fee
			database.DB.Where("enrollment_id = ? AND month = ? AND year = ?", e.ID, monthInt, year).First(&existingFee)
			if existingFee.ID == 0 || existingFee.Status != "paid" {
				pending = append(pending, Pending{
					StudentID: s.ID, StudentName: s.Name,
					RollNo: fmt.Sprintf("%d", e.RollNo), Phone: s.Phone,
				})
			}
		}
	}
	return c.JSON(pending)
}
