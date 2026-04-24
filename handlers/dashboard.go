package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/models"
)

// GET /api/dashboard/summary?month=April&year=2026
func GetDashboardSummary(c *fiber.Ctx) error {
	month := c.Query("month", time.Now().Format("January"))
	year := c.Query("year", fmt.Sprintf("%d", time.Now().Year()))
	monthInt := monthNums[month]
	ayID := currentAcademicYearID()
	sid := schoolID()

	// total students (active enrollments this year)
	var studentCount int64
	database.DB.Model(&models.Enrollment{}).
		Joins("JOIN classes ON classes.id = enrollments.class_id").
		Where("classes.school_id = ? AND enrollments.academic_year_id = ? AND enrollments.status = 'active'", sid, ayID).
		Count(&studentCount)

	// total classes
	var classCount int64
	database.DB.Model(&models.Class{}).Where("school_id = ?", sid).Count(&classCount)

	// fees collected this month
	var totalCollected int
	database.DB.Model(&models.Payment{}).
		Joins("JOIN fees ON fees.id = payments.fee_id").
		Joins("JOIN enrollments ON enrollments.id = fees.enrollment_id").
		Joins("JOIN classes ON classes.id = enrollments.class_id").
		Where("classes.school_id = ? AND fees.month = ? AND fees.year = ?", sid, monthInt, year).
		Select("COALESCE(SUM(payments.amount), 0)").Scan(&totalCollected)

	// expected total (sum of tuition fees for all enrolled students)
	type EnrollmentClass struct {
		TuitionFee int
	}
	var ecs []EnrollmentClass
	database.DB.Table("enrollments").
		Select("classes.tuition_fee").
		Joins("JOIN classes ON classes.id = enrollments.class_id").
		Where("classes.school_id = ? AND enrollments.academic_year_id = ? AND enrollments.status = 'active'", sid, ayID).
		Scan(&ecs)

	expectedTotal := 0
	for _, ec := range ecs {
		expectedTotal += ec.TuitionFee
	}

	// pending count (students who haven't fully paid this month)
	var paidEnrollmentIDs []uint
	database.DB.Model(&models.Fee{}).
		Joins("JOIN enrollments ON enrollments.id = fees.enrollment_id").
		Joins("JOIN classes ON classes.id = enrollments.class_id").
		Where("classes.school_id = ? AND fees.month = ? AND fees.year = ? AND fees.status = 'paid'",
			sid, monthInt, year).
		Pluck("fees.enrollment_id", &paidEnrollmentIDs)

	pendingCount := int(studentCount) - len(paidEnrollmentIDs)
	if pendingCount < 0 {
		pendingCount = 0
	}

	// overdue (partial fees from previous months)
	currentMonthInt := int(time.Now().Month())
	currentYear := time.Now().Year()
	var overdueCount int64
	database.DB.Model(&models.Fee{}).
		Joins("JOIN enrollments ON enrollments.id = fees.enrollment_id").
		Joins("JOIN classes ON classes.id = enrollments.class_id").
		Where("classes.school_id = ? AND fees.status = 'partial' AND (fees.year < ? OR (fees.year = ? AND fees.month < ?))",
			sid, currentYear, currentYear, currentMonthInt).
		Count(&overdueCount)

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
