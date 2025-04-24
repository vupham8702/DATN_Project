package controller

import (
	"datn_backend/message"
	"datn_backend/middleware"
	"datn_backend/payload"
	"datn_backend/payload/response"
	"datn_backend/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// Login godoc
// @Summary Login
// @Description Login
// @Tags SecurityController
// @Accept json
// @Produce json
// @Param Login body payload.UserLogin true "User credentials"
// @Success 200 {object} response.VResponse
// @Router /datn_backend/security/login [post]
// @Security BearerAuth
func Login(c *gin.Context) {
	var userLogin payload.UserLogin
	deviceID := uuid.New().String()
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		middleware.Log(err)
		response.Response(c, message.Message{Message: message.ValidationError, Code: http.StatusBadRequest})
		return
	}
	token, err := service.Login(c, &userLogin, deviceID)
	if err != nil {
		response.Response(c, err)
		return
	}
	response.Response(c, token, message.Success)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, password and full name
// @Tags SecurityController
// @Accept json
// @Produce json
// @Param Register body payload.UserRegister true "User registration data"
// @Success 200 {object} response.VResponse
// @Router /datn_backend/security/register [post]
func Register(c *gin.Context) {
	var userRegister payload.UserRegister
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		middleware.Log(err)
		response.Response(c, message.Message{Message: message.ValidationError, Code: http.StatusBadRequest})
		return
	}

	result, err := service.Register(c, &userRegister)
	if err != nil {
		response.Response(c, err)
		return
	}

	response.Response(c, result, message.RegistrationSuccess)
}

// VerifyEmail godoc
// @Summary Verify user email
// @Description Verify user email with token
// @Tags SecurityController
// @Accept json
// @Produce json
// @Param VerifyEmail body payload.VerifyEmail true "Email verification data"
// @Success 200 {object} response.VResponse
// @Router /datn_backend/security/verify-email [post]
func VerifyEmail(c *gin.Context) {
	var verifyEmail payload.VerifyEmail
	if err := c.ShouldBindJSON(&verifyEmail); err != nil {
		middleware.Log(err)
		response.Response(c, message.Message{Message: message.ValidationError, Code: http.StatusBadRequest})
		return
	}

	result, err := service.VerifyEmail(c, &verifyEmail)
	if err != nil {
		response.Response(c, err)
		return
	}

	response.Response(c, result, message.EmailVerifySuccess)
}
