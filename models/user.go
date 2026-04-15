package models

import "gorm.io/gorm"

type User struct {
	BaseModel
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"-"`
	Role      string `json:"role"` // "admin", "teacher", "student"
	EpunjabID string `json:"epunjab_id" gorm:"uniqueIndex"`
	StudentID uint   `json:"student_id"`
}

func MigrateUser(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
