package model

import (
	"time"
)

type Role struct {
	VModel
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Permissions []*Permission `gorm:"many2many:role_permission" json:"permissions"`
	IsDeleted   bool          `json:"is_deleted"`
	CreatedAt   time.Time     `json:"created_at"`
	Users       []*User       `gorm:"many2many:user_role" json:"users"`
}
