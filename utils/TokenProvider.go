package utils

import (
	"datn_backend/config"
	"datn_backend/message"
	"datn_backend/payload/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateToken(uid uint, issupper bool, roles interface{}, userType string) response.UserToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      fmt.Sprintf("%d", uid),
		"issupper": issupper,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"type":     userType,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 168).Unix(),
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	refreshTokenString, err := refreshToken.SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return response.UserToken{Error: err}
	}

	return response.UserToken{Token: tokenString, RefreshToken: refreshTokenString, Error: err}
}

func ValidateToken(c *gin.Context) (*jwt.Token, *message.Message) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Abort()
		return nil, &message.UnAuthorizedError
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.Abort()
		return nil, &message.UnAuthorizedError
	}
	tokenDecode, errVerify := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SECRET_KEY), nil
	})

	if errVerify != nil {
		return tokenDecode, &message.TokenExpired
	}

	return tokenDecode, nil
}

func GetToken(c *gin.Context, claims jwt.MapClaims) (interface{}, interface{}) {
	if uid, ok := claims["uid"].(string); ok {
		deviceid := claims["deviceid"]
		key := fmt.Sprintf("%s:%s:%s", config.TOKEN, uid, deviceid)
		token, err := config.RedisClient.Get(c, key).Result()
		if err != nil {
			return nil, message.NotFound
		}
		return token, nil
	}
	return nil, errors.New("Token does not exist!")
}
func VerifyToken(c *gin.Context) (jwt.MapClaims, *message.Message) {
	tokenDecode, err := ValidateToken(c)
	if err != nil {
		return nil, err
	}
	claims, ok := tokenDecode.Claims.(jwt.MapClaims)
	if !ok || !tokenDecode.Valid {
		return nil, &message.UnAuthorizedError
	}
	return claims, nil
}

func RefreshToken(token string, claims jwt.MapClaims) response.UserToken {
	tokenDecode, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return response.UserToken{Error: err}
	}
	if !tokenDecode.Valid {
		return response.UserToken{Error: errors.New("token decode invalid")}
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"deviceid":    claims["deviceid"],
		"issupper":    claims["issupper"],
		"uid":         claims["uid"],
		"roles":       claims["roles"],
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
		"resource_id": claims["resource_id"],
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 168).Unix(),
	})
	tokenString, errSigned := accessToken.SignedString([]byte(SECRET_KEY))
	refreshTokenString, _ := refreshToken.SignedString([]byte(SECRET_KEY))

	if errSigned != nil {
		log.Panic(errSigned)
		return response.UserToken{Error: errSigned}
	}

	return response.UserToken{Token: tokenString, RefreshToken: refreshTokenString, Error: err}
}

func GetClaimsByTokenExpired(c *gin.Context) (jwt.MapClaims, error) {
	tokenDecode, _ := ValidateToken(c)
	claims, _ := tokenDecode.Claims.(jwt.MapClaims)
	return claims, nil
}

func GetClaimsByToken(c *gin.Context) (jwt.MapClaims, error) {

	tokenDecode, err := ValidateToken(c)
	if err != nil {
		return nil, fmt.Errorf("Validate token error")
	}

	claims, _ := tokenDecode.Claims.(jwt.MapClaims)
	return claims, nil
}

func GetUidByClaim(c *gin.Context) (*uint, interface{}) {
	tokenDecode, err := ValidateToken(c)

	if err != nil {
		return nil, message.UnAuthorizedError
	}

	claims, _ := tokenDecode.Claims.(jwt.MapClaims)
	if uid, ok := claims["uid"].(string); ok {
		value, errParsing := strconv.ParseUint(uid, 10, 32)
		if errParsing != nil {
			return nil, message.UserActionNotFound
		}
		valueUint := uint(value)
		return &valueUint, nil
	}
	return nil, message.UserActionNotFound
}

func GetUserTypeInToken(c *gin.Context) (*string, interface{}) {
	claims, errorGetClaim := GetClaimsByToken(c)
	if errorGetClaim != nil {
		return nil, errorGetClaim
	}
	if provider, ok := claims["type"].(string); ok {
		return &provider, nil
	}
	return nil, message.UserActionNotFound
}
func GetFieldInToken(c *gin.Context, key string) (*string, interface{}) {
	claims, errorGetClaim := GetClaimsByToken(c)
	if errorGetClaim != nil {
		return nil, errorGetClaim
	}
	value, err := parseClaimToString(claims, key)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
func parseClaimToString(claims jwt.MapClaims, key string) (string, *message.Message) {
	if value, ok := claims[key]; ok {
		switch v := value.(type) {
		case string:
			return v, nil
		case float64:
			return fmt.Sprintf("%.0f", v), nil
		case int:
			return fmt.Sprintf("%d", v), nil
		case uint:
			return fmt.Sprintf("%d", v), nil
		case bool:
			return fmt.Sprintf("%t", v), nil
		default:
			return "", &message.FieldNotExist
		}
	}
	return "", &message.FieldNotExist
}
