package service

import (
	"context"
	"datn_backend/config"
	"datn_backend/domain/model"
	repo "datn_backend/domain/repository"
	"datn_backend/message"
	"datn_backend/middleware"
	"datn_backend/payload"
	"datn_backend/payload/response"
	"datn_backend/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
	"time"
)

func Login(c *gin.Context, userLogin *payload.UserLogin, deviceId string) (*response.UserToken, interface{}) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(userLogin.Username))

	// Log login attempt
	middleware.Log(fmt.Sprintf("Login attempt for email: %s, device: %s", email, deviceId))

	var user model.User
	userMail, err := repo.GetUserByMail(email)
	if err != nil {
		middleware.Log(fmt.Sprintf("Login failed: Email not found: %s", email))
		return nil, message.EmailNotExist
	}

	user = *userMail
	if &user == nil {
		middleware.Log(fmt.Sprintf("Login failed: User object is nil for email: %s", email))
		return nil, message.EmailNotExist
	}

	// Check if email is verified
	if user.IsActive == false {
		middleware.Log(fmt.Sprintf("Login failed: Email not verified: %s", email))
		return nil, message.EmailNotVerified
	}

	if user.IsLocked == true {
		middleware.Log(fmt.Sprintf("Login failed: Account locked: %s", email))
		return nil, message.UserHasBeenLocked
	}
	provider := repo.GetUserProvider(user)

	// Check if employer account is approved
	if provider.UserType == "employer" && !provider.IsApproved {
		middleware.Log(fmt.Sprintf("Login failed: Employer account not approved yet: %s", email))
		return nil, message.ApprovalAccountPenning
	}

	verify, _, err := utils.VerifyPassword(userLogin.Password, user.Password)
	if !verify || err != nil {
		middleware.Log(fmt.Sprintf("Login failed: Incorrect password for email: %s", email))
		return nil, message.PasswordNotCorrect
	}
	userType := repo.GetUserProviderByUserID(&user)

	token, tokenErr := CreateToken(c, &user, userType)
	if tokenErr != nil {
		middleware.Log(fmt.Errorf("Failed to create token for user %s: %v", email, tokenErr))
		return nil, tokenErr
	}

	// Log successful login
	middleware.Log(fmt.Sprintf("Login successful for user: %s (ID: %d, Type: %s)", email, user.ID, provider.UserType))

	return token, nil
}

func CreateToken(c *gin.Context, user *model.User, userType string) (*response.UserToken, interface{}) {
	var roles []string

	for _, role := range user.Roles {
		roles = append(roles, fmt.Sprintf("%d", role.ID))
	}

	token := utils.GenerateToken(
		user.ID,
		user.IsSupper,
		roles,
		userType,
	)
	uidStr := fmt.Sprintf("%d", user.ID)
	errSaveToken := CreateTokenRedis(c, &token, uidStr)
	if errSaveToken != nil {
		return nil, message.ExcuteDatabaseError
	}

	return &token, nil
}

func CreateTokenRedis(c *gin.Context, token *response.UserToken, uid string) interface{} {
	key := config.TOKEN + ":" + uid
	value, err := json.Marshal(token)
	if err != nil {
		return message.InternalServerError
	}
	status := config.RedisClient.Set(c, key, value, 0)
	if status.Val() != config.OK {
		middleware.Log(fmt.Errorf("Save token error Redis ...."))
		return nil
	}
	return nil
}

