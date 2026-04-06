package models

import "gorm.io/gorm"

type Student struct {
	ID                    int    `json:"id"`
	ClassID               int    `json:"class_id"`
	Name                  string `json:"name"`
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
	CreatedAt             string `json:"created_at"`
}

func MigrateStudent(db *gorm.DB) error {
	return db.AutoMigrate(&Student{})
}
