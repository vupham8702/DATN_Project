package response

import (
	"datn_backend/message"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VResponse struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Response(c *gin.Context, response ...interface{}) {
	vResponse := VResponse{}

	for _, item := range response {
		switch item.(type) {
		case message.Message:
			vResponse.Code = item.(message.Message).Code
			vResponse.Message = item.(message.Message).Message.(string)
		case *message.Message:
			vResponse.Code = item.(*message.Message).Code
			vResponse.Message = item.(*message.Message).Message.(string)
		default:
			if vResponse.Data == nil {
				vResponse.Data = item
			}
			break
		}
	}

	c.JSON(http.StatusOK, vResponse)
}
