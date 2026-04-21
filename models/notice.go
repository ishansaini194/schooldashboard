package models

import "gorm.io/gorm"

type Notice struct {
	BaseModel
	SchoolID      uint   `json:"school_id"`
	PostedBy      uint   `json:"posted_by"` // user_id
	Title         string `json:"title"`
	Body          string `json:"body"`
	TargetType    string `json:"target_type"`     // all/class
	TargetClassID *uint  `json:"target_class_id"` // nullable
}

func MigrateNotice(db *gorm.DB) error {
	return db.AutoMigrate(&Notice{})
}
