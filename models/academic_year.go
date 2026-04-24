package models

type AcademicYear struct {
	BaseModel
	SchoolID  uint   `json:"school_id"`
	Name      string `json:"name"`       // "2025-26"
	StartDate string `json:"start_date"` // "2025-04-01"
	EndDate   string `json:"end_date"`   // "2026-03-31"
	IsCurrent bool   `json:"is_current" gorm:"default:false"`
}
