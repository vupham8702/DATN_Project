package controller

import (
	"datn_backend/domain/model"
	"datn_backend/message"
	"datn_backend/payload/response"
	"datn_backend/service"
	"datn_backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// CreatePostJob godoc
// @Summary Tạo mới bài đăng tuyển dụng
// @Description Đăng bài tuyển dụng (chờ duyệt bởi Admin)
// @Tags JobPost
// @Accept json
// @Produce json
// @Param payload body object true "Job Post Payload với start_date và end_date (YYYY-MM-DD)"
// @Success 201 {object} response.VResponse{data=model.PostJob}
// @Failure 400 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/post-job/create [post]
// @Security BearerAuth
func CreatePostJob(c *gin.Context) {
	// Tạo struct để nhận dữ liệu từ frontend
	type PostJobRequest struct {
		model.PostJob
		StartDate string `json:"start_date"` // Ngày bắt đầu từ date-time-picker
		EndDate   string `json:"end_date"`   // Ngày kết thúc từ date-time-picker
	}

	// 1. Bind JSON vào struct
	var req PostJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Response(c, err)
		return
	}

	// 2. Xử lý trường time_frame
	if req.StartDate != "" && req.EndDate != "" {
		// Parse ngày bắt đầu
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			response.Response(c, nil, message.Message{Message: "Invalid start date format. Expected YYYY-MM-DD", Code: 400})
			return
		}

		// Parse ngày kết thúc
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			response.Response(c, nil, message.Message{Message: "Invalid end date format. Expected YYYY-MM-DD", Code: 400})
			return
		}

		// Định dạng lại theo dd/mm/yyyy - dd/mm/yyyy
		req.TimeFrame = fmt.Sprintf("%s - %s",
			startDate.Format("02/01/2006"),
			endDate.Format("02/01/2006"))
	} else if req.TimeFrame == "" {
		// Nếu không có ngày bắt đầu và kết thúc, và không có time_frame,
		// thì mặc định là từ ngày hiện tại đến 30 ngày sau
		now := time.Now()
		endDate := now.AddDate(0, 0, 30)
		req.TimeFrame = fmt.Sprintf("%s - %s",
			now.Format("02/01/2006"),
			endDate.Format("02/01/2006"))
	}

	// 3. Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// 4. Gọi service
	created, msg := service.CreatePost(&req.PostJob, *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// 5. Trả về 201 + object vừa tạo
	response.Response(c, created, message.Success)
}

// GetJobPost godoc
// @Summary Lấy thông tin chi tiết bài đăng tuyển dụng
// @Description Lấy thông tin chi tiết của một bài đăng tuyển dụng theo ID
// @Tags JobPost
// @Accept json
// @Produce json
// @Param id path int true "Job Post ID"
// @Success 200 {object} response.VResponse{data=model.PostJob}
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/post-job/{id} [get]
func GetJobPost(c *gin.Context) {
	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Gọi service
	post, msg := service.GetPostByID(uint(id))
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về kết quả
	response.Response(c, post, message.Success)
}

// GetAllJobPosts godoc
// @Summary Lấy danh sách bài đăng tuyển dụng
// @Description Lấy danh sách tất cả các bài đăng tuyển dụng với phân trang
// @Tags JobPost
// @Accept json
// @Produce json
// @Param page query int false "Số trang (mặc định: 1)"
// @Param size query int false "Số lượng mỗi trang (mặc định: 10)"
// @Param status query string false "Lọc theo trạng thái (pending, approved, rejected)"
// @Success 200 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/post-job [get]
func GetAllJobPosts(c *gin.Context) {
	// Lấy tham số phân trang
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	status := c.Query("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	// Gọi service
	result, msg := service.GetAllPosts(page, size, status)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về kết quả
	response.Response(c, result, message.Success)
}

// GetMyJobPosts godoc
// @Summary Lấy danh sách bài đăng tuyển dụng của tôi
// @Description Lấy danh sách các bài đăng tuyển dụng của người dùng hiện tại
// @Tags JobPost
// @Accept json
// @Produce json
// @Param page query int false "Số trang (mặc định: 1)"
// @Param size query int false "Số lượng mỗi trang (mặc định: 10)"
// @Success 200 {object} response.VResponse
// @Failure 401 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/post-job/my-post [get]
// @Security BearerAuth
func GetMyJobPosts(c *gin.Context) {
	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Lấy tham số phân trang
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	// Gọi service
	result, msg := service.GetPostsByEmployer(*uid, page, size)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về kết quả
	response.Response(c, result, message.Success)
}

// UpdateJobPost godoc
// @Summary Cập nhật bài đăng tuyển dụng
// @Description Cập nhật thông tin của một bài đăng tuyển dụng
// @Tags JobPost
// @Accept json
// @Produce json
// @Param id path int true "Job Post ID"
// @Param payload body object true "Job Post Payload với start_date và end_date (YYYY-MM-DD)"
// @Success 200 {object} response.VResponse{data=model.PostJob}
// @Failure 400 {object} response.VResponse
// @Failure 401 {object} response.VResponse
// @Failure 403 {object} response.VResponse
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/post-job/{id} [put]
// @Security BearerAuth
func UpdateJobPost(c *gin.Context) {
	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Tạo struct để nhận dữ liệu từ frontend
	type PostJobRequest struct {
		model.PostJob
		StartDate string `json:"start_date"` // Ngày bắt đầu từ date-time-picker
		EndDate   string `json:"end_date"`   // Ngày kết thúc từ date-time-picker
	}

	// Bind JSON vào struct
	var req PostJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Response(c, err)
		return
	}

	// Đảm bảo ID trong path và body khớp nhau
	req.ID = uint(id)

	// Xử lý trường time_frame
	if req.StartDate != "" && req.EndDate != "" {
		// Parse ngày bắt đầu
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			response.Response(c, nil, message.Message{Message: "Invalid start date format. Expected YYYY-MM-DD", Code: 400})
			return
		}

		// Parse ngày kết thúc
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			response.Response(c, nil, message.Message{Message: "Invalid end date format. Expected YYYY-MM-DD", Code: 400})
			return
		}

		// Định dạng lại theo dd/mm/yyyy - dd/mm/yyyy
		req.TimeFrame = fmt.Sprintf("%s - %s",
			startDate.Format("02/01/2006"),
			endDate.Format("02/01/2006"))
	}

	// Gọi service
	updated, msg := service.UpdatePost(&req.PostJob, *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về kết quả
	response.Response(c, updated, message.Success)
}

