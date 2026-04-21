package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	BaseModel
	FeeID       uint      `json:"fee_id"`
	CollectedBy uint      `json:"collected_by"` // user_id
	Amount      int       `json:"amount"`
	ReceiptNo   string    `json:"receipt_no" gorm:"uniqueIndex"`
	PaymentMode string    `json:"payment_mode" gorm:"default:'cash'"` // cash/online/cheque
	PaidAt      time.Time `json:"paid_at"`
	Notes       string    `json:"notes"`
}

func MigratePayment(db *gorm.DB) error {
	return db.AutoMigrate(&Payment{})
}
