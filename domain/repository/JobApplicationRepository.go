package repository

import (
	"datn_backend/config"
	m "datn_backend/domain/model"
	"errors"
	"gorm.io/gorm"
)

// GetJobApplicationsByUserID lấy danh sách ứng tuyển của người dùng
func GetJobApplicationsByUserID(userID uint) ([]*m.JobApplication, error) {
	var applications []*m.JobApplication
	db := config.DB

	result := db.
		//Preload("PostJob").
		//Preload("CV").
		Where("created_by = ? AND is_deleted = ?", userID, false).
		Order("created_at DESC").
		Find(&applications)

	if result.Error != nil {
		return nil, result.Error
	}

	return applications, nil
}

// GetJobApplicationsByPostID lấy danh sách ứng tuyển cho bài đăng
func GetJobApplicationsByPostID(postJobID uint) ([]*m.JobApplication, error) {
	var applications []*m.JobApplication
	db := config.DB

	result := db.Preload("User").
		Preload("CV").
		Where("post_job_id = ? AND is_deleted = ?", postJobID, false).
		Order("created_at DESC").
		Find(&applications)

	if result.Error != nil {
		return nil, result.Error
	}

	return applications, nil
}

// GetJobApplicationByID lấy thông tin ứng tuyển theo ID
func GetJobApplicationByID(id uint) (*m.JobApplication, error) {
	var application m.JobApplication
	db := config.DB

	result := db.Preload("User").
		Preload("PostJob").
		Preload("CV").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&application)

	if result.Error != nil {
		return nil, result.Error
	}

	return &application, nil
}

// UpdateJobApplicationStatus cập nhật trạng thái ứng tuyển
func UpdateJobApplicationStatus(id uint, status string) error {
	db := config.DB

	result := db.Model(&m.JobApplication{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("application not found")
	}

	return nil
}

// CheckIfUserOwnsPost kiểm tra xem người dùng có sở hữu bài đăng không
func CheckIfUserOwnsPost(userID uint, postJobID uint) (bool, error) {
	var post m.PostJob
	db := config.DB

	result := db.Where("id = ? AND created_by = ? AND is_deleted = ?", postJobID, userID, false).
		First(&post)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

//// CheckIfUserOwnsApplication kiểm tra xem người dùng có sở hữu đơn ứng tuyển không
//func CheckIfUserOwnsApplication(userID uint, applicationID uint) (bool, error) {
//	var application m.JobApplication
//	db := config.DB
//
//	result := db.Where("id = ? AND user_id = ? AND is_deleted = ?", applicationID, userID, false).
//		First(&application)
//
//	if result.Error != nil {
//		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//			return false, nil
//		}
//		return false, result.Error
//	}
//
//	return true, nil
//}

// GetPostOwnerByApplicationID lấy ID người sở hữu bài đăng từ ID đơn ứng tuyển
func GetPostOwnerByApplicationID(applicationID uint) (uint, error) {
	var application m.JobApplication
	db := config.DB

	// Lấy thông tin đơn ứng tuyển
	result := db.Where("id = ? AND is_deleted = ?", applicationID, false).
		First(&application)

	if result.Error != nil {
		return 0, result.Error
	}

	// Lấy thông tin bài đăng
	var post m.PostJob
	result = db.Where("id = ? AND is_deleted = ?", application.PostJobID, false).
		First(&post)

	if result.Error != nil {
		return 0, result.Error
	}

	return post.CreatedBy, nil
}
