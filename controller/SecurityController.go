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
// @Router /sso/security/login [post]
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
