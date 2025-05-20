package model

type PostJob struct {
	VModel
	Title             string `json:"title" binding:"required"`
	Company           string `json:"company"`
	Logo              string `json:"logo"`
	Location          string `json:"location"`
	Salary            string `json:"salary"`
	Status            string `json:"status"`
	Type              string `json:"type"`
	TimeFrame         string `json:"time_frame"` // ví dụ "01/05/2025 - 01/06/2025"
	Experience        string `json:"experience"`
	Gender            string `json:"gender"`
	Description       string `json:"description" binding:"required"`
	ApplicationsCount int    `json:"applications_count"`
	Requirements      string `json:"requirements"`
}
