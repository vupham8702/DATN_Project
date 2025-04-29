package model

import (
	"database/sql"
	"time"
)

type JobseekerProfile struct {
	VModel
	UserID          uint           `json:"user_id" gorm:"uniqueIndex"`
	User            User           `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	DateOfBirth     *time.Time     `json:"date_of_birth"`
	Gender          string         `json:"gender"` // "male", "female", "other"
	PhoneNumber     string         `json:"phone_number"`
	Address         string         `json:"address"`
	City            string         `json:"city"`
	Country         string         `json:"country"`
	ProfileTitle    string         `json:"profile_title"`
	About           string         `json:"about"`
	Skills          string         `json:"skills"`
	Education       sql.NullString `json:"education" gorm:"type:jsonb"`
	Experience      sql.NullString `json:"experience" gorm:"type:jsonb"`
	Certifications  sql.NullString `json:"certifications" gorm:"type:jsonb"`
	Languages       sql.NullString `json:"languages" gorm:"type:jsonb"`
	ResumeURL       string         `json:"resume_url"`
	ProfilePicture  string         `json:"profile_picture"`
	ProfileComplete bool           `json:"profile_complete" gorm:"default:false"`
	Availability    string         `json:"availability"` // "immediate", "2_weeks", "1_month"
	LinkedinProfile string         `json:"linkedin_profile"`
	GithubProfile   string         `json:"github_profile"`
	WebsiteURL      string         `json:"website_url"`
}
