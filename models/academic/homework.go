package academic

import (
	"time"

	"gorm.io/gorm"
)

type Homework struct {
	gorm.Model
	Class     string    `json:"class"`
	Section   string    `json:"section"`
	Subject   string    `json:"subject"`
	Content   string    `json:"content"`
	CreatedBy string    `json:"created_by"`
	Date      time.Time `json:"date"`
}

func MigrateHomework(db *gorm.DB) error {
	return db.AutoMigrate(&Homework{})
}
