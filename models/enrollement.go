package models

import "gorm.io/gorm"

type Enrollment struct {
	BaseModel
	StudentID      uint   `json:"student_id"`
	ClassID        uint   `json:"class_id"`
	AcademicYearID uint   `json:"academic_year_id"`
	RollNo         int    `json:"roll_no"`
	Status         string `json:"status" gorm:"default:'active'"` // active/promoted/left
}

func MigrateEnrollment(db *gorm.DB) error {
	return db.AutoMigrate(&Enrollment{})
}
