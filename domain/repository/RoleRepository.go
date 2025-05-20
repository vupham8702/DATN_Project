package repository

import (
	"datn_backend/domain/model"
	"datn_backend/message"
	"datn_backend/middleware"
	"gorm.io/gorm"
)

func GetRoleByID(db *gorm.DB, id uint) (*model.Role, interface{}) {
	var role model.Role
	result := db.Preload("Resource").
		Preload("Permissions").
		First(&role, id)
	if result.Error != nil {
		middleware.Log(result.Error)
		return nil, message.NotFound
	}
	return &role, nil
}
