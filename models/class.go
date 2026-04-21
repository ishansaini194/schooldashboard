package models

import "gorm.io/gorm"

type Class struct {
	BaseModel
	SchoolID       uint   `json:"school_id"`
	Number         int    `json:"number"`
	Section        string `json:"section"`
	ClassTeacherID *uint  `json:"class_teacher_id"`
	TuitionFee     int    `json:"tuition_fee"`
	TransportFee   int    `json:"transport_fee"`
}

func MigrateClass(db *gorm.DB) error {
	return db.AutoMigrate(&Class{})
}
