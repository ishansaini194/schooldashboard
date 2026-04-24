package database

import (
	"log"
	"os"

	"github.com/ishansaini194/dashboard/models"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.School{},
		&models.AcademicYear{},
		&models.Teacher{},
		&models.Student{},
		&models.User{},
		&models.Class{},
		&models.Enrollment{},
		&models.Fee{},
		&models.Payment{},
		&models.Homework{},
		&models.Notice{},
		&models.Exam{},
		&models.Result{},
		&models.Paper{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("All migrations done")

	seedDefaults(db)
}

// seedDefaults creates the default School and AcademicYear rows
// if they don't already exist. Safe to run on every startup.
func seedDefaults(db *gorm.DB) {
	schoolName := os.Getenv("SCHOOL_NAME")
	if schoolName == "" {
		schoolName = "KRB School"
	}
	schoolCode := os.Getenv("SCHOOL_CODE")
	if schoolCode == "" {
		schoolCode = "KRB"
	}

	var school models.School
	db.Where("code = ?", schoolCode).First(&school)
	if school.ID == 0 {
		school = models.School{Name: schoolName, Code: schoolCode}
		db.Create(&school)
		log.Printf("Created school: %s (id=%d)\n", schoolName, school.ID)
	}

	var ay models.AcademicYear
	db.Where("school_id = ? AND is_current = true", school.ID).First(&ay)
	if ay.ID == 0 {
		ayName := os.Getenv("ACADEMIC_YEAR")
		if ayName == "" {
			ayName = "2025-26"
		}
		ay = models.AcademicYear{
			SchoolID:  school.ID,
			Name:      ayName,
			StartDate: "2025-04-01",
			EndDate:   "2026-03-31",
			IsCurrent: true,
		}
		db.Create(&ay)
		log.Printf("Created academic year: %s (id=%d)\n", ayName, ay.ID)
	}
}
