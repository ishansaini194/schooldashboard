package models

import "gorm.io/gorm"

type Fee struct {
	gorm.Model
	StudentID int    `json:"student_id"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
	Amount    int    `json:"amount"`
	Status    string `json:"status"`
	PaidAt    string `json:"paid_at"`
}

func MigrateFee(db *gorm.DB) error {
	return db.AutoMigrate(&Fee{})
}
