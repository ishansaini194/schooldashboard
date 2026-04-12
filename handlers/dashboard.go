package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/models"
)

func GetDashboardSummary(c *fiber.Ctx) error {
	month := c.Query("month", time.Now().Format("January"))
	year := c.Query("year", fmt.Sprintf("%d", time.Now().Year()))

	// total students
	var studentCount int64
	database.DB.Model(&models.Student{}).Count(&studentCount)

	// total classes
	var classCount int64
	database.DB.Model(&models.Class{}).Count(&classCount)

	// fees this month
	var fees []models.Fee
	database.DB.Where("month = ? AND year = ?", month, year).Find(&fees)

	totalCollected := 0
	for _, f := range fees {
		totalCollected += f.PaidAmount
	}

	// students who paid fully this month
	paidStudentIDs := map[uint]bool{}
	for _, f := range fees {
		if f.Status == "paid" {
			paidStudentIDs[f.StudentID] = true
		}
	}

	pendingCount := int(studentCount) - len(paidStudentIDs)

	// overdue count
	today := time.Now().Format("2006-01-02")
	var overdueCount int64
	database.DB.Model(&models.Fee{}).
		Where("status = ? AND due_date != '' AND due_date < ?", "partial", today).
		Count(&overdueCount)

	// expected total — 2 queries only
	var students []models.Student
	database.DB.Find(&students)

	var classes []models.Class
	database.DB.Find(&classes)

	classMap := map[string]int{}
	for _, cls := range classes {
		classMap[fmt.Sprintf("%d", cls.Class)] = cls.TuitionFee
	}

	expectedTotal := 0
	for _, s := range students {
		expectedTotal += classMap[s.Class]
	}

	return c.JSON(fiber.Map{
		"total_students":  studentCount,
		"total_classes":   classCount,
		"total_collected": totalCollected,
		"expected_total":  expectedTotal,
		"pending_count":   pendingCount,
		"overdue_count":   overdueCount,
		"month":           month,
		"year":            year,
	})
}

// GET /api/fees/recent
func GetRecentPayments(c *fiber.Ctx) error {
	var fees []models.Fee
	database.DB.Where("paid_amount > 0").
		Order("created_at desc").
		Limit(5).
		Find(&fees)

	return c.JSON(fees)
}

// GET /api/fees/overdue
func GetOverdueFees(c *fiber.Ctx) error {
	today := time.Now().Format("2006-01-02")

	var fees []models.Fee
	database.DB.Where("status = ? AND due_date != '' AND due_date < ?", "partial", today).
		Order("due_date asc").
		Find(&fees)

	return c.JSON(fees)
}

// GET /api/fees/pending/all?month=April&year=2026
func GetAllPendingFees(c *fiber.Ctx) error {
	month := c.Query("month", time.Now().Format("January"))
	year := c.Query("year", fmt.Sprintf("%d", time.Now().Year()))

	// get all students
	var students []models.Student
	database.DB.Order("class asc, roll_no asc").Find(&students)

	// get all fees for this month/year
	var fees []models.Fee
	database.DB.Where("month = ? AND year = ?", month, year).Find(&fees)

	// map fees by student_id
	feeMap := map[uint]models.Fee{}
	for _, f := range fees {
		feeMap[f.StudentID] = f
	}

	type PendingStudent struct {
		StudentID   uint   `json:"student_id"`
		StudentName string `json:"student_name"`
		RollNo      string `json:"roll_no"`
		Class       string `json:"class"`
		Phone       string `json:"phone"`
		Status      string `json:"status"`
		PaidAmount  int    `json:"paid_amount"`
		Remaining   int    `json:"remaining"`
		DueDate     string `json:"due_date"`
		ReceiptNo   string `json:"receipt_no"`
	}

	var pending []PendingStudent
	for _, s := range students {
		f, hasFee := feeMap[s.ID]
		if !hasFee {
			// completely unpaid
			pending = append(pending, PendingStudent{
				StudentID:   s.ID,
				StudentName: s.Name,
				RollNo:      s.RollNo,
				Class:       s.Class,
				Phone:       s.Phone,
				Status:      "unpaid",
				PaidAmount:  0,
				Remaining:   0,
			})
		} else if f.Status == "partial" {
			// partial payment
			pending = append(pending, PendingStudent{
				StudentID:   s.ID,
				StudentName: s.Name,
				RollNo:      s.RollNo,
				Class:       s.Class,
				Phone:       s.Phone,
				Status:      "partial",
				PaidAmount:  f.PaidAmount,
				Remaining:   f.Remaining,
				DueDate:     f.DueDate,
				ReceiptNo:   f.ReceiptNo,
			})
		}
		// skip fully paid students
	}

	return c.JSON(fiber.Map{
		"month":   month,
		"year":    year,
		"count":   len(pending),
		"pending": pending,
	})
}
