package handler

import (
	"datn_backend/config"
	"datn_backend/message"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

const KEY_PREFIX = config.ROLE + ":%s:" + config.PERMISSIONS

func SetRolePermissionHandler(c *gin.Context, roleId string, value []string) interface{} {
	key := fmt.Sprintf(KEY_PREFIX, roleId)
	permissionsJson, err := json.Marshal(value)
	if err != nil {
		return message.ExcuteDatabaseError
	}
	status := config.RedisClient.Set(c, key, permissionsJson, 0)
	if (status.Val()) == config.OK {
		return nil
	}
	fmt.Errorf("Lỗi lưu role permissions redis tại : %s với key : %s", "SetRolePermissionHandler", key)
	return message.InternalServerError
}
func GetRolePermissionHandler(c *gin.Context, roleId string) (interface{}, interface{}) {
	key := fmt.Sprintf(KEY_PREFIX, roleId)
	permissionsData, err := config.RedisClient.Get(c, key).Result()
	if err != nil {
		(fmt.Errorf(fmt.Sprintf("Lỗi get role permissions redis tại : %s với key : %s", "GetRolePermissionHandler", key)))
		return nil, message.NotFound
	}
	var list []string
	errUnMarhshal := json.Unmarshal([]byte(permissionsData), &list)
	if errUnMarhshal != nil {
		(fmt.Errorf("Lỗi errUnMarhshal tại : %s với key : %s", "GetRolePermissionHandler", key))
		return nil, message.InternalServerError
	}
	return list, nil
}
func DeleteRolePermissionHandler(c *gin.Context, roleId string) interface{} {
	key := fmt.Sprintf(KEY_PREFIX, roleId)
	status := config.RedisClient.Del(c, key)
	if status == nil {
		(fmt.Errorf("Lỗi xóa quyền trên redis tại : %s với key : %s", "DeleteRolePermissionHandler", key))
		return message.NotFound
	}
	return nil
}
func GetPermissionsByRoleHandler(c *gin.Context, roleIds []string) ([]string, interface{}) {
	var list []string
	for _, roleId := range roleIds {
		key := fmt.Sprintf(KEY_PREFIX, roleId)
		permissionsData, err := config.RedisClient.Get(c, key).Result()
		if err != nil {
			(fmt.Errorf("Not found %s: %v", roleId, err))
			return nil, message.ExcuteDatabaseError
		}
		var permissions []string
		errUnmarshal := json.Unmarshal([]byte(permissionsData), &permissions)
		if errUnmarshal != nil {
			(fmt.Errorf("failed to unmarshal permissions for role %s: %v", roleId, err))
			return nil, message.ValidationError
		}
		list = append(list, permissions...)
	}
	return list, nil
}
