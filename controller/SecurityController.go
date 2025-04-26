package controller

import (
	m "datn_backend/domain/model"
	"datn_backend/message"
	"datn_backend/middleware"
	"datn_backend/payload"
	"datn_backend/payload/response"
	"datn_backend/service"
	"datn_backend/utils"
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
		response.Response(c, message.ValidationError, http.StatusBadRequest)
		return
	}

	result, err := service.Register(c, &userRegister)
	if err != nil {
		response.Response(c, err)
		return
	}

	response.Response(c, result, message.RegistrationSuccess)
}

// ApproveEmployer godoc
// @Summary Approve or reject an employer account
// @Description Admin can approve or reject employer account registrations
// @Tags AdminController
// @Accept json
// @Produce json
// @Param ApproveEmployer body payload.ApproveEmployer true "Approval data"
// @Success 200 {object} response.VResponse
// @Router /datn_backend/security/approve-employer [post]
// @Security BearerAuth
func ApproveEmployer(c *gin.Context) {
	// Check if the current user is an admin
	//user, exists := c.Get("currentUser")
	//if !exists || !user.(m.User).IsSupper {
	//	response.Response(c, message.Message{Message: "Unauthorized access", Code: http.StatusUnauthorized})
	//	return
	//}
	uid, errGet := utils.GetUidByClaim(c)
	if errGet != nil {
		response.Response(c, errGet)
		return
	}
	isSupper, errSuper := utils.GetFieldInToken(c, "issupper")
	if errSuper != nil {
		response.Response(c, message.InternalServerError)
		return
	}
	if *isSupper == "false" {
		response.Response(c, message.Message{Message: "Unauthorized access", Code: http.StatusUnauthorized})
		return
	}

	var approveRequest payload.ApproveEmployer
	if err := c.ShouldBindJSON(&approveRequest); err != nil {
		middleware.Log(err)
		response.Response(c, message.Message{Message: message.ValidationError, Code: http.StatusBadRequest})
		return
	}

	result, err := service.ApproveEmployer(c, &approveRequest, uid)
	if err != nil {
		response.Response(c, err)
		return
	}

	response.Response(c, result, message.Success)
}

// GetPendingEmployers godoc
// @Summary Get a list of employers pending approval
// @Description Admin can view all employer accounts waiting for approval
// @Tags AdminController
// @Accept json
// @Produce json
// @Success 200 {object} response.VResponse
// @Router /datn_backend/security/pending-employers [get]
// @Security BearerAuth
func GetPendingEmployers(c *gin.Context) {
	// Check if the current user is an admin
	user, exists := c.Get("currentUser")
	if !exists || !user.(m.User).IsSupper {
		response.Response(c, message.Message{Message: "Unauthorized access", Code: http.StatusUnauthorized})
		return
	}

	result, err := service.GetPendingEmployers()
	if err != nil {
		response.Response(c, err)
		return
	}

	response.Response(c, result, message.Success)
}
