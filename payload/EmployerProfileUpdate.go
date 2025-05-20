package payload

import (
	"database/sql"
)

// EmployerProfileUpdate là payload cho cập nhật hồ sơ nhà tuyển dụng
type EmployerProfileUpdate struct {
	CompanyName        string         `json:"company_name"`
	CompanySize        string         `json:"company_size"`
	Industry           string         `json:"industry"`
	CompanyDescription string         `json:"company_description"`
	CompanyLogo        string         `json:"company_logo"`
	CompanyWebsite     string         `json:"company_website"`
	CompanyAddress     string         `json:"company_address"`
	CompanyCity        string         `json:"company_city"`
	CompanyCountry     string         `json:"company_country"`
	ContactPerson      string         `json:"contact_person"`
	ContactPosition    string         `json:"contact_position"`
	ContactEmail       string         `json:"contact_email"`
	ContactPhone       string         `json:"contact_phone"`
	Benefits           sql.NullString `json:"benefits"`
	Culture            string         `json:"culture"`
	LinkedinProfile    string         `json:"linkedin_profile"`
	FacebookProfile    string         `json:"facebook_profile"`
	TwitterProfile     string         `json:"twitter_profile"`
	FoundedYear        int            `json:"founded_year"`
	TaxID              string         `json:"tax_id"`
	BusinessLicense    string         `json:"business_license"`
	BusinessLicenseURL string         `json:"business_license_url"`
}
