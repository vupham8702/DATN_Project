package router

import (
	"datn_backend/controller"
	"github.com/gin-gonic/gin"
)

func PostJobRouter(routerGroup *gin.RouterGroup) {
	postJobGroup := routerGroup.Group("/post-job")
	{
		postJobGroup.POST("/create", controller.CreatePostJob)
		postJobGroup.GET("/:id", controller.GetJobPost)
		postJobGroup.GET("", controller.GetAllJobPosts)
		postJobGroup.GET("/my-post", controller.GetMyJobPosts)
		postJobGroup.PUT("/:id", controller.UpdateJobPost)
		postJobGroup.DELETE("/:id", controller.DeleteJobPost)
		postJobGroup.PUT("/:id/status", controller.UpdateJobPostStatus)
		postJobGroup.POST("/apply/:id", controller.ApplyPostJob)
		postJobGroup.GET("/my-applications", controller.GetMyApplications)
		postJobGroup.GET("/:id/applications", controller.GetJobApplications)
		postJobGroup.PUT("/applications/:id/status", controller.UpdateApplicationStatus)
	}
}
