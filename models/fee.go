package models

import "gorm.io/gorm"

type Fee struct {
	gorm.Model
	StudentID      uint   `json:"student_id"`
	RollNo         string `json:"roll_no"`
	EpunjabID      string `json:"epunjab_id"`
	StudentName    string `json:"student_name"`
	Class          string `json:"class"`
	Month          string `json:"month"`
	Year           int    `json:"year"`
	FeeType        string `json:"fee_type"` // "tuition" / "transport"
	BaseAmount     int    `json:"base_amount"`
	Discount       int    `json:"discount"`
	DiscountReason string `json:"discount_reason"`
	FinalAmount    int    `json:"final_amount"` // base - discount
	PaidAmount     int    `json:"paid_amount"`
	Remaining      int    `json:"remaining"` // final - paid
	Status         string `json:"status"`    // "paid" / "partial" / "unpaid"
	DueDate        string `json:"due_date"`
	ReceiptNo      string `json:"receipt_no"`
	PaidAt         string `json:"paid_at"`
}

func MigrateFee(db *gorm.DB) error {
	return db.AutoMigrate(&Fee{})
}
