package models

import "gorm.io/gorm"

type Result struct {
	BaseModel
	ExamID       uint `json:"exam_id"`
	EnrollmentID uint `json:"enrollment_id"`
	Marks        int  `json:"marks"`
	EnteredBy    uint `json:"entered_by"` // user_id
}

func MigrateResult(db *gorm.DB) error {
	return db.AutoMigrate(&Result{})
}
