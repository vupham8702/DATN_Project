package model

import (
	"gorm.io/gorm"
	"time"
)

type Role struct {
	VModel
	Name        string `json:"name"`
	Description string `json:"description"`
	ResourceID  uint   `json:"resource_id"`
	//Resource    *Resource     `json:"resource"`
	Permissions []*Permission `gorm:"many2many:role_permission" json:"permissions"`
	IsDeleted   bool          `json:"is_deleted"`
	CreatedAt   time.Time     `json:"created_at"`
	IsActive    bool          `json:"is_active"`
	Users       []*User       `gorm:"many2many:user_role" json:"users"`
}

func (v *Role) BeforeCreateRole(tx *gorm.DB) (err error) {
	v.IsActive = true
	return
}
