package service

import (
	"datn_backend/config"
	"datn_backend/domain/model"
	"datn_backend/domain/repository"
	"datn_backend/message"
	"datn_backend/middleware"
	"errors"
	"fmt"
	"strings"
	"time"
)

// CreatePost tạo mới một bài đăng tuyển dụng
func CreatePost(post *model.PostJob, userID uint) (interface{}, interface{}) {
	// Thiết lập các giá trị mặc định
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.CreatedBy = userID
	post.UpdatedBy = userID
	post.IsDeleted = false

	var typeMapping = map[string]string{
		"1": "Full-time",
		"2": "Part-time",
		"3": "Remote online",
		"4": "Freelance",
	}
	if val, ok := typeMapping[post.Type]; ok {
		post.Type = val
	}

	var genderMapping = map[string]string{
		"1": "Nam",
		"2": "Nũ",
		"3": "Không yêu cầu",
	}
	if val, ok := genderMapping[post.Gender]; ok {
		post.Gender = val
	}
	// Mặc định trạng thái là "pending" (chờ duyệt)
	if post.Status == "" {
		post.Status = "pending"
	}

	//lấy profile bằng Id
	profile, errProfile := repository.GetEmployerProfileByUserID(&userID)
	if errProfile != nil {
		middleware.Log(fmt.Errorf("Failed to get user profile: %v", errProfile))
		return nil, message.UserNotFound
	}

	post.Company = profile.CompanyName
	post.Logo = profile.CompanyLogo

	// Gọi repository để lưu vào database
	err := repository.CreatePost(post)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to create job post: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return post, nil
}

// GetPostByID lấy thông tin bài đăng theo ID
func GetPostByID(postID uint) (interface{}, interface{}) {
	post, err := repository.GetPostByID(postID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get job post: %v", err))
		return nil, message.PostNotFound
	}

	return post, nil
}

