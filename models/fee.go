package models

import "gorm.io/gorm"

type Fee struct {
	BaseModel
	EnrollmentID   uint   `json:"enrollment_id"`
	FeeType        string `json:"fee_type"`   // tuition / transport
	Month          int    `json:"month"`      // 1-12
	Year           int    `json:"year"`
	Amount         int    `json:"amount"`     // base fee
	Discount       int    `json:"discount"`
	DiscountReason string `json:"discount_reason"`
	NetAmount      int    `json:"net_amount"` // amount - discount
	Status         string `json:"status" gorm:"default:'unpaid'"` // unpaid / partial / paid
}

func MigrateFee(db *gorm.DB) error {
	return db.AutoMigrate(&Fee{})
}
