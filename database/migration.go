package database

import (
	"log"

	"github.com/ishansaini194/dashboard/models"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	migrations := []func(*gorm.DB) error{
		models.MigrateClass,
		models.MigrateStudent,
		models.MigrateFee,
	}

	for _, m := range migrations {
		if err := m(db); err != nil {
			log.Fatal("migration failed: ", err)
		}
	}

	log.Println("all migrations done")
}
