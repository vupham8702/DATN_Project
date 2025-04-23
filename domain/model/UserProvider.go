package model

type UserProvider struct {
	VModel
	Email            string `json:"email"`
	AccountInfo      []byte `gorm:"type:jsonb;default:'{}'"`
	UserID           uint   `json:"gorm:not null"`
	User             User   `gorm:"constraint:OnDelete:CASCADE;"`
	Provider         string `gorm:"not null"`
	AppleId          string `json:"apple_id"`
	UserType         string `json:"user_type"`
	ProviderIdentify string `json:"provider_identify"`
}
