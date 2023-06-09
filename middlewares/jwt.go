package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hq2005001/modules/auth"
	"github.com/hq2005001/modules/i18n"
	"github.com/hq2005001/modules/response"
	"net/http"
)

// Jwt 检查jwt
func Jwt(conf *auth.JwtConf, contextKey string, claims func() jwt.Claims, isContinue bool, locale *i18n.I18N) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.Parser(conf).
			SetContextKey(contextKey).
			SetClaims(claims()).
			SetContinue(isContinue).Verify(c)
		if err != nil {
			response.New(c, locale).Error(http.StatusUnauthorized, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
