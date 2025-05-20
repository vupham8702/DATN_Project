package payload

import (
	"datn_backend/domain/model"
	"time"
)

type RoleDto struct {
	Id          uint                `json:"id"`
	Name        string              `json:"name"`
	Code        string              `json:"code"`
	Description string              `json:"description"`
	ResourceID  uint                `json:"resource_id"`
	Permissions []*model.Permission `json:"permissions"`
	IsDeleted   *bool               `json:"is_deleted"`
	IsActive    bool                `json:"is_active"`
	CountUser   int64               `json:"count_user"`
	Users       []*UserDto          `json:"users"`
	CreatedAt   time.Time           `json:"created_at"`
}
