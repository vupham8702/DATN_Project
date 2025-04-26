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

	token, tokenErr := CreateToken(c, &user, provider.UserType, deviceId)
	if tokenErr != nil {
		middleware.Log(fmt.Errorf("Failed to create token for user %s: %v", email, tokenErr))
		return nil, tokenErr
	}

	// Log successful login
	middleware.Log(fmt.Sprintf("Login successful for user: %s (ID: %d, Type: %s)", email, user.ID, provider.UserType))

	return token, nil
}

func CreateToken(c *gin.Context, user *model.User, userType string, deviceId string) (*response.UserToken, interface{}) {
	var roles []string

	for _, role := range user.Roles {
		roles = append(roles, fmt.Sprintf("%d", role.ID))
	}

	token := utils.GenerateToken(
		user.ID,
		user.IsSupper,
		roles,
		userType,
		deviceId,
	)
	uidStr := fmt.Sprintf("%d", user.ID)
	errSaveToken := CreateTokenRedis(c, &token, uidStr, deviceId)
	if errSaveToken != nil {
		return nil, message.ExcuteDatabaseError
	}

	return &token, nil
}

func CreateTokenRedis(c *gin.Context, token *response.UserToken, uid string, deviceId string) interface{} {
	key := config.TOKEN + ":" + uid + ":" + deviceId
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
	tx := config.DB.Begin()
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

	defaultRole, errRole := repo.GetRoleByName(config.DB, config.DEFAULT_USER_PERMISSION)
	if errRole != nil {
		return nil, defaultRole
	}
	roles := []*model.Role{defaultRole}
	for _, rdto := range userRegister.Roles {
		extraRole, err := repo.GetRoleByID(tx, rdto.Id)
		if err != nil {
			tx.Rollback()
			return nil, message.RoleNotFound
		}
		if extraRole.ID != defaultRole.ID {
			roles = append(roles, extraRole)
		}
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
		Roles:     roles,
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
	middleware.Log(fmt.Sprintf("User registered successfully: %s (ID: %d, Type: %s, Role: DEFAULT_USER)", email, user.ID, userRegister.UserType))

	return user, nil
}

// ApproveEmployer handles the approval or rejection of an employer account
func ApproveEmployer(c *gin.Context, approveRequest *payload.ApproveEmployer, adminID *uint) (interface{}, interface{}) {
	// Get the user by ID
	var user model.User
	if err := config.DB.First(&user, approveRequest.UserID).Error; err != nil {
		middleware.Log(fmt.Errorf("User not found: %v", err))
		return nil, message.UserNotFound
	}

	// Get the user provider
	userProvider, err := repo.GetUserProviderByUserID(&user)
	if err != nil {
		middleware.Log(fmt.Errorf("User provider not found: %v", err))
		return nil, message.Message{Message: "Không tìm thấy thông tin tài khoản", Code: 404}
	}

	// Check if this is actually an employer account
	if userProvider.UserType != config.USER_TYPE_EMPLOYER {
		return nil, message.Message{Message: "Tài khoản này không phải tài khoản của nhà tuyển dụng", Code: 400}
	}

	// Check if already approved
	if userProvider.IsApproved {
		return nil, message.Message{Message: "Tài khoản này đã được kiểm duyệt", Code: 400}
	}

	// Process approval/rejection
	isApproved := approveRequest.Status == "approved"
	err = repo.UpdateUserProviderApprovalStatus(userProvider.ID, isApproved, adminID, approveRequest.Note)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to update approval status: %v", err))
		return nil, message.InternalServerError
	}

	// Prepare notification to the user
	// TODO: Send email notification to the user about approval status

	if isApproved {
		return map[string]interface{}{
			"message": "Employer account has been approved",
			"userId":  user.ID,
		}, nil
	} else {
		return map[string]interface{}{
			"message": "Employer account has been rejected",
			"userId":  user.ID,
		}, nil
	}
}

// GetPendingEmployers retrieves all employer accounts pending approval
func GetPendingEmployers() (interface{}, interface{}) {
	var userProviders []model.UserProvider
	db := config.DB

	// Get all user providers that are employer type and not approved
	result := db.Where("user_type = ? AND is_approved = ?", "employer", false).
		Preload("User"). // Load associated user information
		Find(&userProviders)

	if result.Error != nil {
		middleware.Log(fmt.Errorf("Failed to retrieve pending employers: %v", result.Error))
		return nil, message.InternalServerError
	}

	// Format data for response
	var response []map[string]interface{}
	for _, provider := range userProviders {
		userData := map[string]interface{}{
			"userId":     provider.UserID,
			"email":      provider.Email,
			"providerId": provider.ID,
			"fullName":   provider.User.FirstName + " " + provider.User.LastName,
			"createdAt":  provider.CreatedAt,
		}
		response = append(response, userData)
	}

	return map[string]interface{}{
		"pendingEmployers": response,
		"count":            len(response),
	}, nil
}
