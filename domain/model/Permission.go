package model

type Permission struct {
	VModel
	Name        string `json:"name"`
	Description string `json:"description"`
	ResourceID  uint   `json:"resource_id"`
	//Resource    *Resource    `json:"resource"`
	ParentID  *uint        `json:"parent_id"`
	Children  []Permission `gorm:"foreignKey:ParentID;" json:"children"`
	IsChecked bool         `json:"is_checked" gorm:"-"`
	Priority  int          `gorm:"column:priority" json:"priority"`
	Roles     []Role       `json:"roles" gorm:"many2many:role_permission;"`
}
