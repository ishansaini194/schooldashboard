package models

import "gorm.io/gorm"

type Homework struct {
	BaseModel
	ClassID   uint   `json:"class_id"`
	TeacherID uint   `json:"teacher_id"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
	DueDate   string `json:"due_date"`
}

func MigrateHomework(db *gorm.DB) error {
	return db.AutoMigrate(&Homework{})
}
