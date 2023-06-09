package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/hq2005001/modules/exception"
	"github.com/hq2005001/modules/i18n"
	"github.com/hq2005001/modules/response"
)

type Handler[T any, U any] func(ctx *gin.Context, param T) (U, exception.Exception)

func ResponseWrapper[T any, U any](handler Handler[T, U], hasQuery bool, locale *i18n.I18N) gin.HandlerFunc {
	return func(c *gin.Context) {
		var param T
		if err := c.ShouldBind(&param); err != nil {
			response.New(c, locale).Fail(exception.NewParamsError(err))
			return
		}
		if hasQuery {
			if err := c.ShouldBind(&param); err != nil {
				response.New(c, locale).Fail(exception.NewParamsError(err))
				return
			}
		}
		result, err := handler(c, param)
		if err != nil {
			response.New(c, locale).Fail(err)
			return
		}
		response.New(c, locale).Data(result)
	}
}
