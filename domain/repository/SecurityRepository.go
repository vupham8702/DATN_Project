package repository

import (
	"datn_backend/config"
	m "datn_backend/domain/model"
	"datn_backend/message"
	"datn_backend/middleware"
)

func GetUserByPhone(phone string) (*m.User, interface{}) {
	db := config.DB
	var user m.User
	result := db.Preload("Providers").
		Preload("Roles", "is_deleted = ? ", false).
		Where("phone = ? and is_active = ? and is_deleted = ?", phone, true, false).
		First(&user)
	if result.RowsAffected == 0 {
		middleware.Log(result.Error)
		return nil, message.UserNotFound
	}
	return &user, nil
}

func GetUserByMail(email string) (*m.User, error) {
	var user m.User
	db := config.DB
	result := db.Preload("Providers").Preload("Roles", "is_deleted = ? ", false).Where("email = ? ", email).
		Where("is_deleted = ?", false).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
