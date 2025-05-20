package service

import (
	"archive/zip"
	"bytes"
	"datn_backend/domain/model"
	"datn_backend/domain/repository"
	"datn_backend/message"
	"datn_backend/middleware"
	"fmt"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetAllCVTemplates lấy tất cả mẫu CV
func GetAllCVTemplates() (interface{}, interface{}) {
	templates, err := repository.GetAllCVTemplates()
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get CV templates: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return templates, nil
}

// GetCVTemplateByID lấy mẫu CV theo ID
func GetCVTemplateByID(id uint) (interface{}, interface{}) {
	template, err := repository.GetCVTemplateByID(id)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get CV template: %v", err))
		return nil, message.CVTemplateNotFound
	}

	return template, nil
}

// CreateCVTemplate tạo mẫu CV mới
func CreateCVTemplate(template *model.CVTemplate, userID uint) (interface{}, interface{}) {
	// Thiết lập các giá trị mặc định
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	template.CreatedBy = userID
	template.UpdatedBy = userID
	template.IsDeleted = false

	if err := repository.CreateCVTemplate(template); err != nil {
		middleware.Log(fmt.Errorf("Failed to create CV template: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return template, nil
}

// UpdateCVTemplate cập nhật mẫu CV
func UpdateCVTemplate(template *model.CVTemplate, userID uint) (interface{}, interface{}) {
	// Kiểm tra xem mẫu CV có tồn tại không
	existingTemplate, err := repository.GetCVTemplateByID(template.ID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get CV template: %v", err))
		return nil, message.CVTemplateNotFound
	}

	// Cập nhật các trường
	template.CreatedAt = existingTemplate.CreatedAt
	template.CreatedBy = existingTemplate.CreatedBy
	template.UpdatedAt = time.Now()
	template.UpdatedBy = userID
	template.IsDeleted = existingTemplate.IsDeleted

	if err := repository.UpdateCVTemplate(template); err != nil {
		middleware.Log(fmt.Errorf("Failed to update CV template: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return template, nil
}

// DeleteCVTemplate xóa mềm mẫu CV
func DeleteCVTemplate(id uint, userID uint) (interface{}, interface{}) {
	// Kiểm tra xem mẫu CV có tồn tại không
	_, err := repository.GetCVTemplateByID(id)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get CV template: %v", err))
		return nil, message.CVTemplateNotFound
	}

	if err := repository.DeleteCVTemplate(id); err != nil {
		middleware.Log(fmt.Errorf("Failed to delete CV template: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return nil, nil
}

// DownloadCVTemplate tải mẫu CV
func DownloadCVTemplate(id uint) (string, interface{}) {
	template, err := repository.GetCVTemplateByID(id)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get CV template: %v", err))
		return "", message.CVTemplateNotFound
	}

	return template.FilePath, nil
}

// GenerateCV tạo CV từ mẫu và thông tin người dùng
func GenerateCV(templateID uint, userID uint) (string, interface{}) {
	// Lấy thông tin mẫu CV
	template, err := repository.GetCVTemplateByID(templateID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get CV template: %v", err))
		return "", message.CVTemplateNotFound
	}

	// Lấy thông tin người dùng
	user, err := repository.GetUserById(userID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get user: %v", err))
		return "", message.UserNotFound
	}

	// Lấy thông tin profile của người dùng
	profile, err := repository.GetJobseekerProfileByUserID(&userID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get jobseeker profile: %v", err))
		// Không return lỗi ở đây, vẫn tiếp tục với profile nil
	}

	// Tạo tên file mới
	fileExt := filepath.Ext(template.FilePath)
	fileName := fmt.Sprintf("%s-%s%s", user.Username, uuid.New().String()[:8], fileExt)
	outputPath := filepath.Join("./uploads/cv", fileName)

	// Đảm bảo thư mục tồn tại
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	// Đọc file mẫu
	templatePath := "." + template.FilePath
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		middleware.Log(fmt.Errorf("Template file not found: %v", err))
		return "", message.Message{Message: "Template file not found", Code: 404}
	}

	// Xử lý file docx
	if strings.HasSuffix(templatePath, ".docx") {
		repls := map[string]string{
			"{{name}}":  user.Username,
			"{{email}}": user.Email,
			// … các placeholder khác
		}
		if err := processDocxTemplate(templatePath, outputPath, user, profile, repls); err != nil {
			middleware.Log(fmt.Errorf("Failed to process docx template: %v", err))
			return "", message.Message{Message: "Failed to generate CV", Code: 500}
		}
	} else if strings.HasSuffix(templatePath, ".xlsx") {
		if err := processExcelTemplate(templatePath, outputPath, user, profile); err != nil {
			middleware.Log(fmt.Errorf("Failed to process excel template: %v", err))
			return "", message.Message{Message: "Failed to generate CV", Code: 500}
		}
	} else {
		middleware.Log(fmt.Errorf("Unsupported file format: %s", templatePath))
		return "", message.Message{Message: "Unsupported file format", Code: 400}
	}

	// Tăng số lượt tải xuống
	if err := repository.IncrementDownloadCount(templateID); err != nil {
		middleware.Log(fmt.Errorf("Failed to increment download count: %v", err))
		// Không return lỗi ở đây, vẫn trả về file đã tạo
	}

	return "/uploads/cv/" + fileName, nil
}

// processDocxTemplate xử lý file docx template
func processDocxTemplate(templatePath, outputPath string, user *model.User, profile *model.JobseekerProfile, repls map[string]string) error {
	r, err := zip.OpenReader(templatePath)
	if err != nil {
		return err
	}
	defer r.Close()

	// 2. Tạo zip mới
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	zw := zip.NewWriter(outFile)
	defer zw.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		buf, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}

		// Khi gặp document.xml thì replace thô
		if f.Name == "word/document.xml" {
			content := string(buf)
			for placeholder, val := range repls {
				content = strings.ReplaceAll(content, placeholder, val)
			}
			buf = []byte(content)
		}

		// Ghi vào zip mới
		w, err := zw.Create(f.Name)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, bytes.NewReader(buf)); err != nil {
			return err
		}
	}
	return nil
}

// processExcelTemplate xử lý file excel template
func processExcelTemplate(templatePath, outputPath string, user *model.User, profile *model.JobseekerProfile) error {
	// Mở file excel
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			middleware.Log(fmt.Errorf("Failed to close excel file: %v", err))
		}
	}()

	// Lấy tất cả các sheet
	sheets := f.GetSheetList()

	// Duyệt qua từng sheet
	for _, sheet := range sheets {
		// Lấy tất cả các cell có giá trị
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}

		// Duyệt qua từng hàng
		for i, row := range rows {
			// Duyệt qua từng cột
			for j, cell := range row {
				// Thay thế thông tin từ user
				if user != nil {
					// Thay thế họ tên
					if strings.Contains(cell, "{{name}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{name}}", user.Username))
					}
					if strings.Contains(cell, "{{fullname}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{fullname}}", user.Username))
					}

					// Thay thế email
					if strings.Contains(cell, "{{email}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{email}}", user.Email))
					}
				}

				// Thay thế số điện thoại (chỉ lấy từ profile)
				if strings.Contains(cell, "{{phone}}") {
					if profile != nil && profile.PhoneNumber != "" {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{phone}}", profile.PhoneNumber))
					} else {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{phone}}", ""))
					}
				}

				// Thay thế thông tin từ profile
				if profile != nil {
					// Thông tin cá nhân
					if strings.Contains(cell, "{{address}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{address}}", profile.Address))
					}

					// Ngày sinh
					if strings.Contains(cell, "{{dob}}") || strings.Contains(cell, "{{date_of_birth}}") {
						if profile.DateOfBirth != nil {
							dobStr := profile.DateOfBirth.Format("02/01/2006")
							if strings.Contains(cell, "{{dob}}") {
								f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{dob}}", dobStr))
							}
							if strings.Contains(cell, "{{date_of_birth}}") {
								f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{date_of_birth}}", dobStr))
							}
						} else {
							if strings.Contains(cell, "{{dob}}") {
								f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{dob}}", ""))
							}
							if strings.Contains(cell, "{{date_of_birth}}") {
								f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{date_of_birth}}", ""))
							}
						}
					}

					// Giới tính
					if strings.Contains(cell, "{{gender}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{gender}}", profile.Gender))
					}

					// Thành phố
					if strings.Contains(cell, "{{city}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{city}}", profile.City))
					}

					// Quốc gia
					if strings.Contains(cell, "{{country}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{country}}", profile.Country))
					}

					// Tiêu đề hồ sơ
					if strings.Contains(cell, "{{profile_title}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{profile_title}}", profile.ProfileTitle))
					}

					// Giới thiệu bản thân
					if strings.Contains(cell, "{{about}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{about}}", profile.About))
					}

					// Kỹ năng
					if strings.Contains(cell, "{{skills}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{skills}}", profile.Skills))
					}

					// Học vấn
					if strings.Contains(cell, "{{education}}") {
						if profile.Education.Valid {
							f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{education}}", profile.Education.String))
						} else {
							f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{education}}", ""))
						}
					}

					// Kinh nghiệm làm việc
					if strings.Contains(cell, "{{experience}}") {
						if profile.Experience.Valid {
							f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{experience}}", profile.Experience.String))
						} else {
							f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{experience}}", ""))
						}
					}

					// Chứng chỉ
					if strings.Contains(cell, "{{certifications}}") {
						if profile.Certifications.Valid {
							f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{certifications}}", profile.Certifications.String))
						} else {
							f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{certifications}}", ""))
						}
					}

					// Ngôn ngữ
					if strings.Contains(cell, "{{languages}}") {
						if profile.Languages.Valid {
							f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{languages}}", profile.Languages.String))
						} else {
							f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{languages}}", ""))
						}
					}

					// Tình trạng sẵn sàng làm việc
					if strings.Contains(cell, "{{availability}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{availability}}", profile.Availability))
					}

					// Liên kết mạng xã hội
					if strings.Contains(cell, "{{linkedin}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{linkedin}}", profile.LinkedinProfile))
					}
					if strings.Contains(cell, "{{github}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{github}}", profile.GithubProfile))
					}
					if strings.Contains(cell, "{{website}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{website}}", profile.WebsiteURL))
					}

					// Họ và tên riêng
					if strings.Contains(cell, "{{first_name}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{first_name}}", profile.FistName))
					}
					if strings.Contains(cell, "{{last_name}}") {
						f.SetCellValue(sheet, getCellName(j, i), strings.ReplaceAll(cell, "{{last_name}}", profile.LastName))
					}
				}
			}
		}
	}

	// Lưu file mới
	return f.SaveAs(outputPath)
}

