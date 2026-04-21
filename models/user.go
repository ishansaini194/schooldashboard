package models

import "gorm.io/gorm"

type User struct {
	BaseModel
	SchoolID  uint   `json:"school_id"`
	Username  string `json:"username" gorm:"uniqueIndex"`
	Password  string `json:"-"`
	Role      string `json:"role"` // admin/teacher/student
	EpunjabID string `json:"epunjab_id"`
	StudentID *uint  `json:"student_id"` // nullable
	TeacherID *uint  `json:"teacher_id"` // nullable (new)
	LastLogin string `json:"last_login"`
}

func MigrateUser(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
