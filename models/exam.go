package models

import "gorm.io/gorm"

type Exam struct {
	BaseModel
	ClassID        uint   `json:"class_id"`
	AcademicYearID uint   `json:"academic_year_id"`
	Name           string `json:"name"` // "mid-term", "final", "unit-1"
	Subject        string `json:"subject"`
	MaxMarks       int    `json:"max_marks"`
	Date           string `json:"date"`
}

func MigrateExam(db *gorm.DB) error {
	return db.AutoMigrate(&Exam{})
}