// getCellName chuyển đổi vị trí cột, hàng thành tên cell (A1, B2, ...)
func getCellName(col, row int) string {
	colName := ""
	for col >= 0 {
		colName = string(rune('A'+col%26)) + colName
		col = col/26 - 1
	}
	return fmt.Sprintf("%s%d", colName, row+1)
}

// UploadCVTemplate tải lên mẫu CV mới
func UploadCVTemplate(file io.Reader, fileName string, fileType string, userID uint) (interface{}, interface{}) {
	// Tạo tên file mới
	fileExt := filepath.Ext(fileName)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	filePath := filepath.Join("./uploads/templates", newFileName)

	// Đảm bảo thư mục tồn tại
	os.MkdirAll(filepath.Dir(filePath), 0755)

	// Tạo file mới
	dst, err := os.Create(filePath)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to create file: %v", err))
		return nil, message.Message{Message: "Failed to create file", Code: 500}
	}
	defer dst.Close()

	// Sao chép nội dung từ file tải lên vào file mới
	if _, err = io.Copy(dst, file); err != nil {
		middleware.Log(fmt.Errorf("Failed to copy file: %v", err))
		return nil, message.Message{Message: "Failed to copy file", Code: 500}
	}

	// Tạo bản ghi mẫu CV mới
	template := &model.CVTemplate{
		Name:     strings.TrimSuffix(fileName, fileExt),
		FilePath: "/uploads/templates/" + newFileName,
		FileType: fileType,
	}

	// Lưu vào database
	result, msg := CreateCVTemplate(template, userID)
	if msg != nil {
		// Xóa file nếu lưu database thất bại
		os.Remove(filePath)
		return nil, msg
	}

	return result, nil
}

