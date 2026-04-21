package models

import "gorm.io/gorm"

type Paper struct {
	BaseModel
	ClassID   uint   `json:"class_id"`
	TeacherID uint   `json:"teacher_id"`
	Subject   string `json:"subject"`
	ExamType  string `json:"exam_type"` // midterm/final
	Year      int    `json:"year"`
	DriveLink string `json:"drive_link"`
}

func MigratePaper(db *gorm.DB) error {
	return db.AutoMigrate(&Paper{})
}
