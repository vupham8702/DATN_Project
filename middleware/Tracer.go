package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
			c.Request.Header.Set("X-Trace-ID", traceID)
		}
		c.Set("traceID", traceID)
		c.Next()
	}
}
