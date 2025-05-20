package controller

import (
	"datn_backend/domain/model"
	"datn_backend/message"
	"datn_backend/middleware"
	"datn_backend/payload/response"
	"datn_backend/service"
	"datn_backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strconv"
	"strings"
)

// GetAllCVTemplates godoc
// @Summary Lấy tất cả mẫu CV
// @Description Lấy danh sách tất cả các mẫu CV
// @Tags CV
// @Accept json
// @Produce json
// @Success 200 {object} response.VResponse{data=[]model.CVTemplate}
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/cv/templates [get]
func GetAllCVTemplates(c *gin.Context) {
	templates, msg := service.GetAllCVTemplates()
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	response.Response(c, templates, message.Success)
}

// GetCVTemplateByID godoc
// @Summary Lấy mẫu CV theo ID
// @Description Lấy thông tin chi tiết của mẫu CV theo ID
// @Tags CV
// @Accept json
// @Produce json
// @Param id path int true "Template ID"
// @Success 200 {object} response.VResponse{data=model.CVTemplate}
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/cv/templates/{id} [get]
func GetCVTemplateByID(c *gin.Context) {
	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	template, msg := service.GetCVTemplateByID(uint(id))
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	response.Response(c, template, message.Success)
}

// GetCVPreview godoc
// @Summary Xem trước mẫu CV
// @Description Xem trước hình ảnh thumbnail của mẫu CV
// @Tags CV
// @Accept json
// @Produce image/jpeg
// @Param id path int true "Template ID"
// @Success 200 {file} file "CV Template thumbnail"
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/cv/templates/{id}/preview [get]
func GetCVPreview(c *gin.Context) {
	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	template, msg := service.GetCVTemplateByID(uint(id))
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Kiểm tra xem có thumbnail không
	cvTemplate, ok := template.(*model.CVTemplate)
	if !ok {
		response.Response(c, nil, message.InternalServerError)
		return
	}

	if cvTemplate.ThumbnailPath == "" {
		response.Response(c, nil, message.Message{Message: "Thumbnail not found", Code: 404})
		return
	}

	// Trả về file thumbnail
	c.File("." + cvTemplate.ThumbnailPath)
}

// DownloadCVTemplate godoc
// @Summary Tải xuống mẫu CV gốc
// @Description Tải xuống file mẫu CV gốc (không điền thông tin)
// @Tags Admin
// @Accept json
// @Produce octet-stream
// @Param id path int true "Template ID"
// @Success 200 {file} file "CV Template file"
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/cv/admin/templates/{id}/download-original [get]
// @Security BearerAuth
func DownloadCVTemplate(c *gin.Context) {
	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Lấy đường dẫn file
	filePath, msg := service.DownloadCVTemplate(uint(id))
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về file
	c.File("." + filePath)
}

