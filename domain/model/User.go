package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type User struct {
	VModel
	Age      int    `json:"age"`
	Username string `json:"username"`
	Email    string `json:"email"`
	//Phone     string `json:"phone"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsSupper  bool   `json:"is_supper"`
	IsActive  bool   `json:"is_active"`
	Avatar    string `json:"avatar"`
	//ResourceID   *uint           `json:"resource_id"`
	//Resource     *Resource       `json:"resource"`
	Roles     []*Role         `gorm:"many2many:user_role" json:"roles"`
	Providers []*UserProvider `gorm:"foreignKey:UserID"`
	IsLocked  bool            `json:"is_locked"`
	//MemberShip   Membership      `json:"membership" gorm:"foreignKey:MembershipId"`
	//MembershipId uint            `json:"membership_id"`
	//Gender       int        `json:"gender"`
	//DateOfBirth  *time.Time `json:"date_of_birth"`
}

func (v *User) BeforeCreate(tx *gorm.DB) (err error) {
	atTime := time.Now()
	//v.Gender = 0
	//v.MembershipId = 1
	v.CreatedAt = atTime
	v.UpdatedAt = atTime
	v.IsDeleted = false
	v.IsActive = true
	return
}

func (v *User) BeforeUpdate(tx *gorm.DB) (err error) {
	v.UpdatedAt = time.Now()
	return
}

func (v *User) BeforeDelete(tx *gorm.DB) (err error) {
	v.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	v.IsDeleted = true
	return
}
