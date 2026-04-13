package academic

import "gorm.io/gorm"

type Result struct {
	gorm.Model
	StudentID uint   `json:"student_id"`
	Subject   string `json:"subject"`
	ExamType  string `json:"exam_type"` // "midterm" or "final"
	Marks     int    `json:"marks"`
	MaxMarks  int    `json:"max_marks"`
	Year      int    `json:"year"`
	Class     string `json:"class"`
	Section   string `json:"section"`
	EnteredBy string `json:"entered_by"`
}

func MigrateResult(db *gorm.DB) error {
	return db.AutoMigrate(&Result{})
}
