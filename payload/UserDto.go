package payload

import (
	"datn_backend/domain/model"
	"datn_backend/payload/response"
	"time"
)

type UserDto struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	//Age        int           `json:"age" binding:"gte=1,lte=130"`
	Password string `json:"-" binding:"required,min=8"`
	//ResourceId *uint         `json:"resource_id"`
	Avatar string `json:"avatar"`
	//Status     int           `json:"status" `
	IsSupper     bool                             `json:"is_supper"`
	Id           uint                             `json:"id" `
	Phone        string                           `json:"phone"`
	FirstName    string                           `json:"first_name"`
	LastName     string                           `json:"last_name"`
	Roles        []*model.Role                    `json:"roles"`
	Providers    []*response.UserProviderResponse `json:"providers"`
	IsActive     bool                             `json:"is_active"`
	CreatedAt    time.Time                        `json:"created_at"`
	IsLocked     bool                             `json:"is_locked"`
	Gender       int                              `json:"gender"`
	DateOfBirth  *time.Time                       `json:"date_of_birth"`
	ReceivedNoti bool                             `json:"received_noti"`
}
