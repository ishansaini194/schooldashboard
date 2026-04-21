package models

import "gorm.io/gorm"

type Teacher struct {
	BaseModel
	SchoolID      uint   `json:"school_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	EmployeeID    string `json:"employee_id"`
	Subject       string `json:"subject"`
	Qualification string `json:"qualification"`
	IsActive      bool   `json:"is_active" gorm:"default:true"`
}

func MigrateTeacher(db *gorm.DB) error {
	return db.AutoMigrate(&Teacher{})
}