// Register handles user registration
func Register(c *gin.Context, userRegister *payload.UserRegister) (interface{}, interface{}) {
	// Normalize email (convert to lowercase)
	email := strings.ToLower(strings.TrimSpace(userRegister.Email))

	// Log registration attempt
	middleware.Log(fmt.Sprintf("Registration attempt for email: %s, type: %s", email, userRegister.UserType))

	// Check if email already exists
	existingUser, err := repo.GetUserByMail(email)
	if err == nil && existingUser != nil {
		middleware.Log(fmt.Sprintf("Registration failed: Email already exists: %s", email))
		return nil, message.EmailAlreadyExists
	}

	// Validate password strength
	if !utils.ValidatePassword(userRegister.Password) {
		middleware.Log(fmt.Sprintf("Registration failed: Password requirements not met for email: %s", email))
		return nil, message.PasswordRequirements
	}

	// Validate user type
	if userRegister.UserType != config.USER_TYPE_EMPLOYER && userRegister.UserType != config.USER_TYPE_JOBSEEKER {
		middleware.Log(fmt.Sprintf("Registration failed: Invalid user type: %s", userRegister.UserType))
		return nil, message.Message{Message: "Invalid user type. Must be 'jobseeker' or 'employer'", Code: 400}
	}

	// Create new user with more default values
	user := model.User{
		Username:  email,
		Email:     email,
		FirstName: userRegister.FullName,
		Password:  utils.HashPassword(userRegister.Password),
		IsActive:  false, // Not active until email is verified
		IsSupper:  false,
		IsLocked:  false,
	}

	// Save user to database using repository function
	if err := repo.CreateUser(&user, userRegister.UserType); err != nil {
		middleware.Log(fmt.Errorf("Failed to create user: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	// Generate verification token
	verificationToken := uuid.New().String()

	key := fmt.Sprintf("email_verification:%s", email)
	// Use context with timeout for Redis operations
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	err = config.RedisClient.Set(ctx, key, verificationToken, 24*time.Hour).Err()
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to store verification token in Redis: %v", err))
		return nil, message.InternalServerError
	}

	if userRegister.UserType == "employer" {
		return user, message.RegistrationSuccess
	}

	// Log successful registration
	middleware.Log(fmt.Sprintf("User registered successfully: %s (ID: %d, Type: %s)", email, user.ID, userRegister.UserType))

	return user, nil
}

// VerifyEmail verifies a user's email using the token
func VerifyEmail(c *gin.Context, verifyEmail *payload.VerifyEmail) (interface{}, interface{}) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(verifyEmail.Email))

	// Log verification attempt
	middleware.Log(fmt.Sprintf("Email verification attempt for: %s", email))

	// Get token from Redis with timeout context
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("email_verification:%s", email)
	storedToken, err := config.RedisClient.Get(ctx, key).Result()

	if err != nil {
		middleware.Log(fmt.Errorf("Verification failed: Token not found for email %s: %v", email, err))
		return nil, message.InvalidVerifyToken
	}

	if storedToken != verifyEmail.Token {
		middleware.Log(fmt.Sprintf("Verification failed: Token mismatch for email %s", email))
		return nil, message.InvalidVerifyToken
	}

	// Find user by email
	user, err := repo.GetUserByMail(email)
	if err != nil {
		middleware.Log(fmt.Errorf("Verification failed: User not found for email %s: %v", email, err))
		return nil, message.EmailNotExist
	}

	// Update user to active using repository function
	if err := repo.UpdateUserActiveStatus(user, true); err != nil {
		middleware.Log(fmt.Errorf("Failed to update user status: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	// Delete token from Redis
	config.RedisClient.Del(ctx, key)

	// Log successful verification
	middleware.Log(fmt.Sprintf("Email verified successfully for: %s (ID: %d)", email, user.ID))

	return message.EmailVerifySuccess, nil
}

// UpdateUserProviderApprovalStatus updates the approval status of a user provider
func UpdateUserProviderApprovalStatus(userProviderID uint, isApproved bool, approvedBy uint, note string) error {
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		middleware.Log(fmt.Errorf("Failed to begin transaction: %v", tx.Error))
		return tx.Error
	}

	// Find the user provider
	var userProvider model.UserProvider
	if err := tx.First(&userProvider, userProviderID).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to find user provider: %v", err))
		return err
	}

	// Update approval status
	userProvider.IsApproved = isApproved
	userProvider.ApprovedBy = approvedBy
	userProvider.ApprovalNote = note
	if err := tx.Save(&userProvider).Error; err != nil {
		tx.Rollback()
		middleware.Log(fmt.Errorf("Failed to update user provider: %v", err))
		return err
	}

	// If approved, also activate the user
	if isApproved {
		var user model.User
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
