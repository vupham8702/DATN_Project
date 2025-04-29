package service

import (
	"database/sql"
	"datn_backend/config"
	"datn_backend/domain/model"
	repo "datn_backend/domain/repository"
	"datn_backend/message"
	"datn_backend/middleware"
	"datn_backend/payload/request"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetProfile lấy thông tin hồ sơ của người dùng dựa trên loại người dùng
func GetProfile(userID *uint) (interface{}, interface{}) {
	// Xác định loại người dùng
	userType, err := repo.GetUserType(*userID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get user type: %v", err))
		return nil, message.UserNotFound
	}

	// Lấy thông tin hồ sơ dựa vào loại người dùng
	switch userType {
	case config.USER_TYPE_JOBSEEKER:
		profile, err := repo.GetJobseekerProfileByUserID(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Nếu chưa có hồ sơ, trả về hồ sơ trống
				return &model.JobseekerProfile{
					UserID: *userID,
				}, nil
			}
			middleware.Log(fmt.Errorf("Failed to get jobseeker profile: %v", err))
			return nil, message.InternalServerError
		}
		return profile, nil
	case config.USER_TYPE_EMPLOYER:
		profile, err := repo.GetEmployerProfileByUserID(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Nếu chưa có hồ sơ, trả về hồ sơ trống
				return &model.EmployerProfile{
					UserID: *userID,
				}, nil
			}
			middleware.Log(fmt.Errorf("Failed to get employer profile: %v", err))
			return nil, message.InternalServerError
		}
		return profile, nil
	default:
		return nil, message.Message{Message: "Invalid user type", Code: 400}
	}
}

// UpdateJobseekerProfile cập nhật thông tin hồ sơ của ứng viên
func UpdateJobseekerProfile(userID uint, profileData *request.JobseekerProfileRequest) (interface{}, interface{}) {
	// Kiểm tra loại người dùng
	userType, err := repo.GetUserType(userID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get user type: %v", err))
		return nil, message.UserNotFound
	}

	if userType != config.USER_TYPE_JOBSEEKER {
		return nil, message.Message{Message: "User is not a jobseeker", Code: 400}
	}

	// Lấy hồ sơ hiện tại nếu có
	var profile model.JobseekerProfile
	uid := userID
	existingProfile, err := repo.GetJobseekerProfileByUserID(&uid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		middleware.Log(fmt.Errorf("Failed to get jobseeker profile: %v", err))
		return nil, message.InternalServerError
	}

	if existingProfile != nil {
		profile = *existingProfile
	} else {
		profile = model.JobseekerProfile{UserID: userID}
	}

	// Cập nhật thông tin từ request
	profile.DateOfBirth = profileData.DateOfBirth
	profile.Gender = profileData.Gender
	profile.PhoneNumber = profileData.PhoneNumber
	profile.Address = profileData.Address
	profile.City = profileData.City
	profile.Country = profileData.Country
	profile.ProfileTitle = profileData.ProfileTitle
	profile.About = profileData.About
	profile.Skills = profileData.Skills
	profile.Availability = profileData.Availability
	profile.LinkedinProfile = profileData.LinkedinProfile
	profile.GithubProfile = profileData.GithubProfile
	profile.WebsiteURL = profileData.WebsiteURL

	// Chuyển đổi các mảng thành JSON
	if len(profileData.Education) > 0 {
		educationJSON, err := json.Marshal(profileData.Education)
		if err != nil {
			middleware.Log(fmt.Errorf("Failed to marshal education: %v", err))
			return nil, message.InternalServerError
		}
		profile.Education = sql.NullString{String: string(educationJSON), Valid: true}
	}

	if len(profileData.Experience) > 0 {
		experienceJSON, err := json.Marshal(profileData.Experience)
		if err != nil {
			middleware.Log(fmt.Errorf("Failed to marshal experience: %v", err))
			return nil, message.InternalServerError
		}
		profile.Experience = sql.NullString{String: string(experienceJSON), Valid: true}
	}

	if len(profileData.Certifications) > 0 {
		certificationsJSON, err := json.Marshal(profileData.Certifications)
		if err != nil {
			middleware.Log(fmt.Errorf("Failed to marshal certifications: %v", err))
			return nil, message.InternalServerError
		}
		profile.Certifications = sql.NullString{String: string(certificationsJSON), Valid: true}
	}

	if len(profileData.Languages) > 0 {
		languagesJSON, err := json.Marshal(profileData.Languages)
		if err != nil {
			middleware.Log(fmt.Errorf("Failed to marshal languages: %v", err))
			return nil, message.InternalServerError
		}
		profile.Languages = sql.NullString{String: string(languagesJSON), Valid: true}
	}

	// Kiểm tra hồ sơ có đầy đủ thông tin chưa
	profile.ProfileComplete = isJobseekerProfileComplete(&profile)

	// Lưu vào database
	err = repo.UpsertJobseekerProfile(&profile)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to update jobseeker profile: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return profile, nil
}

