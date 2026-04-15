package validators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reName       = regexp.MustCompile(`^[A-Za-z][A-Za-z\s.'-]{1,49}$`)
	reSection    = regexp.MustCompile(`^[A-Za-z]$`)
	reRollNo     = regexp.MustCompile(`^\d{1,4}$`)
	reEpunjab    = regexp.MustCompile(`^[A-Za-z0-9]{4,20}$`)
	rePhone      = regexp.MustCompile(`^\d{10}$`)
	reAadhar     = regexp.MustCompile(`^\d{12}$`)
	reCaste      = regexp.MustCompile(`^[A-Za-z\s]{2,30}$`)
)

// ValidateStudent returns a map of field -> error message.
// Empty map means valid.
func ValidateStudent(s StudentInput) map[string]string {
	errs := map[string]string{}

	// required: name
	name := strings.TrimSpace(s.Name)
	if name == "" {
		errs["name"] = "Required"
	} else if !reName.MatchString(name) {
		errs["name"] = "Only letters, spaces, dots allowed (2-50 chars)"
	}

	// required: class (1-12)
	cls := strings.TrimSpace(s.Class)
	if cls == "" {
		errs["class"] = "Required"
	} else {
		n, err := strconv.Atoi(cls)
		if err != nil || n < 1 || n > 12 {
			errs["class"] = "Must be a number between 1 and 12"
		}
	}

	// required: section (single letter)
	section := strings.TrimSpace(s.Section)
	if section == "" {
		errs["section"] = "Required"
	} else if !reSection.MatchString(section) {
		errs["section"] = "Must be a single letter (A-Z)"
	}

	// required: roll_no (1-4 digits)
	roll := strings.TrimSpace(s.RollNo)
	if roll == "" {
		errs["roll_no"] = "Required"
	} else if !reRollNo.MatchString(roll) {
		errs["roll_no"] = "Must be 1-4 digits"
	}

	// required: epunjab_id (alphanumeric 4-20)
	epunjab := strings.TrimSpace(s.EpunjabId)
	if epunjab == "" {
		errs["epunjab_id"] = "Required"
	} else if !reEpunjab.MatchString(epunjab) {
		errs["epunjab_id"] = "Must be 4-20 letters/digits"
	}

	// required: phone (10 digits)
	phone := strings.TrimSpace(s.Phone)
	if phone == "" {
		errs["phone"] = "Required"
	} else if !rePhone.MatchString(phone) {
		errs["phone"] = "Must be exactly 10 digits"
	}

	// optional fields below
	if v := strings.TrimSpace(s.AadharNo); v != "" && !reAadhar.MatchString(v) {
		errs["aadhar_no"] = "Must be exactly 12 digits"
	}
	if v := strings.TrimSpace(s.FatherName); v != "" && !reName.MatchString(v) {
		errs["father_name"] = "Only letters, spaces, dots allowed"
	}
	if v := strings.TrimSpace(s.FatherContact); v != "" && !rePhone.MatchString(v) {
		errs["father_contact"] = "Must be exactly 10 digits"
	}
	if v := strings.TrimSpace(s.FatherAadhar); v != "" && !reAadhar.MatchString(v) {
		errs["father_aadhar"] = "Must be exactly 12 digits"
	}
	if v := strings.TrimSpace(s.MotherName); v != "" && !reName.MatchString(v) {
		errs["mother_name"] = "Only letters, spaces, dots allowed"
	}
	if v := strings.TrimSpace(s.MotherContact); v != "" && !rePhone.MatchString(v) {
		errs["mother_contact"] = "Must be exactly 10 digits"
	}
	if v := strings.TrimSpace(s.Caste); v != "" && !reCaste.MatchString(v) {
		errs["caste"] = "Only letters and spaces (2-30 chars)"
	}
	if v := strings.TrimSpace(s.Gender); v != "" && v != "male" && v != "female" {
		errs["gender"] = "Must be 'male' or 'female'"
	}
	if v := strings.TrimSpace(s.DOB); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			errs["dob"] = "Invalid date format"
		} else {
			years := time.Since(t).Hours() / 24 / 365.25
			if years < 3 || years > 25 {
				errs["dob"] = "Age must be between 3 and 25 years"
			}
		}
	}
	if v := strings.TrimSpace(s.Address); v != "" && len(v) < 5 {
		errs["address"] = "Min 5 characters"
	}

	return errs
}

// StudentInput is a validation-only view of student fields.
// Defined here to keep validators independent from the models package.
type StudentInput struct {
	Name          string
	Class         string
	Section       string
	RollNo        string
	EpunjabId     string
	Phone         string
	AadharNo      string
	FatherName    string
	FatherContact string
	FatherAadhar  string
	MotherName    string
	MotherContact string
	Caste         string
	Gender        string
	DOB           string
	Address       string
}

// FormatErrors turns the error map into a single human-readable string.
func FormatErrors(errs map[string]string) string {
	parts := make([]string, 0, len(errs))
	for field, msg := range errs {
		parts = append(parts, fmt.Sprintf("%s: %s", field, msg))
	}
	return strings.Join(parts, "; ")
}
