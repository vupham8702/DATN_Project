package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var Logger *logrus.Entry

type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogResponse(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		contentLength, err := strconv.ParseInt(os.Getenv("MAX_CONTENT_LENGTH"), 10, 64)
		if err != nil {
			contentLength = 5242880

		}
		traceID := c.Request.Header.Get("X-Trace-ID")
		Logger = logger.WithField("traceId", traceID)
		requestBody, _ := io.ReadAll(c.Request.Body)
		requestBodyLog := make([]byte, 0)
		if c.Request.ContentLength <= contentLength/2 {
			requestBodyLog = requestBody
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		writer := &ResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()
		latency := time.Since(start)
		logLine := fmt.Sprintf(
			"REQUEST: %s | RESPONSE: %s\n",
			string(requestBodyLog),
			writer.body.String(),
		)
		if !strings.Contains(c.Request.URL.Path, "swagger") {
			logger.WithField("traceId", traceID).
				WithField("startDate", start.Format(time.DateTime)).
				WithField("method", c.Request.Method).
				WithField("url", c.Request.URL.Path).
				WithField("status", c.Writer.Status()).
				WithField("latency", fmt.Sprintf("%.2f ms", float64(latency.Nanoseconds())/1e6)).
				WithField("clientIP", c.ClientIP()).
				Info(logLine)
		}
	}
}

func Log(error interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Logger.WithField("file", file).
		WithField("line", line).
		Info(error)
}
