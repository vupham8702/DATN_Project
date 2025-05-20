package router

import (
	"datn_backend/controller"
	"github.com/gin-gonic/gin"
)

func CVRouter(routerGroup *gin.RouterGroup) {
	cvGroup := routerGroup.Group("/cv")
	{
		// Các API công khai
		cvGroup.GET("/templates", controller.GetAllCVTemplates)
		cvGroup.GET("/templates/:id", controller.GetCVTemplateByID)
		cvGroup.GET("/templates/:id/preview", controller.GetCVPreview)
		// Tải xuống mẫu CV (với thông tin cá nhân được điền vào)
		cvGroup.GET("/templates/:id/download", controller.DownloadAndFillCVTemplate)

		// Các API dành cho admin
		cvGroup.POST("/admin/templates", controller.CreateCVTemplate)
		cvGroup.PUT("/admin/templates/:id", controller.UpdateCVTemplate)
		cvGroup.DELETE("/admin/templates/:id", controller.DeleteCVTemplate)
		cvGroup.POST("/admin/templates/:id/thumbnail", controller.UploadCVThumbnail)
		cvGroup.GET("/admin/templates/:id/download-original", controller.DownloadCVTemplate)
	}
}
