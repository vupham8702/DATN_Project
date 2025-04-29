package payload

import (
	"database/sql"
	"time"
)

// JobseekerProfileUpdate là payload cho cập nhật hồ sơ ứng viên
type JobseekerProfileUpdate struct {
	DateOfBirth     *time.Time     `json:"date_of_birth"`
	Gender          string         `json:"gender"`
	PhoneNumber     string         `json:"phone_number"`
	Address         string         `json:"address"`
	City            string         `json:"city"`
	Country         string         `json:"country"`
	ProfileTitle    string         `json:"profile_title"`
	About           string         `json:"about"`
	Skills          string         `json:"skills"`
	Education       sql.NullString `json:"education"`
	Experience      sql.NullString `json:"experience"`
	Certifications  sql.NullString `json:"certifications"`
	Languages       sql.NullString `json:"languages"`
	ResumeURL       string         `json:"resume_url"`
	ProfilePicture  string         `json:"profile_picture"`
	Availability    string         `json:"availability"`
	LinkedinProfile string         `json:"linkedin_profile"`
	GithubProfile   string         `json:"github_profile"`
	WebsiteURL      string         `json:"website_url"`
	ExpectedSalary  string         `json:"expected_salary"`
	JobPreferences  sql.NullString `json:"job_preferences"`
	Interests       string         `json:"interests"`
}
