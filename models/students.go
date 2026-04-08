package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Name                  string `json:"name"`
	Class                 string `json:"class"`
	Section               string `json:"section"`
	Phone                 string `json:"phone"`
	RollNo                string `json:"roll_no"`
	AadharNo              string `json:"aadhar_no"`
	EpunjabId             string `json:"epunjab_id"`
	FatherName            string `json:"father_name"`
	FatherContact         string `json:"father_contact"`
	FatherAadhar          string `json:"father_aadhar"`
	MotherName            string `json:"mother_name"`
	MotherContact         string `json:"mother_contact"`
	Address               string `json:"address"`
	Caste                 string `json:"caste"`
	Gender                string `json:"gender"`
	PreviousSchoolDetails string `json:"previous_school_details"`
	DOB                   string `json:"dob"`
}

func MigrateStudent(db *gorm.DB) error {
	return db.AutoMigrate(&Student{})
}
