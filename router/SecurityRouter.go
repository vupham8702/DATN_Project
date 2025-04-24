package router

import (
	"datn_backend/controller"
	"github.com/gin-gonic/gin"
)

func SecurityRouter(routerGroup *gin.RouterGroup) {
	securityGroup := routerGroup.Group("/security")
	{
		securityGroup.POST("/register", controller.Register)
		securityGroup.POST("/verify-email", controller.VerifyEmail)
		//securityGroup.PATCH("/sms/reset-password", controller.ResetPasswordSms)
		securityGroup.POST("/login", controller.Login)
		//securityGroup.POST("/guest/login", controller.GuestLogin)
		//securityGroup.GET("/logout", controller.Logout)
		//securityGroup.POST("/refresh-token", controller.RefreshToken)
		//securityGroup.POST("/forgot-password", controller.ForgotPassword)
		//securityGroup.GET("/verify-otp-reset-pass", controller.VerifyOtpResetPassword)
		//securityGroup.PATCH("/reset-password", controller.ResetPassword)
		//securityGroup.PATCH("/change-password", middleware.Permission([]string{"DEFAULT_USER"}), controller.ChangePassword)
		//securityGroup.POST("/send-login-trace", controller.SendLoginOtp)
		//securityGroup.POST("/login-otp", controller.LoginOtp)
		//securityGroup.POST("/change-password-otp", controller.ChangePasswordOtp)
		//securityGroup.GET("/get-domain-cdn", controller.GetDomainCdn)
	}
}