// GetAllPosts lấy danh sách tất cả các bài đăng
func GetAllPosts(page, pageSize int, status string) (interface{}, interface{}) {
	posts, total, err := repository.GetAllPosts(page, pageSize, status)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get job posts: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return map[string]interface{}{
		"posts": posts,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, nil
}

// GetPostsByEmployer lấy danh sách bài đăng của một nhà tuyển dụng
func GetPostsByEmployer(employerID uint, page, pageSize int) (interface{}, interface{}) {
	posts, total, err := repository.GetPostsByEmployer(employerID, page, pageSize)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get employer job posts: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return map[string]interface{}{
		"posts": posts,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, nil
}

// UpdatePost cập nhật thông tin bài đăng
func UpdatePost(post *model.PostJob, userID uint) (interface{}, interface{}) {
	// Kiểm tra xem bài đăng có tồn tại không
	existingPost, err := repository.GetPostByID(post.ID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get job post: %v", err))
		return nil, message.PostNotFound
	}

	// Kiểm tra quyền (chỉ người tạo mới có thể cập nhật)
	if existingPost.CreatedBy != userID {
		return nil, message.Message{Message: "You don't have permission to update this post", Code: 403}
	}

	// Cập nhật thông tin
	post.UpdatedAt = time.Now()
	post.UpdatedBy = userID

	// Giữ nguyên một số trường không cho phép cập nhật
	post.CreatedAt = existingPost.CreatedAt
	post.CreatedBy = existingPost.CreatedBy
	post.IsDeleted = existingPost.IsDeleted

	// Gọi repository để cập nhật
	err = repository.UpdatePost(post)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to update job post: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return post, nil
}

// DeletePost xóa mềm một bài đăng
func DeletePost(postID uint, userID uint) (interface{}, interface{}) {
	// Kiểm tra xem bài đăng có tồn tại không
	existingPost, err := repository.GetPostByID(postID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get job post: %v", err))
		return nil, message.PostNotFound
	}

	// Kiểm tra quyền (chỉ người tạo mới có thể xóa)
	if existingPost.CreatedBy != userID {
		return nil, message.Message{Message: "You don't have permission to delete this post", Code: 403}
	}

	// Gọi repository để xóa
	err = repository.DeletePost(postID, userID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to delete job post: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return nil, nil
}

// UpdatePostStatus cập nhật trạng thái của bài đăng (dành cho admin)
func UpdatePostStatus(postID uint, isApproved bool, userID uint) (interface{}, interface{}) {
	// Kiểm tra xem bài đăng có tồn tại không
	_, err := repository.GetPostByID(postID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get job post: %v", err))
		return nil, message.PostNotFound
	}

	// Gọi repository để cập nhật trạng thái
	err = repository.UpdatePostStatus(postID, isApproved, userID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to update job post status: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return nil, nil
}

func ApplyPostJob(userID uint, app *model.JobApplication) interface{} {
	// gán metadata chung
	app.CreatedBy = userID
	app.UpdatedBy = userID
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	app.Status = "pending"

	// TODO: kiểm tra xem user đã ứng tuyển chưa, hoặc bài còn mở hay không
	post, err := repository.GetPostByID(app.PostJobID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get job post: %v", err))
		return message.PostNotFound
	}

	open, err := IsOpen(post.TimeFrame)
	if !open || err != nil {
		return message.Message{Message: "The job application has been closed.", Code: 409}
	}

	existed, errCheck := repository.CheckIfUserApplied(userID, app.PostJobID)
	if errCheck != nil {
		return message.ExcuteDatabaseError
	}
	if existed {
		return message.Message{Message: "You have already applied for this job.", Code: 409}
	}

	// nếu hợp lệ, ghi vào DB:
	if err := config.DB.Create(app).Error; err != nil {
		return err
	}
	errUpdate := repository.IncrementApplicationCountRaw(app.PostJobID)
	if errUpdate != nil {
		return message.ExcuteDatabaseError
	}

	return nil
}
func ParseTimeFrame(tf string) (start, end time.Time, err error) {
	parts := strings.Split(tf, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, errors.New("invalid timeframe format")
	}
	layout := "02/01/2006"
	s := strings.TrimSpace(parts[0])
	e := strings.TrimSpace(parts[1])
	start, err = time.Parse(layout, s)
	if err != nil {
		return
	}
	end, err = time.Parse(layout, e)
	if err != nil {
		return
	}
	return
}
func IsOpen(tf string) (bool, error) {
	_, end, err := ParseTimeFrame(tf)
	if err != nil {
		return false, err
	}
	return time.Now().Before(end.Add(24 * time.Hour)), nil
}

// GetMyApplications lấy danh sách ứng tuyển của người dùng
func GetMyApplications(userID uint) (interface{}, interface{}) {
	applications, err := repository.GetJobApplicationsByUserID(userID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get job applications: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return applications, nil
}

// GetJobApplications lấy danh sách ứng tuyển cho bài đăng
func GetJobApplications(postJobID uint, userID uint) (interface{}, interface{}) {
	// Kiểm tra xem người dùng có sở hữu bài đăng không
	isOwner, err := repository.CheckIfUserOwnsPost(userID, postJobID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to check if user owns post: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	if !isOwner {
		return nil, message.ForbidenError
	}

	// Lấy danh sách ứng tuyển
	applications, err := repository.GetJobApplicationsByPostID(postJobID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get job applications: %v", err))
		return nil, message.ExcuteDatabaseError
	}

	return applications, nil
}

// UpdateApplicationStatus cập nhật trạng thái ứng tuyển
func UpdateApplicationStatus(applicationID uint, status string, userID uint) interface{} {
	// Lấy ID người sở hữu bài đăng
	postOwnerID, err := repository.GetPostOwnerByApplicationID(applicationID)
	if err != nil {
		middleware.Log(fmt.Errorf("Failed to get post owner: %v", err))
		return message.ExcuteDatabaseError
	}

	// Kiểm tra xem người dùng có quyền cập nhật không
	if postOwnerID != userID {
		return message.ForbidenError
	}

	// Cập nhật trạng thái
	if err := repository.UpdateJobApplicationStatus(applicationID, status); err != nil {
		middleware.Log(fmt.Errorf("Failed to update application status: %v", err))
		return message.ExcuteDatabaseError
	}

	return nil
}
