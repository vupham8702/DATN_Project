package controller

import (
	repo "datn_backend/domain/repository"
	"datn_backend/message"
	"datn_backend/middleware"
	"datn_backend/payload/request"
	"datn_backend/payload/response"
	"datn_backend/service"
	"datn_backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetUserProfile godoc
// @Summary Lấy thông tin hồ sơ người dùng hiện tại
// @Description Lấy thông tin hồ sơ của người dùng đang đăng nhập
// @Tags ProfileController
// @Accept json
// @Produce json
// @Success 200 {object} response.VResponse
// @Router /datn_backend/profile [get]
// @Security BearerAuth
func GetUserProfile(c *gin.Context) {
	// Lấy thông tin người dùng từ context
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	result, err := service.GetProfile(uid)
	if err != nil {
		response.Response(c, err)
		return
	}

	response.Response(c, result, message.Success)
}

// UpdateJobseekerProfile godoc
// @Summary Cập nhật hồ sơ ứng viên
// @Description Cập nhật thông tin hồ sơ của ứng viên
// @Tags ProfileController
// @Accept json
// @Produce json
// @Param profile body request.JobseekerProfileRequest true "Thông tin hồ sơ ứng viên"
// @Success 200 {object} response.VResponse
// @Router /datn_backend/profile/jobseeker [put]
// @Security BearerAuth
func UpdateJobseekerProfile(c *gin.Context) {
	// Lấy thông tin người dùng từ context
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	var profileRequest request.JobseekerProfileRequest
	if err := c.ShouldBindJSON(&profileRequest); err != nil {
		middleware.Log(err)
		response.Response(c, message.ValidationError, http.StatusBadRequest)
		return
	}

	result, err := service.UpdateJobseekerProfile(*uid, &profileRequest)
	if err != nil {
		response.Response(c, err)
		return
	}

	response.Response(c, result, message.Success)
}

// UpdateEmployerProfile godoc
// @Summary Cập nhật hồ sơ nhà tuyển dụng
// @Description Cập nhật thông tin hồ sơ của nhà tuyển dụng
// @Tags ProfileController
// @Accept json
// @Produce json
// @Param profile body request.EmployerProfileRequest true "Thông tin hồ sơ nhà tuyển dụng"
// @Success 200 {object} response.VResponse
// @Router /datn_backend/profile/employer [put]
// @Security BearerAuth
func UpdateEmployerProfile(c *gin.Context) {
	// Lấy thông tin người dùng từ context
	//isSupper, errSuper := utils.GetFieldInToken(c, "issupper")
	//if errSuper != nil {
	//	response.Response(c, message.InternalServerError)
	//	return
	//}
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	var profileRequest request.EmployerProfileRequest
	if err := c.ShouldBindJSON(&profileRequest); err != nil {
		middleware.Log(err)
		response.Response(c, message.ValidationError, http.StatusBadRequest)
		return
	}

	result, err := service.UpdateEmployerProfile(*uid, &profileRequest)
	if err != nil {
		middleware.Log(err)
		response.Response(c, err)
		return
	}
	response.Response(c, message.Success, result)
	return
}

// UploadProfilePhoto godoc
// @Summary Upload ảnh cho hồ sơ
// @Description Upload ảnh đại diện, CV, logo công ty hoặc banner công ty
// @Tags ProfileController
// @Accept multipart/form-data
// @Produce json
// @Param photoType formData string true "Loại ảnh (profile, resume, company_logo, company_banner)"
// @Param file formData file true "File ảnh cần upload"
// @Success 200 {object} response.VResponse
// @Router /datn_backend/profile/upload-photo [post]
// @Security BearerAuth
func UploadProfilePhoto(c *gin.Context) {
	// Lấy thông tin người dùng từ context
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Lấy loại ảnh từ form
	photoType := c.PostForm("photoType")
	if photoType == "" {
		response.Response(c, message.Message{Message: "Photo type is required", Code: http.StatusBadRequest})
		return
	}

	// Kiểm tra photoType hợp lệ
	validPhotoTypes := map[string]bool{
		"profile":        true,
		"resume":         true,
		"company_logo":   true,
		"company_banner": true,
	}

	if !validPhotoTypes[photoType] {
		response.Response(c, message.Message{Message: "Invalid photo type", Code: http.StatusBadRequest})
		return
	}

	// Lấy file từ form
	file, err := c.FormFile("file")
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get file: %v", err))
		response.Response(c, message.Message{Message: "File is required", Code: http.StatusBadRequest})
		return
	}

	result, apiErr := service.UploadProfilePhoto(c, *uid, photoType, file)
	if apiErr != nil {
		response.Response(c, apiErr)
		return
	}

	response.Response(c, result, message.Success)
}

// GetJobseekerPublicProfile godoc
// @Summary Lấy thông tin hồ sơ công khai của ứng viên
// @Description Lấy thông tin hồ sơ công khai của ứng viên theo ID
// @Tags ProfileController
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.VResponse
// @Router /datn_backend/profile/jobseeker/{id} [get]
func GetJobseekerPublicProfile(c *gin.Context) {
	// Lấy ID từ path parameter
	userIDStr := c.Param("id")
	if userIDStr == "" {
		response.Response(c, message.ValidationError, http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.Response(c, message.Message{Message: "Invalid user ID", Code: http.StatusBadRequest})
		return
	}

	uid := uint(userID)
	result, apiErr := repo.GetJobseekerProfileByUserID(&uid)
	if apiErr != nil {
		response.Response(c, apiErr)
		return
	}

	response.Response(c, result, message.Success)
}

// GetEmployerPublicProfile godoc
// @Summary Lấy thông tin hồ sơ công khai của nhà tuyển dụng
// @Description Lấy thông tin hồ sơ công khai của nhà tuyển dụng theo ID
// @Tags ProfileController
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.VResponse
// @Router /datn_backend/profile/employer/{id} [get]
func GetEmployerPublicProfile(c *gin.Context) {
	// Lấy ID từ path parameter
	userIDStr := c.Param("id")
	if userIDStr == "" {
		response.Response(c, message.ValidationError, http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.Response(c, message.Message{Message: "Invalid user ID", Code: http.StatusBadRequest})
		return
	}

	uid := uint(userID)
	result, apiErr := repo.GetEmployerProfileByUserID(&uid)
	if apiErr != nil {
		response.Response(c, apiErr)
		return
	}

	response.Response(c, result, message.Success)
}
