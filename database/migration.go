package database

import (
	"log"

	"github.com/ishansaini194/dashboard/models"
	"github.com/ishansaini194/dashboard/models/academic"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	migrations := []func(*gorm.DB) error{
		models.MigrateClass,
		models.MigrateStudent,
		models.MigrateFee,
		models.MigrateUser,
		academic.MigrateHomework,
		academic.MigrateNotice,
		academic.MigrateResult,
		academic.MigratePaper,
	}

	for _, m := range migrations {
		if err := m(db); err != nil {
			log.Fatal("migration failed: ", err)
		}
	}

	// backfill: populate fee.section from student.section for older records
	backfillFeeSection(db)

	log.Println("all migrations done")
}

// one-time backfill: copy section from students into fees where section is empty
func backfillFeeSection(db *gorm.DB) {
	res := db.Exec(`
		UPDATE fees
		SET section = (SELECT section FROM students WHERE students.id = fees.student_id)
		WHERE (section IS NULL OR section = '')
		  AND student_id IN (SELECT id FROM students)
	`)
	if res.Error != nil {
		log.Println("fee section backfill skipped:", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		log.Printf("backfilled section on %d fee records\n", res.RowsAffected)
	}
}