// UpdateEmployerProfile cập nhật thông tin hồ sơ của nhà tuyển dụng
func UpdateEmployerProfile(userID uint, profileData *request.EmployerProfileRequest) (interface{}, interface{}) {
	// Kiểm tra loại người dùng
	userType, err := repo.GetUserType(userID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get user type: %v", err))
		return nil, message.UserNotFound
	}

	if userType != config.USER_TYPE_EMPLOYER {
		return nil, message.Message{Message: "User is not an employer", Code: 400}
	}

	// Lấy hồ sơ hiện tại nếu có
	var profile model.EmployerProfile
	uid := userID
	existingProfile, err := repo.GetEmployerProfileByUserID(&uid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		middleware.Log(fmt.Errorf("Failed to get employer profile: %v", err))
		return nil, message.InternalServerError
	}

	if existingProfile != nil {
		profile = *existingProfile
	} else {
		profile = model.EmployerProfile{UserID: userID}
	}

	// Cập nhật thông tin từ request
	profile.CompanyName = profileData.CompanyName
	profile.CompanySize = profileData.CompanySize
	profile.Industry = profileData.Industry
	profile.Website = profileData.Website
	profile.Founded = profileData.Founded
	profile.About = profileData.About
	profile.Mission = profileData.Mission
	profile.PhoneNumber = profileData.PhoneNumber
	profile.Email = profileData.Email
	profile.Address = profileData.Address
	profile.City = profileData.City
	profile.Country = profileData.Country
	profile.FacebookURL = profileData.FacebookURL
	profile.TwitterURL = profileData.TwitterURL
	profile.LinkedinURL = profileData.LinkedinURL
	profile.TaxCode = profileData.TaxCode
	profile.BusinessLicense = profileData.BusinessLicense
	profile.ContactPersonName = profileData.ContactPersonName
	profile.ContactPersonRole = profileData.ContactPersonRole

	// Chuyển đổi các mảng thành JSON
	if len(profileData.Benefits) > 0 {
		benefitsJSON, err := json.Marshal(profileData.Benefits)
		if err != nil {
			middleware.Log(fmt.Errorf("Failed to marshal benefits: %v", err))
			return nil, message.InternalServerError
		}
		profile.Benefits = sql.NullString{String: string(benefitsJSON), Valid: true}
	}

	if len(profileData.Culture) > 0 {
		cultureJSON, err := json.Marshal(profileData.Culture)
		if err != nil {
			middleware.Log(fmt.Errorf("Failed to marshal culture: %v", err))
			return nil, message.InternalServerError
		}
		profile.Culture = sql.NullString{String: string(cultureJSON), Valid: true}
	}

	// Kiểm tra hồ sơ có đầy đủ thông tin chưa
	profile.ProfileComplete = isEmployerProfileComplete(&profile)

	// Lưu vào database
	err = repo.UpsertEmployerProfile(&profile)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to update employer profile: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return profile, nil
}

// UploadProfilePhoto xử lý upload ảnh cho hồ sơ
func UploadProfilePhoto(c *gin.Context, userID uint, photoType string, file *multipart.FileHeader) (interface{}, interface{}) {
	// Kiểm tra loại người dùng và loại ảnh
	userType, err := repo.GetUserType(userID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get user type: %v", err))
		return nil, message.UserNotFound
	}

	// Xác thực loại ảnh phù hợp với loại người dùng
	if (photoType == "resume" && userType != "jobseeker") ||
		((photoType == "company_logo" || photoType == "company_banner") && userType != "employer") {
		return nil, message.Message{Message: "Invalid photo type for user type", Code: 400}
	}

	// Kiểm tra định dạng file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	}

	// Cho phép file PDF cho hồ sơ
	if photoType == "resume" {
		allowedExtensions[".pdf"] = true
	}

	if !allowedExtensions[ext] {
		return nil, message.Message{Message: "Unsupported file type", Code: 400}
	}

	// Tạo tên file duy nhất
	filename := fmt.Sprintf("%d_%s_%d%s", userID, photoType, time.Now().Unix(), ext)

	// Tạo thư mục lưu trữ nếu chưa tồn tại
	uploadDir := filepath.Join("uploads", photoType)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		middleware.Log(fmt.Errorf("Failed to create upload directory: %v", err))
		return nil, message.InternalServerError
	}

	// Lưu file
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		middleware.Log(fmt.Errorf("Failed to save uploaded file: %v", err))
		return nil, message.InternalServerError
	}

	// Tạo URL cho file
	// Trong môi trường thực tế, bạn có thể sử dụng CDN hoặc S3
	photoURL := fmt.Sprintf("/uploads/%s/%s", photoType, filename)

	// Cập nhật URL vào database
	err = repo.UpdateProfilePhoto(userID, photoType, photoURL)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to update photo URL: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return map[string]string{"url": photoURL}, nil
}

// Helper functions to check if profiles are complete

func isJobseekerProfileComplete(profile *model.JobseekerProfile) bool {
	// Kiểm tra các trường bắt buộc
	return profile.DateOfBirth != nil &&
		profile.PhoneNumber != "" &&
		profile.Address != "" &&
		profile.City != "" &&
		profile.Country != "" &&
		profile.ProfileTitle != "" &&
		profile.About != "" &&
		profile.Skills != "" &&
		profile.Experience.Valid &&
		profile.Education.Valid
}

func isEmployerProfileComplete(profile *model.EmployerProfile) bool {
	// Kiểm tra các trường bắt buộc
	return profile.CompanyName != "" &&
		profile.CompanySize != "" &&
		profile.Industry != "" &&
		profile.About != "" &&
		profile.PhoneNumber != "" &&
		profile.Email != "" &&
		profile.Address != "" &&
		profile.City != "" &&
		profile.Country != "" &&
		profile.TaxCode != "" &&
		profile.ContactPersonName != ""
}
