package model

// CVTemplate đại diện cho mẫu CV trong hệ thống
type CVTemplate struct {
	VModel
	Name          string `json:"name" gorm:"not null"`
	Description   string `json:"description"`
	FilePath      string `json:"file_path" gorm:"not null"`
	ThumbnailPath string `json:"thumbnail_path"`
	FileType      string `json:"file_type" gorm:"not null"` // docx, xlsx
	Category      string `json:"category"`
	Tags          string `json:"tags"`
	DownloadCount int    `json:"download_count" gorm:"default:0"`
}

func (CVTemplate) TableName() string {
	return "cv_template"
}
