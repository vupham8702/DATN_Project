package repository

import (
	"datn_backend/config"
	m "datn_backend/domain/model"
	"errors"
	"gorm.io/gorm"
)

// GetAllCVTemplates lấy tất cả mẫu CV
func GetAllCVTemplates() ([]*m.CVTemplate, error) {
	var templates []*m.CVTemplate
	db := config.DB

	result := db.Where("is_deleted = ?", false).
		Order("created_at DESC").
		Find(&templates)

	if result.Error != nil {
		return nil, result.Error
	}

	return templates, nil
}

// GetCVTemplateByID lấy mẫu CV theo ID
func GetCVTemplateByID(id uint) (*m.CVTemplate, error) {
	var template m.CVTemplate
	db := config.DB

	result := db.Where("id = ? AND is_deleted = ?", id, false).
		First(&template)

	if result.Error != nil {
		return nil, result.Error
	}

	return &template, nil
}

// CreateCVTemplate tạo mẫu CV mới
func CreateCVTemplate(template *m.CVTemplate) error {
	db := config.DB

	result := db.Create(template)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateCVTemplate cập nhật mẫu CV
func UpdateCVTemplate(template *m.CVTemplate) error {
	db := config.DB

	result := db.Save(template)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteCVTemplate xóa mềm mẫu CV
func DeleteCVTemplate(id uint) error {
	db := config.DB

	result := db.Model(&m.CVTemplate{}).
		Where("id = ?", id).
		Update("is_deleted", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("template not found")
	}

	return nil
}

// IncrementDownloadCount tăng số lượt tải xuống của mẫu CV
func IncrementDownloadCount(id uint) error {
	db := config.DB

	result := db.Model(&m.CVTemplate{}).
		Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + ?", 1))

	if result.Error != nil {
		return result.Error
	}

	return nil
}
