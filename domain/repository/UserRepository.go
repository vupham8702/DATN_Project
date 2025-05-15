package repository

import (
	"datn_backend/config"
	"datn_backend/domain/model"
	"datn_backend/middleware"
	"fmt"
)

// GetUserType lấy loại người dùng (jobseeker hoặc employer) dựa trên ID
func GetUserType(userID uint) (string, error) {
	var user model.User
	db := config.DB
	result := db.First(&user, userID)
	if result.Error != nil {
		middleware.Log(fmt.Errorf("Failed to get user by ID: %v", result.Error))
		return "", result.Error
	}

	// Kiểm tra loại người dùng dựa trên provider
	var provider model.UserProvider
	result = db.Where("user_id = ?", userID).First(&provider)
	if result.Error != nil {
		// Nếu không có provider, mặc định là jobseeker
		return config.USER_TYPE_JOBSEEKER, nil
	}

	if provider.UserType == config.USER_TYPE_EMPLOYER {
		return config.USER_TYPE_EMPLOYER, nil
	}

	return config.USER_TYPE_JOBSEEKER, nil
}

// GetPendingEmployers lấy danh sách nhà tuyển dụng đang chờ phê duyệt
func GetPendingEmployers() ([]*model.EmployerProfile, error) {
	var profiles []*model.EmployerProfile
	db := config.DB

	// Join với bảng user_providers để lấy các nhà tuyển dụng chưa được phê duyệt
	result := db.Preload("User").
		Joins("JOIN user_provider ON employer_profile.user_id = user_provider.user_id").
		Where("user_provider.user_type = ?", config.USER_TYPE_EMPLOYER).
		Where("user_provider.is_approved = ?", false).
		Find(&profiles)

	if result.Error != nil {
		middleware.Log(fmt.Errorf("Failed to get pending employers: %v", result.Error))
		return nil, result.Error
	}

	return profiles, nil
}
