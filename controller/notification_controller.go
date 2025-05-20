package controller

//import (
//	"net/http"
//	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
//	"datn_backend/payload"
//	"datn_backend/payload/response"
//	"datn_backend/service"
//	"datn_backend/utils"
//)
//
//// NotificationController xử lý các request liên quan đến thông báo
//type NotificationController struct {
//	svc *service.NotificationService
//}
//
//// NewNotificationController tạo instance mới của NotificationController
//func NewNotificationController(s *service.NotificationService) *NotificationController {
//	return &NotificationController{svc: s}
//}
//
//// GetUserNotifications lấy danh sách thông báo của người dùng
//func (ctr *NotificationController) GetUserNotifications(ctx *gin.Context) {
//	claims, ok := ctx.Get("claims")
//	if !ok {
//		response.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
//		return
//	}
//
//	user := claims.(*utils.JWTClaims)
//	userID, err := uuid.Parse(user.UserID)
//	if err != nil {
//		response.ErrorResponse(ctx, http.StatusBadRequest, "invalid user id")
//		return
//	}
//
//	page, limit := utils.GetPaginationParams(ctx)
//	notifications, total, err := ctr.svc.GetUserNotifications(userID, page, limit)
//	if err != nil {
//		response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	response.PaginationResponse(ctx, http.StatusOK, "notifications retrieved successfully", notifications, page, limit, total)
//}
//
//// MarkAsRead đánh dấu thông báo đã đọc
//func (ctr *NotificationController) MarkAsRead(ctx *gin.Context) {
//	claims, ok := ctx.Get("claims")
//	if !ok {
//		response.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
//		return
//	}
//
//	user := claims.(*utils.JWTClaims)
//	userID, err := uuid.Parse(user.UserID)
//	if err != nil {
//		response.ErrorResponse(ctx, http.StatusBadRequest, "invalid user id")
//		return
//	}
//
//	notificationID := ctx.Param("id")
//	notifID, err := uuid.Parse(notificationID)
//	if err != nil {
//		response.ErrorResponse(ctx, http.StatusBadRequest, "invalid notification id")
//		return
//	}
//
//	if err := ctr.svc.MarkAsRead(userID, notifID); err != nil {
//		response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	response.SuccessResponse(ctx, http.StatusOK, "notification marked as read", nil)
//}
//
//// MarkAllAsRead đánh dấu tất cả thông báo đã đọc
//func (ctr *NotificationController) MarkAllAsRead(ctx *gin.Context) {
//	claims, ok := ctx.Get("claims")
//	if !ok {
//		response.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
//		return
//	}
//
//	user := claims.(*utils.JWTClaims)
//	userID, err := uuid.Parse(user.UserID)
//	if err != nil {
//		response.ErrorResponse(ctx, http.StatusBadRequest, "invalid user id")
//		return
//	}
//
//	if err := ctr.svc.MarkAllAsRead(userID); err != nil {
//		response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	response.SuccessResponse(ctx, http.StatusOK, "all notifications marked as read", nil)
//}
//
//// DeleteNotification xóa một thông báo
//func (ctr *NotificationController) DeleteNotification(ctx *gin.Context) {
//	claims, ok := ctx.Get("claims")
//	if !ok {
//		response.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
//		return
//	}
//
//	user := claims.(*utils.JWTClaims)
//	userID, err := uuid.Parse(user.UserID)
//	if err != nil {
//		response.ErrorResponse(ctx, http.StatusBadRequest, "invalid user id")
//		return
//	}
//
//	notificationID := ctx.Param("id")
//	notifID, err := uuid.Parse(notificationID)
//	if err != nil {
//		response.ErrorResponse(ctx, http.StatusBadRequest, "invalid notification id")
//		return
//	}
//
//	if err := ctr.svc.DeleteNotification(userID, notifID); err != nil {
//		response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	response.SuccessResponse(ctx, http.StatusOK, "notification deleted successfully", nil)
//}
//
//// CreateEmployerNotification tạo thông báo cho nhà tuyển dụng (chỉ admin)
//func (ctr *NotificationController) CreateEmployerNotification(ctx *gin.Context) {
//	claims, ok := ctx.Get("claims")
//	if !ok {
//		response.ErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
//		return
//	}
//
//	user := claims.(*utils.JWTClaims)
//	if user.Role != "admin" {
//		response.ErrorResponse(ctx, http.StatusForbidden, "only admin can create employer notifications")
//		return
//	}
//
//	var req payload.CreateNotificationRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		response.ValidationErrorResponse(ctx, err)
//		return
//	}
//
//	employerID, err := uuid.Parse(req.RecipientID)
//	if err != nil {
//		response.ErrorResponse(ctx, http.StatusBadRequest, "invalid employer id")
//		return
//	}
//
//	notification, err := ctr.svc.CreateEmployerNotification(employerID, req.Title, req.Content, req.Type)
//	if err != nil {
//		response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	response.SuccessResponse(ctx, http.StatusCreated, "notification created successfully", notification)
//}null
