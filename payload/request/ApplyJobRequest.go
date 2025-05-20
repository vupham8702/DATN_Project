package request

import "mime/multipart"

type ApplyJobRequest struct {
	ResumeFile  *multipart.FileHeader `json:"resume_file" binding:"required,url" form:"resume_url"`
	CoverLetter string                `json:"cover_letter" form:"cover_letter"`
}
