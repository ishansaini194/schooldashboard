package models

import "gorm.io/gorm"

type AcademicYear struct {
	BaseModel
	SchoolID  uint   `json:"school_id"`
	Name      string `json:"name"`       // "2025-26"
	StartDate string `json:"start_date"` // "2025-04-01"
	EndDate   string `json:"end_date"`   // "2026-03-31"
	IsCurrent bool   `json:"is_current"`
}

func MigrateAcademicYear(db *gorm.DB) error {
	return db.AutoMigrate(&AcademicYear{})
}
