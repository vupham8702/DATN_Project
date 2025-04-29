package request

import (
	"time"
)

// JobseekerProfileRequest là request body cho cập nhật hồ sơ ứng viên
type JobseekerProfileRequest struct {
	DateOfBirth  *time.Time `json:"date_of_birth"`
	Gender       string     `json:"gender"`
	PhoneNumber  string     `json:"phone_number"`
	Address      string     `json:"address"`
	City         string     `json:"city"`
	Country      string     `json:"country"`
	ProfileTitle string     `json:"profile_title"`
	About        string     `json:"about"`
	Skills       string     `json:"skills"`
	Education    []struct {
		Institution   string     `json:"institution"`
		Degree        string     `json:"degree"`
		Field         string     `json:"field"`
		StartDate     *time.Time `json:"start_date"`
		EndDate       *time.Time `json:"end_date"`
		CurrentlyHere bool       `json:"currently_here"`
		Description   string     `json:"description"`
	} `json:"education"`
	Experience []struct {
		Company       string     `json:"company"`
		Position      string     `json:"position"`
		StartDate     *time.Time `json:"start_date"`
		EndDate       *time.Time `json:"end_date"`
		CurrentlyHere bool       `json:"currently_here"`
		Description   string     `json:"description"`
	} `json:"experience"`
	Certifications []struct {
		Name         string     `json:"name"`
		Organisation string     `json:"organisation"`
		IssueDate    *time.Time `json:"issue_date"`
		ExpiryDate   *time.Time `json:"expiry_date"`
		CredentialID string     `json:"credential_id"`
		URL          string     `json:"url"`
	} `json:"certifications"`
	Languages []struct {
		Language string `json:"language"`
		Level    string `json:"level"` // "beginner", "intermediate", "advanced", "native"
	} `json:"languages"`
	Availability    string `json:"availability"`
	LinkedinProfile string `json:"linkedin_profile"`
	GithubProfile   string `json:"github_profile"`
	WebsiteURL      string `json:"website_url"`
}

// EmployerProfileRequest là request body cho cập nhật hồ sơ nhà tuyển dụng
type EmployerProfileRequest struct {
	CompanyName       string   `json:"company_name"`
	CompanySize       string   `json:"company_size"`
	Industry          string   `json:"industry"`
	Website           string   `json:"website"`
	Founded           uint     `json:"founded"`
	About             string   `json:"about"`
	Mission           string   `json:"mission"`
	PhoneNumber       string   `json:"phone_number"`
	Email             string   `json:"email"`
	Address           string   `json:"address"`
	City              string   `json:"city"`
	Country           string   `json:"country"`
	FacebookURL       string   `json:"facebook_url"`
	TwitterURL        string   `json:"twitter_url"`
	LinkedinURL       string   `json:"linkedin_url"`
	Benefits          []string `json:"benefits"`
	Culture           []string `json:"culture"`
	TaxCode           string   `json:"tax_code"`
	BusinessLicense   string   `json:"business_license"`
	ContactPersonName string   `json:"contact_person_name"`
	ContactPersonRole string   `json:"contact_person_role"`
}

// ProfilePhotoRequest là request body cho cập nhật ảnh đại diện
type ProfilePhotoRequest struct {
	PhotoType string `json:"photo_type"` // "profile", "resume", "company_logo", "company_banner"
	// Không cần field cho file vì sẽ được gửi qua form-data
}
