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
	IsApproved       bool   `json:"is_approved"`   // For employer accounts
	ApprovedBy       uint   `json:"approved_by"`   // ID of admin who approved
	ApprovalNote     string `json:"approval_note"` // Note from admin
	ReceivedNoti     bool   `json:"received_noti"`
}
