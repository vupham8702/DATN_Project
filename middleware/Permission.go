package middleware

import (
	"datn_backend/message"
	"datn_backend/middleware/handler"
	"datn_backend/payload/response"
	"datn_backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"slices"
)

func Permission(requiredPermission []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, errVerify := utils.VerifyToken(c)
		if errVerify != nil {
			response.Response(c, errVerify)
			c.Abort()
			return
		}

		_, err := utils.GetToken(c, claims)

		if err != nil {
			response.Response(c, message.TokenExpired)
			c.Abort()
			return
		}
		hasPermission, _ := hasPermission(c, claims, requiredPermission)

		if !hasPermission {
			response.Response(c, message.ForbidenError)
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasPermission(c *gin.Context, claims jwt.MapClaims, requiredPermission []string) (bool, interface{}) {
	var rolesName []string
	hasPermission := false
	isSuper := claims["issupper"]
	if isSuper == true {
		return true, nil
	}
	if roles, ok := claims["roles"].([]interface{}); ok {
		for _, role := range roles {
			if roleStr, ok := role.(string); ok {
				rolesName = append(rolesName, roleStr)
			}
		}
	}
	userPermissions, errGet := handler.GetPermissionsByRoleHandler(c, rolesName)
	if errGet != nil {
		return hasPermission, errGet
	}

	for _, p := range userPermissions {
		if slices.Contains(requiredPermission, p) {
			hasPermission = true
			break
		}
	}

	return hasPermission, nil
}
