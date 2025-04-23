package service

import (
	"datn_backend/config"
	"datn_backend/domain/model"
	repo "datn_backend/domain/repository"
	"datn_backend/message"
	"datn_backend/middleware"
	"datn_backend/payload"
	"datn_backend/payload/response"
	"datn_backend/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, userLogin *payload.UserLogin, deviceId string) (*response.UserToken, interface{}) {
	var user model.User
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
		return nil, message.UserHasBeenLocked
	}

	verify, _, err := utils.VerifyPassword(userLogin.Password, user.Password)
	if !verify || err != nil {
		return nil, message.PasswordNotCorrect
	}

	token, _ := CreateToken(c, &user, UserTypeByProviderForToken(&user))
	return token, nil
}
func UserTypeByProviderForToken(user *model.User) string {
	if user.Providers == nil || len(user.Providers) == 0 {
		return config.USER_TYPE_ANONYMOUS
	}
	for _, v := range user.Providers {
		if v.Provider == config.SYSTEM_ACC {
			return config.USER_TYPE_CMS
		}
	}
	return config.USER_TYPE_MOBILE

}

func CreateToken(c *gin.Context, user *model.User, userType string) (*response.UserToken, interface{}) {
	var roles []string

	for _, role := range user.Roles {
		roles = append(roles, fmt.Sprintf("%d", role.ID))
	}

	token := utils.GenerateToken(
		user.ID,
		user.IsSupper,
		roles,
		userType,
	)
	uidStr := fmt.Sprintf("%d", user.ID)
	errSaveToken := CreateTokenRedis(c, &token, uidStr)
	if errSaveToken != nil {
		return nil, message.ExcuteDatabaseError
	}

	return &token, nil
}

func CreateTokenRedis(c *gin.Context, token *response.UserToken, uid string) interface{} {
	key := config.TOKEN + ":" + uid
	value, err := json.Marshal(token)
	if err != nil {
		return message.InternalServerError
	}
	status := config.RedisClient.Set(c, key, value, 0)
	if status.Val() != config.OK {
		middleware.Log(fmt.Errorf("Save token error Redis ...."))
		return nil
	}
	return nil
}
