package repository

import (
	"datn_backend/config"
	m "datn_backend/domain/model"
	"datn_backend/middleware"
	"fmt"
	"gorm.io/gorm"
	"math"
	"time"
)

// GetJobseekerProfileByUserID gets a jobseeker profile by user ID
func GetJobseekerProfileByUserID(userID *uint) (*m.JobseekerProfile, error) {
	var profile m.JobseekerProfile
	db := config.DB
	result := db.Where("user_id = ?", userID).First(&profile)

	if result.Error != nil {
		return nil, result.Error
	}

	return &profile, nil
}

// GetEmployerProfileByUserID gets an employer profile by user ID
func GetEmployerProfileByUserID(userID *uint) (*m.EmployerProfile, error) {
	var profile m.EmployerProfile
	db := config.DB
	result := db.Where("user_id = ?", userID).First(&profile)

	if result.Error != nil {
		return nil, result.Error
	}

	return &profile, nil
}

// CreateJobseekerProfile creates a new jobseeker profile
func CreateJobseekerProfile(profile *m.JobseekerProfile) error {
	db := config.DB

	var user m.User
	maxRetries := 5
	retryDelay := time.Millisecond * 100

	for i := 0; i < maxRetries; i++ {
		if err := db.First(&user, profile.UserID).Error; err == nil {
			// User đã tồn tại, tiếp tục tạo profile
			break
		}

		if i == maxRetries-1 {
			// Đã thử tối đa số lần cho phép
			return fmt.Errorf("user with ID %d does not exist after %d attempts", profile.UserID, maxRetries)
		}
		// Chờ một chút trước khi thử lại
		time.Sleep(retryDelay)
		retryDelay *= 2
	}

	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	if err := tx.First(&user, profile.UserID).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("User %d disappeared before creating profile: %v", profile.UserID, err))
		return err
	}

	if err := tx.Create(profile).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to create jobseeker profile: %v", err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	return nil
}

// CreateEmployerProfile creates a new employer profile
func CreateEmployerProfile(profile *m.EmployerProfile) error {
	db := config.DB

	var user m.User
	maxRetries := 5
	retryDelay := time.Millisecond * 500 // Bắt đầu với 500ms

	for i := 0; i < maxRetries; i++ {
		err := db.First(&user, profile.UserID).Error

		if err == nil {
			// User đã tồn tại, tiếp tục tạo profile
			middleware.Log(fmt.Sprintf("Found user %d on attempt %d. Creating employer profile...", profile.UserID, i+1))
			break
		}

		if i == maxRetries-1 {
			// Đã thử tối đa số lần cho phép
			middleware.Log(fmt.Errorf("user with ID %d does not exist after %d attempts", profile.UserID, maxRetries))
			return fmt.Errorf("user with ID %d does not exist after %d attempts", profile.UserID, maxRetries)
		}

		middleware.Log(fmt.Sprintf("User %d not found yet, retrying in %v (attempt %d/%d)",
			profile.UserID, retryDelay, i+1, maxRetries))

		// Chờ một chút trước khi thử lại
		time.Sleep(retryDelay)
		// Tăng thời gian chờ theo cấp số nhân nhưng không quá 2 giây
		retryDelay = time.Duration(math.Min(float64(retryDelay)*2, float64(2*time.Second)))
	}

	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	if err := tx.First(&user, profile.UserID).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("User %d disappeared before creating profile: %v", profile.UserID, err))
		return err
	}

	if err := tx.Create(profile).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to create employer profile: %v", err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	middleware.Log(fmt.Sprintf("Successfully created employer profile for user %d", profile.UserID))
	return nil
}

// UpdateJobseekerProfile updates a jobseeker profile
func UpdateJobseekerProfile(profile *m.JobseekerProfile) error {
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	if err := tx.Save(profile).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to update jobseeker profile: %v", err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	return nil
}

// UpdateEmployerProfile updates an employer profile
func UpdateEmployerProfile(profile *m.EmployerProfile) error {
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	if err := tx.Save(profile).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to update employer profile: %v", err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	return nil
}

