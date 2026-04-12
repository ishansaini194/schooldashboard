package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

func MigrateUser(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
