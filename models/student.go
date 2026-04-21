package models

import "gorm.io/gorm"

type Student struct {
	BaseModel
	SchoolID              uint   `json:"school_id"`
	Name                  string `json:"name"`
	Phone                 string `json:"phone"`
	Gender                string `json:"gender"`
	DOB                   string `json:"dob"`
	AadharNo              string `json:"aadhar_no"`
	EpunjabId             string `json:"epunjab_id"`
	FatherName            string `json:"father_name"`
	FatherContact         string `json:"father_contact"`
	FatherAadhar          string `json:"father_aadhar"`
	MotherName            string `json:"mother_name"`
	MotherContact         string `json:"mother_contact"`
	Address               string `json:"address"`
	Caste                 string `json:"caste"`
	PreviousSchoolDetails string `json:"previous_school_details"`
	IsActive              bool   `json:"is_active" gorm:"default:true"`
}

func MigrateStudent(db *gorm.DB) error {
	return db.AutoMigrate(&Student{})
}
