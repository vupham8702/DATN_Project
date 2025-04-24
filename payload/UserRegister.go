package payload

type UserRegister struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	UserType string `json:"userType" binding:"required"` // "jobseeker" or "employer"
}

type VerifyEmail struct {
	Token string `json:"token" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type ResendVerification struct {
	Email string `json:"email" binding:"required,email"`
}

// ApproveEmployer is used by admins to approve employer accounts
type ApproveEmployer struct {
	UserID uint   `json:"userId" binding:"required"`
	Status string `json:"status" binding:"required"` // "approved" or "rejected"
	Note   string `json:"note"`
}
