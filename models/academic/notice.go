package academic

import "gorm.io/gorm"

type Notice struct {
	gorm.Model
	Title    string `json:"title"`
	Body     string `json:"body"`
	Target   string `json:"target"` // "all" or specific class like "8-A"
	PostedBy string `json:"posted_by"`
}

func MigrateNotice(db *gorm.DB) error {
	return db.AutoMigrate(&Notice{})
}