// DownloadAndFillCVTemplate godoc
// @Summary Tải xuống mẫu CV với thông tin cá nhân
// @Description Tải xuống file mẫu CV với thông tin cá nhân được điền vào
// @Tags CV
// @Accept json
// @Produce octet-stream
// @Param id path int true "Template ID"
// @Success 200 {file} file "CV Template file with personal info"
// @Failure 401 {object} response.VResponse
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/cv/templates/{id}/download [get]
// @Security BearerAuth
func DownloadAndFillCVTemplate(c *gin.Context) {
	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Tạo CV từ mẫu và thông tin người dùng
	filePath, msg := service.GenerateCV(uint(id), *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về file
	c.File("." + filePath)
}

// CreateCVTemplate godoc
// @Summary Tạo mẫu CV mới
// @Description Tải lên file mẫu CV mới
// @Tags Admin
// @Accept multipart/form-data
// @Produce json
// @Param template_file formData file true "CV Template file (.docx, .xlsx)"
// @Param name formData string true "Template name"
// @Param description formData string false "Template description"
// @Param category formData string false "Template category"
// @Param tags formData string false "Template tags"
// @Success 200 {object} response.VResponse{data=model.CVTemplate}
// @Failure 400 {object} response.VResponse
// @Failure 401 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/cv/admin/templates [post]
// @Security BearerAuth
func CreateCVTemplate(c *gin.Context) {
	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Lấy file template từ form
	templateFile, templateHeader, err := c.Request.FormFile("template_file")
	if err != nil {
		response.Response(c, nil, message.Message{Message: "Template file is required", Code: 400})
		return
	}
	defer templateFile.Close()

	// Kiểm tra định dạng file template
	fileExt := filepath.Ext(templateHeader.Filename)
	if fileExt != ".docx" && fileExt != ".xlsx" {
		response.Response(c, nil, message.Message{Message: "Only .docx and .xlsx files are supported for template", Code: 400})
		return
	}

	// Lấy các thông tin khác từ form
	name := c.PostForm("name")
	if name == "" {
		name = strings.TrimSuffix(templateHeader.Filename, fileExt)
	}

	description := c.PostForm("description")
	category := c.PostForm("category")
	tags := c.PostForm("tags")

	// Xác định loại file
	fileType := "docx"
	if fileExt == ".xlsx" {
		fileType = "xlsx"
	}

	// Tải lên file template
	template, msg := service.UploadCVTemplate(templateFile, templateHeader.Filename, fileType, *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Cập nhật thông tin bổ sung
	cvTemplate, ok := template.(*model.CVTemplate)
	if !ok {
		response.Response(c, nil, message.InternalServerError)
		return
	}

	cvTemplate.Name = name
	cvTemplate.Description = description
	cvTemplate.Category = category
	cvTemplate.Tags = tags

	// Lưu lại thông tin
	updatedTemplate, msg := service.UpdateCVTemplate(cvTemplate, *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Kiểm tra xem có file thumbnail không
	thumbnailFile, thumbnailHeader, err := c.Request.FormFile("thumbnail_file")
	if err == nil {
		// Có file thumbnail, tiến hành tải lên
		defer thumbnailFile.Close()

		// Tải lên thumbnail
		finalTemplate, msg := service.UploadCVThumbnail(thumbnailFile, thumbnailHeader.Filename, cvTemplate.ID, *uid)
		if msg != nil {
			// Không return lỗi ở đây, vẫn trả về template đã tạo
			middleware.Log(fmt.Errorf("Failed to upload thumbnail: %v", msg))
		} else {
			// Nếu tải thumbnail thành công, trả về template với thumbnail
			updatedTemplate = finalTemplate
		}
	}

	response.Response(c, updatedTemplate, message.Success)
}

// UpdateCVTemplate godoc
// @Summary Cập nhật mẫu CV
// @Description Cập nhật thông tin mẫu CV
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path int true "Template ID"
// @Param template body model.CVTemplate true "Template info"
// @Success 200 {object} response.VResponse{data=model.CVTemplate}
// @Failure 400 {object} response.VResponse
// @Failure 401 {object} response.VResponse
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/cv/admin/templates/{id} [put]
// @Security BearerAuth
func UpdateCVTemplate(c *gin.Context) {
	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Bind JSON vào struct
	var template model.CVTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		response.Response(c, err)
		return
	}

	// Đảm bảo ID khớp với path parameter
	template.ID = uint(id)

	// Cập nhật mẫu CV
	updatedTemplate, msg := service.UpdateCVTemplate(&template, *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	response.Response(c, updatedTemplate, message.Success)
}

// DeleteCVTemplate godoc
// @Summary Xóa mẫu CV
// @Description Xóa mẫu CV
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path int true "Template ID"
// @Success 200 {object} response.VResponse
// @Failure 401 {object} response.VResponse
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/cv/admin/templates/{id} [delete]
// @Security BearerAuth
func DeleteCVTemplate(c *gin.Context) {
	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Xóa mẫu CV
	_, msg := service.DeleteCVTemplate(uint(id), *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	response.Response(c, nil, message.Success)
}

// UploadCVThumbnail godoc
// @Summary Tải lên thumbnail cho mẫu CV
// @Description Tải lên hoặc cập nhật thumbnail cho mẫu CV
// @Tags Admin
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Template ID"
// @Param thumbnail_file formData file true "Thumbnail file (.jpg, .jpeg, .png, .gif)"
// @Success 200 {object} response.VResponse{data=model.CVTemplate}
// @Failure 400 {object} response.VResponse
// @Failure 401 {object} response.VResponse
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/cv/admin/templates/{id}/thumbnail [post]
// @Security BearerAuth
func UploadCVThumbnail(c *gin.Context) {
	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Lấy file thumbnail từ form
	thumbnailFile, thumbnailHeader, err := c.Request.FormFile("thumbnail_file")
	if err != nil {
		response.Response(c, nil, message.Message{Message: "Thumbnail file is required", Code: 400})
		return
	}
	defer thumbnailFile.Close()

	// Tải lên thumbnail
	updatedTemplate, msg := service.UploadCVThumbnail(thumbnailFile, thumbnailHeader.Filename, uint(id), *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	response.Response(c, updatedTemplate, message.Success)
}
