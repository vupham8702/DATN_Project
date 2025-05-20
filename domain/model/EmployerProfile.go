package model

import (
	"database/sql"
)

type EmployerProfile struct {
	VModel
	UserID            uint           `json:"user_id" gorm:"uniqueIndex"`
	User              User           `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	CompanyName       string         `json:"company_name"`
	CompanySize       string         `json:"company_size"` // "1-10", "11-50", "51-200", "201-500", "501+"
	Industry          string         `json:"industry"`
	CompanyLogo       string         `json:"company_logo"`
	CompanyBanner     string         `json:"company_banner"`
	Website           string         `json:"website"`
	Founded           uint           `json:"founded"` // Year
	About             string         `json:"about"`
	Mission           string         `json:"mission"`
	PhoneNumber       string         `json:"phone_number"`
	Email             string         `json:"email"`
	Address           string         `json:"address"`
	City              string         `json:"city"`
	Country           string         `json:"country"`
	FacebookURL       string         `json:"facebook_url"`
	TwitterURL        string         `json:"twitter_url"`
	LinkedinURL       string         `json:"linkedin_url"`
	Benefits          sql.NullString `json:"benefits" gorm:"type:jsonb"`
	Culture           sql.NullString `json:"culture" gorm:"type:jsonb"`
	ProfileComplete   bool           `json:"profile_complete" gorm:"default:false"`
	TaxCode           string         `json:"tax_code"`
	BusinessLicense   string         `json:"business_license"`
	ContactPersonName string         `json:"contact_person_name"`
	ContactPersonRole string         `json:"contact_person_role"`
}
