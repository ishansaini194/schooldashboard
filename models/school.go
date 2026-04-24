package models

type School struct {
	BaseModel
	Name    string `json:"name"`
	Code    string `json:"code" gorm:"uniqueIndex"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	LogoURL string `json:"logo_url"`
}

func MigrateSchool(db interface{ AutoMigrate(...interface{}) error }) error {
	return nil // handled centrally in migration.go
}
