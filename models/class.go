package models

import "gorm.io/gorm"

type Class struct {
	ID             int    `json:"id"`
	ClassNo        int    `json:"class_no"`
	Section        string `json:"section"`
	TeacherName    string `json:"teacher_name"`
	TeacherContact string `json:"teacher_contact"`
}

func MigrateClass(db *gorm.DB) error {
	return db.AutoMigrate(&Class{})
}
