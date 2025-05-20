package payload

type UserRegister struct {
	FullName string     `json:"fullName" binding:"required"`
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=8"`
	UserType string     `json:"userType" binding:"required"` // "jobseeker" or "employer"
	Roles    []*RoleDto `json:"roles"`
}