// UploadCVThumbnail tải lên thumbnail cho mẫu CV
func UploadCVThumbnail(file io.Reader, fileName string, templateID uint, userID uint) (interface{}, interface{}) {
	// Lấy thông tin mẫu CV
	template, err := repository.GetCVTemplateByID(templateID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get CV template: %v", err))
		return nil, message.CVTemplateNotFound
	}

	// Tạo tên file mới
	fileExt := filepath.Ext(fileName)
	// Chỉ chấp nhận các định dạng hình ảnh phổ biến
	if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" && fileExt != ".gif" {
		return nil, message.Message{Message: "Only image files are supported for thumbnail (.jpg, .jpeg, .png, .gif)", Code: 400}
	}

	newFileName := fmt.Sprintf("thumb_%d_%s%s", templateID, uuid.New().String()[:8], fileExt)
	filePath := filepath.Join("./uploads/templates/thumbnails", newFileName)

	// Đảm bảo thư mục tồn tại
	os.MkdirAll(filepath.Dir(filePath), 0755)

	// Tạo file mới
	dst, err := os.Create(filePath)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to create thumbnail file: %v", err))
		return nil, message.Message{Message: "Failed to create thumbnail file", Code: 500}
	}
	defer dst.Close()

	// Sao chép nội dung từ file tải lên vào file mới
	if _, err = io.Copy(dst, file); err != nil {
		middleware.Log(fmt.Errorf("Failed to copy thumbnail file: %v", err))
		return nil, message.Message{Message: "Failed to copy thumbnail file", Code: 500}
	}

	// Cập nhật đường dẫn thumbnail
	template.ThumbnailPath = "/uploads/templates/thumbnails/" + newFileName

	// Lưu lại thông tin
	updatedTemplate, msg := UpdateCVTemplate(template, userID)
	if msg != nil {
		// Xóa file nếu cập nhật database thất bại
		os.Remove(filePath)
		return nil, msg
	}

	return updatedTemplate, nil
}
