package repository

import (
	"datn_backend/config"
	m "datn_backend/domain/model"
	"datn_backend/message"
	"datn_backend/middleware"
	"fmt"
	"gorm.io/gorm"
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

func GetUserProviderByUserID(user *m.User) (*m.UserProvider, interface{}) {
	var userProvider m.UserProvider
	db := config.DB
	result := db.Where("user_id = ?", user.ID).First(&userProvider)
	if result.Error != nil {
		return nil, result.Error
	}

	return &userProvider, nil
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

	// Save user to database within transaction first to get ID
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to create user: %v", err))
		return err
	}

	// Now create UserProvider with the user ID
	userProvider := &m.UserProvider{
		Email:            user.Email,
		UserID:           uint(user.ID),
		ProviderIdentify: user.Email,
		UserType:         userType,
		ReceivedNoti:     true, // Default to receive notifications
	}

	// Fix duplicate conditional logic
	if userType == config.USER_TYPE_EMPLOYER {
		userProvider.IsApproved = false // Employer needs approval
	} else if userType == config.USER_TYPE_JOBSEEKER {
		userProvider.IsApproved = true // Jobseeker doesn't need approval
	}

	if err := tx.Create(userProvider).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to create user provider: %v", err))
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

// UpdateUserProviderApprovalStatus updates the approval status of a user provider
func UpdateUserProviderApprovalStatus(userProviderID uint, isApproved bool, approvedBy *uint, note string) error {
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	// Find the user provider
	var userProvider m.UserProvider
	if err := tx.First(&userProvider, userProviderID).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to find user provider: %v", err))
		return err
	}

	// Update approval status
	userProvider.IsApproved = isApproved
	userProvider.ApprovedBy = *approvedBy
	userProvider.ApprovalNote = note
	if err := tx.Save(&userProvider).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to update user provider: %v", err))
		return err
	}

	// If approved, also activate the user
	if isApproved {
		var user m.User
		if err := tx.First(&user, userProvider.UserID).Error; err != nil {
			tx.Rollback()
			middleware.Log(fmt.Errorf("Failed to find user: %v", err))
			return err
		}

		// Only activate if the email has been verified (check your logic here)
		if user.IsActive {
			user.IsLocked = false // Ensure the user is unlocked
			if err := tx.Save(&user).Error; err != nil {
				tx.Rollback()
				middleware.Log(fmt.Errorf("Failed to update user: %v", err))
				return err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		middleware.Log(fmt.Errorf("Failed to commit transaction: %v", err))
		return err
	}

	return nil
}

func GetRoleByName(db *gorm.DB, name string) (*m.Role, *message.Message) {
	var role m.Role
	result := db.Where("(name = ?)", name).First(&role)
	if result.Error != nil {
		middleware.Log(result.Error)
		return nil, &message.NotFound
	}
	return &role, nil
}
