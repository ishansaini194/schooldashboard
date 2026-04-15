package academic

import (
	"github.com/ishansaini194/dashboard/models"
	"gorm.io/gorm"
)

type Paper struct {
	models.BaseModel
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
