package repository

import (
	"datn_backend/config"
	"datn_backend/domain/model"
	"datn_backend/middleware"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// CreatePost tạo mới một bài đăng tuyển dụng
func CreatePost(post *model.PostJob) error {
	db := config.DB

	// Thực hiện lưu vào database
	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to create job post: %v", err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	return nil
}

// GetPostByID lấy thông tin bài đăng theo ID
func GetPostByID(postID uint) (*model.PostJob, error) {
	var post model.PostJob
	db := config.DB

	result := db.Where("id = ? AND is_deleted = ?", postID, false).First(&post)
	if result.Error != nil {
		return nil, result.Error
	}

	return &post, nil
}

// GetAllPosts lấy danh sách tất cả các bài đăng
func GetAllPosts(page, pageSize int, status string) ([]*model.PostJob, int64, error) {
	var posts []*model.PostJob
	var total int64
	db := config.DB

	query := db.Model(&model.PostJob{}).Where("is_deleted = ? and is_approved = ?", false, true)

	// Lọc theo trạng thái nếu được chỉ định
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Đếm tổng số bản ghi
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Phân trang
	offset := (page - 1) * pageSize
	result := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return posts, total, nil
}

// GetPostsByEmployer lấy danh sách bài đăng của một nhà tuyển dụng
func GetPostsByEmployer(employerID uint, page, pageSize int) ([]*model.PostJob, int64, error) {
	var posts []*model.PostJob
	var total int64
	db := config.DB

	query := db.Model(&model.PostJob{}).
		Where("created_by = ? AND is_deleted = ?", employerID, false)

	// Đếm tổng số bản ghi
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Phân trang
	offset := (page - 1) * pageSize
	result := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return posts, total, nil
}

// UpdatePost cập nhật thông tin bài đăng
func UpdatePost(post *model.PostJob) error {
	db := config.DB

	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	// Cập nhật thời gian cập nhật
	post.UpdatedAt = time.Now()

	if err := tx.Save(post).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to update job post: %v", err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	return nil
}

// DeletePost xóa mềm một bài đăng
func DeletePost(postID uint, userID uint) error {
	db := config.DB

	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	// Thực hiện xóa mềm bằng cách cập nhật trường is_deleted
	result := tx.Model(&model.PostJob{}).
		Where("id = ?", postID).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"updated_at": time.Now(),
			"updated_by": userID,
		})

	if result.Error != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to delete job post: %v", result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	return nil
}

// UpdatePostStatus cập nhật trạng thái của bài đăng
func UpdatePostStatus(postID uint, isApproved bool, userID uint) error {
	db := config.DB

	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	result := tx.Model(&model.PostJob{}).
		Where("id = ? AND is_deleted = ?", postID, false).
		Updates(map[string]interface{}{
			"is_approve": isApproved,
			"updated_at": time.Now(),
			"updated_by": userID,
		})

	if result.Error != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to update job post status: %v", result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	return nil
}
func CheckIfUserApplied(uid uint, pid uint) (bool, interface{}) {
	var count int64
	err := config.DB.Model(&model.JobApplication{}).Where("created_by = ? and post_job_id = ? and is_deleted = false", uid, pid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func IncrementApplicationCountRaw(postID uint) error {
	db := config.DB

	// Bắt đầu transaction
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("không thể bắt đầu transaction: %w", tx.Error)
	}

	// Câu lệnh UPDATE sử dụng COUNT(*) trên bảng job_applications
	result := tx.Exec(`
        UPDATE post_job
        SET applications_count = (
            SELECT COUNT(*)
            FROM post_job_application
            WHERE post_job_id = ? AND is_deleted = FALSE
        )
        WHERE id = ?
    `, postID, postID)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("lỗi khi cập nhật application_count: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("không tìm thấy post_job với id = %d", postID)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("lỗi khi commit transaction: %w", err)
	}
	return nil
}
