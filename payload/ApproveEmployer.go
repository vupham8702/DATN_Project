package payload

// ApproveEmployer is used by admins to approve employer accounts
type ApproveEmployer struct {
	UserID uint   `json:"userId" binding:"required"`
	Status string `json:"status" binding:"required"` // "approved" or "rejected"
	Note   string `json:"note"`
}
