package model

type JobApplication struct {
	VModel
	PostJobID   uint   `json:"post_job_id" binding:"required"`
	Status      string `json:"status"`
	CoverLetter string `json:"cover_letter"`
	ResumeURL   string `json:"resume " binding:"required"`
	Notes       string `json:"notes"`
}

func (JobApplication) TableName() string {
	return "post_job_application"
}
