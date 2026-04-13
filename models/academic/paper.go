package academic

import "gorm.io/gorm"

type Paper struct {
	gorm.Model
	Class      string `json:"class"`
	Section    string `json:"section"`
	Subject    string `json:"subject"`
	ExamType   string `json:"exam_type"` // "midterm" or "final"
	Year       int    `json:"year"`
	DriveLink  string `json:"drive_link"`
	UploadedBy string `json:"uploaded_by"`
}

func MigratePaper(db *gorm.DB) error {
	return db.AutoMigrate(&Paper{})
}
