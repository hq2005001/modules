package middlewares

import (
	"bytes"
	"github.com/hq2005001/modules/logger"
	"github.com/hq2005001/modules/utils"

	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Logger(logf *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		w := &responseBodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: c.Writer,
		}

		c.Writer = w

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		start := time.Now()
		c.Next()

		cost := time.Since(start)
		responseStatus := c.Writer.Status()

		logFields := []zap.Field{
			zap.Int("status", responseStatus),
			zap.String("request", fmt.Sprintf("[%s] %s", c.Request.Method, c.Request.URL)),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", utils.GetClientIP(c)),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("time", utils.MicrosecondsStr(cost)),
			zap.Any("headers", c.Request.Header),
		}

		if utils.InArray([]string{"POST", "PUT", "DELETE"}, c.Request.Method) && c.Request.MultipartForm == nil {
			logFields = append(logFields, zap.String("body", string(requestBody)))
		}

		if c.Request.Header.Get("Content-Type") == "application/json" {
			logFields = append(logFields, zap.Any("response", w.body))
		}

		if responseStatus > 400 && responseStatus < 500 {
			logf.Warn("HTTP Warning "+cast.ToString(responseStatus), logFields...)
		} else if responseStatus > 500 && responseStatus <= 599 {
			logf.Error("HTTP Error "+cast.ToString(responseStatus), logFields...)
		} else {
			logf.Debug("HTTP Access Log", logFields...)
		}
	}
}
