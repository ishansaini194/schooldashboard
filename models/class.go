package models

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	Class          int    `json:"class"`
	Section        string `json:"section"`
	TeacherName    string `json:"teacher_name"`
	TeacherContact string `json:"teacher_contact"`
}

func MigrateClass(db *gorm.DB) error {
	return db.AutoMigrate(&Class{})
}
