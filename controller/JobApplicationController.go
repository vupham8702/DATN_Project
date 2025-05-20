package controller

import (
	"datn_backend/message"
	"datn_backend/payload/response"
	"datn_backend/service"
	"datn_backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetMyApplications godoc
// @Summary Lấy danh sách ứng tuyển của tôi
// @Description Lấy danh sách các công việc mà người dùng đã ứng tuyển
// @Tags JobApplication
// @Accept json
// @Produce json
// @Success 200 {object} response.VResponse{data=[]model.JobApplication}
// @Failure 401 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/post-job/my-applications [get]
// @Security BearerAuth
func GetMyApplications(c *gin.Context) {
	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Gọi service
	applications, msg := service.GetMyApplications(*uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về danh sách
	response.Response(c, applications, message.Success)
}

// GetJobApplications godoc
// @Summary Lấy danh sách ứng viên cho bài đăng
// @Description Lấy danh sách các ứng viên đã ứng tuyển vào bài đăng của nhà tuyển dụng
// @Tags JobApplication
// @Accept json
// @Produce json
// @Param id path int true "Post Job ID"
// @Success 200 {object} response.VResponse{data=[]model.JobApplication}
// @Failure 401 {object} response.VResponse
// @Failure 403 {object} response.VResponse
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/post-job/{id}/applications [get]
// @Security BearerAuth
func GetJobApplications(c *gin.Context) {
	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Lấy ID bài đăng từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Gọi service
	applications, msg := service.GetJobApplications(uint(id), *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về danh sách
	response.Response(c, applications, message.Success)
}

// UpdateApplicationStatus godoc
// @Summary Cập nhật trạng thái ứng tuyển
// @Description Nhà tuyển dụng cập nhật trạng thái ứng tuyển của ứng viên
// @Tags JobApplication
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Param status body string true "Status (pending, reviewing, accepted, rejected)"
// @Success 200 {object} response.VResponse
// @Failure 400 {object} response.VResponse
// @Failure 401 {object} response.VResponse
// @Failure 403 {object} response.VResponse
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/post-job/applications/{id}/status [put]
// @Security BearerAuth
func UpdateApplicationStatus(c *gin.Context) {
	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Lấy ID ứng tuyển từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Bind JSON vào struct
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Response(c, err)
		return
	}

	// Kiểm tra trạng thái hợp lệ
	validStatuses := []string{"pending", "reviewing", "accepted", "rejected"}
	isValid := false
	for _, status := range validStatuses {
		if req.Status == status {
			isValid = true
			break
		}
	}

	if !isValid {
		response.Response(c, nil, message.Message{
			Code:    http.StatusBadRequest,
			Message: "Trạng thái không hợp lệ. Các trạng thái hợp lệ: pending, reviewing, accepted, rejected",
		})
		return
	}

	// Gọi service
	msg := service.UpdateApplicationStatus(uint(id), req.Status, *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về thành công
	response.Response(c, nil, message.Success)
}
