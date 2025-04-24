package repository

import (
	"datn_backend/config"
	m "datn_backend/domain/model"
	"datn_backend/message"
	"datn_backend/middleware"
	"fmt"
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

func GetUserProviderByUserID(user *m.User) string {
	var userProvider m.UserProvider
	db := config.DB
	result := db.Where("user_id = ?", user.ID).First(&userProvider)
	if result.Error != nil {
		return ""
	}
	userType := userProvider.UserType
	return userType
}

// CreateUser creates a new user in the database using a transaction
func CreateUser(user *m.User, userType string) error {
	// Start a database transaction
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	userProvider := &m.UserProvider{
		Email:            user.Email,
		UserID:           uint(user.ID),
		ProviderIdentify: user.Email,
		UserType:         userType,
		ReceivedNoti:     true, // Default to receive notifications
	}
	if userType == config.USER_TYPE_EMPLOYER {
		userProvider.IsApproved = true
	} else if userType == config.USER_TYPE_JOBSEEKER {
		userProvider.IsApproved = false
	}

	if err := tx.Create(userProvider).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to create user provider: %v", err))
		return err
	}

	// Save user to database within transaction
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to create user: %v", err))
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	return nil
}

// UpdateUserActiveStatus updates the IsActive status of a user
func UpdateUserActiveStatus(user *m.User, isActive bool) error {
	// Start transaction for updating user
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	// Update user active status
	user.IsActive = isActive
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to update user status: %v", err))
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	return nil
}

func GetUserProvider(user m.User) m.UserProvider {
	db := config.DB
	query := db.Select("*").Where("user_id = ?", user.ID).Find(&[]m.UserProvider{})
	if query.Error != nil {
		panic(query.Error)
	}
	return m.UserProvider{}
}