// GetAllEmployerProfiles gets all employer profiles
func GetAllEmployerProfiles(page, pageSize int) ([]*m.EmployerProfile, int64, error) {
	var profiles []*m.EmployerProfile
	var total int64
	db := config.DB

	// Count total records
	if err := db.Model(&m.EmployerProfile{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * pageSize
	result := db.Preload("User").
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&profiles)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return profiles, total, nil
}

// GetVerifiedEmployerProfiles gets all verified employer profiles
func GetVerifiedEmployerProfiles(page, pageSize int) ([]*m.EmployerProfile, int64, error) {
	var profiles []*m.EmployerProfile
	var total int64
	db := config.DB

	// Count total verified records
	if err := db.Model(&m.EmployerProfile{}).
		Where("verification_status = ?", "verified").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated verified records
	offset := (page - 1) * pageSize
	result := db.Preload("User").
		Where("verification_status = ?", "verified").
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&profiles)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return profiles, total, nil
}

// GetPendingVerificationEmployerProfiles gets all employer profiles pending verification
func GetPendingVerificationEmployerProfiles() ([]*m.EmployerProfile, error) {
	var profiles []*m.EmployerProfile
	db := config.DB

	result := db.Preload("User").
		Where("verification_status = ?", "pending").
		Order("created_at ASC").
		Find(&profiles)

	if result.Error != nil {
		return nil, result.Error
	}

	return profiles, nil
}

// UpdateProfilePhoto updates a profile photo URL
func UpdateProfilePhoto(userID uint, photoType string, photoURL string) error {
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Get user type
	userType, err := GetUserType(userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update the appropriate profile based on user type and photo type
	if userType == "jobseeker" {
		var profile m.JobseekerProfile
		if err := tx.Where("user_id = ?", userID).First(&profile).Error; err != nil {
			// If profile doesn't exist, create it
			if err == gorm.ErrRecordNotFound {
				profile = m.JobseekerProfile{UserID: userID}
				if photoType == "profile_picture" {
					profile.ProfilePicture = photoURL
				} else if photoType == "resume" {
					profile.ResumeURL = photoURL
				}
				if err := tx.Create(&profile).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				tx.Rollback()
				return err
			}
		} else {
			// Update existing profile
			if photoType == "profile_picture" {
				profile.ProfilePicture = photoURL
			} else if photoType == "resume" {
				profile.ResumeURL = photoURL
			}
			if err := tx.Save(&profile).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else if userType == "employer" {
		var profile m.EmployerProfile
		if err := tx.Where("user_id = ?", userID).First(&profile).Error; err != nil {
			// If profile doesn't exist, create it
			if err == gorm.ErrRecordNotFound {
				profile = m.EmployerProfile{UserID: userID}
				if photoType == "company_logo" {
					profile.CompanyLogo = photoURL
				} else if photoType == "company_banner" {
					profile.CompanyBanner = photoURL
				}
				if err := tx.Create(&profile).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				tx.Rollback()
				return err
			}
		} else {
			// Update existing profile
			if photoType == "company_logo" {
				profile.CompanyLogo = photoURL
			} else if photoType == "company_banner" {
				profile.CompanyBanner = photoURL
			}
			if err := tx.Save(&profile).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		tx.Rollback()
		return fmt.Errorf("invalid user type: %s", userType)
	}

	return tx.Commit().Error
}

// UpsertJobseekerProfile creates or updates a jobseeker profile
func UpsertJobseekerProfile(profile *m.JobseekerProfile) error {
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var existingProfile m.JobseekerProfile
	result := tx.Where("user_id = ?", profile.UserID).First(&existingProfile)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new profile
			if err := tx.Create(profile).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return result.Error
		}
	} else {
		// Update existing profile
		if err := tx.Model(&existingProfile).Updates(profile).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// UpsertEmployerProfile creates or updates an employer profile
func UpsertEmployerProfile(profile *m.EmployerProfile) error {
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var existingProfile m.EmployerProfile
	result := tx.Where("user_id = ?", profile.UserID).First(&existingProfile)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new profile
			if err := tx.Create(profile).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return result.Error
		}
	} else {
		// Update existing profile
		if err := tx.Model(&existingProfile).Updates(profile).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
