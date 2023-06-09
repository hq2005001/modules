package middlewares

import (
	"errors"
	"github.com/hq2005001/modules/exception"
	"github.com/hq2005001/modules/i18n"
	"github.com/hq2005001/modules/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ForceUA(locale *i18n.I18N) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("User-Agent") == "" {
			response.New(c, locale).Error(http.StatusBadRequest, exception.NewUAError(errors.New("invalid user-agent")))
			return
		}
		c.Next()
	}
}