// DeleteJobPost godoc
// @Summary Xóa bài đăng tuyển dụng
// @Description Xóa một bài đăng tuyển dụng (xóa mềm)
// @Tags JobPost
// @Accept json
// @Produce json
// @Param id path int true "Job Post ID"
// @Success 200 {object} response.VResponse
// @Failure 401 {object} response.VResponse
// @Failure 403 {object} response.VResponse
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/post-job/{id} [delete]
// @Security BearerAuth
func DeleteJobPost(c *gin.Context) {
	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Gọi service
	_, msg := service.DeletePost(uint(id), *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về kết quả
	response.Response(c, nil, message.Success)
}

// UpdateJobPostStatus godoc
// @Summary Cập nhật trạng thái bài đăng tuyển dụng
// @Description Cập nhật trạng thái của một bài đăng tuyển dụng (dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path int true "Job Post ID"
// @Param status body object true "Status Payload"
// @Success 200 {object} response.VResponse
// @Failure 400 {object} response.VResponse
// @Failure 401 {object} response.VResponse
// @Failure 403 {object} response.VResponse
// @Failure 404 {object} response.VResponse
// @Failure 500 {object} response.VResponse
// @Router /datn_backend/admin/post-job/{id}/status [put]
// @Security BearerAuth
func UpdateJobPostStatus(c *gin.Context) {
	// Lấy ID từ path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	// Bind JSON vào struct
	var req struct {
		IsApproved bool `json:"is_approved" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Response(c, err)
		return
	}

	// Kiểm tra giá trị status hợp lệ
	if req.IsApproved != true && req.IsApproved != false {
		response.Response(c, message.Message{Message: "Invalid status value", Code: http.StatusBadRequest})
		return
	}

	// Gọi service
	_, msg := service.UpdatePostStatus(uint(id), req.IsApproved, *uid)
	if msg != nil {
		response.Response(c, nil, msg)
		return
	}

	// Trả về kết quả
	response.Response(c, nil, message.Success)

}

// ApplyPostJob godoc
// @Summary     Ứng tuyển vào bài đăng tuyển dụng
// @Description Người tìm việc gửi file CV và thư xin việc
// @Tags        JobSeeker
// @Accept      multipart/form-data
// @Produce     json
// @Param       id            path     int                  true  "Job Post ID"
// @Param       resume_file   formData file                 true  "File CV (.pdf|.docx)"
// @Param       cover_letter  formData string               false "Cover letter"
// @Success     200           {object} response.VResponse
// @Failure     400           {object} response.VResponse
// @Failure     401           {object} response.VResponse
// @Failure     403           {object} response.VResponse
// @Failure     404           {object} response.VResponse
// @Failure     500           {object} response.VResponse
// @Router      /datn_backend/post-job/apply/{id} [post]
// @Security    BearerAuth
func ApplyPostJob(c *gin.Context) {
	// 1. Lấy ID bài đăng từ path
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Response(c, message.ValidationError)
		return
	}

	// 2. Lấy userID từ JWT claim
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}

	resumeHdr, err := c.FormFile("resume_file")
	if err != nil {
		response.Response(c, nil, message.Message{Message: "No resume file.", Code: 400})
		return
	}
	coverLetter := c.PostForm("cover_letter")

	savedPath := "/uploads/" + resumeHdr.Filename
	if err := c.SaveUploadedFile(resumeHdr, "."+savedPath); err != nil {
		response.Response(c, nil, message.Message{Message: "Save file failed.", Code: 409})
		return
	}

	// 4. Gọi service tạo ứng tuyển
	application := model.JobApplication{
		PostJobID:   uint(postID),
		ResumeURL:   savedPath,
		CoverLetter: coverLetter,
	}
	if err := service.ApplyPostJob(uint(*uid), &application); err != nil {
		// Ví dụ service trả về lỗi nếu bài đã ứng tuyển rồi hoặc bài đã đóng
		response.Response(c, nil, err)
		return
	}

	// 5. Thành công
	response.Response(c, nil, message.Success)
	return
}
