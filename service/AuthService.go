package service

import (
	"datn_backend/config"
	"datn_backend/domain/model"
	repo "datn_backend/domain/repository"
	"datn_backend/message"
	"datn_backend/middleware"
	"datn_backend/payload"
	"datn_backend/payload/response"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, userLogin *payload.UserLogin, deviceId string) (*response.UserToken, interface{}) {
	var user model.User
	deviceID := c.Request.Header.Get("deviceId")
	deviceName := c.Request.Header.Get("deviceName")
	devicePlatform := c.Request.Header.Get("devicePlatform")
	fcmToken := c.Request.Header.Get("fcmToken")
	if userLogin.Type == config.Phone {
		db := config.DB
		tx := db.Begin()
		if &deviceID == nil || deviceID == "" {
			deviceID = deviceId
		}
		userPhone, err := repo.GetUserByPhone(userLogin.Username)
		if err != nil {
			tx.Rollback()
			middleware.Log(err)
			return nil, message.AccountDoesNotExist
		}
		if userPhone.IsLocked == true {
			return nil, message.UserHasBeenLocked
		}
		user = *userPhone
		errSave := SaveDeviceInfo(userPhone, deviceID, deviceName, devicePlatform, tx)
		if errSave != nil {
			tx.Rollback()
			return nil, errSave
		}
		if fcmToken != "" {
			errSaveToken := SaveFcmToken(tx, fcmToken, userPhone)
			if errSaveToken != nil {
				tx.Rollback()
				return nil, errSaveToken
			}
		}
		tx.Commit()
	} else {
		userMail, err := repo.GetUserByMail(userLogin.Username)
		if err != nil {
			return nil, message.EmailNotExist
		}
		user = *userMail
		if &user == nil {
			return nil, message.EmailNotExist
		}
		//if user.IsActive == false {
		//	return nil, message.EmailNotActive
		//}
		if user.IsLocked == true {
			return nil, sso_message.UserHasBeenLocked
		}
	}

	verify, _, err := utils.VerifyPassword(userLogin.Password, user.Password)
	if !verify || err != nil {
		return nil, message.PasswordNotCorrect
	}

	token, _ := CreateToken(c, &user, deviceID, UserTypeByProviderForToken(&user))
	return token, nil
}
