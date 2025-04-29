package router

import (
	"datn_backend/controller"
	"github.com/gin-gonic/gin"
)

func ProfileRouter(routerGroup *gin.RouterGroup) {
	profileGroup := routerGroup.Group("/profile")
	{
		profileGroup.GET("", controller.GetUserProfile)
		profileGroup.PUT("/jobseeker", controller.UpdateJobseekerProfile)
		//profileGroup.PUT("/employer", middleware.Permission([]string{"DEFAULT_USER"}), controller.UpdateEmployerProfile) phải làm luồng assign role cho role đẩy ln redis trước rôid mới thêm đc middleware này vào
		profileGroup.PUT("/employer", controller.UpdateEmployerProfile)
		profileGroup.POST("/upload-photo", controller.UploadProfilePhoto)
		// API công khai không yêu cầu đăng nhập
		profileGroup.GET("/jobseeker/:id", controller.GetJobseekerPublicProfile)
		profileGroup.GET("/employer/:id", controller.GetEmployerPublicProfile)
	}
}
