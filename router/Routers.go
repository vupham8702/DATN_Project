// Package router routes/Routers.go
package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"strconv"
)

func RegisterRoutes(router *gin.Engine) {
	apiPath := os.Getenv("API_PATH")
	prefixRoute := router.Group(apiPath)
	{
		debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
		if debug {
			prefixRoute.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
		//UserRouter(prefixRoute)
		//ContentRouter(prefixRoute)
		SecurityRouter(prefixRoute)
		ProfileRouter(prefixRoute)
		//ResourceRouter(prefixRoute)
		//PermissionRouter(prefixRoute)
		//RoleRouter(prefixRoute)
		//GoogleAuthRouter(prefixRoute)
		//AppleIDAuthRouter(prefixRoute)
		//SocialLoginRouter(prefixRoute)
		//SmsRouter(prefixRoute)
		//DeviceRouter(prefixRoute)
		//MobiAuthRouter(prefixRoute)
	}
}
