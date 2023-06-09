package response

import (
	"github.com/hq2005001/modules/exception"
	"github.com/hq2005001/modules/i18n"
	"net/http"

	"github.com/gin-gonic/gin"
)

const Ok = "ok"

type Result struct {
	Code exception.Code `json:"code"`
	Msg  string         `json:"msg"`
	Data interface{}    `json:"data"`
}

type Response struct {
	c      *gin.Context
	params map[string]interface{}
	locale *i18n.I18N
}

func (r *Response) WithParams(params map[string]interface{}) *Response {
	r.params = params
	return r
}

func (r *Response) Data(data interface{}) {
	r.c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  r.locale.Tr("success", nil),
		Data: data,
	})
}

func (r *Response) Json(data interface{}) {
	r.c.JSON(http.StatusOK, data)
}

func (r *Response) Fail(msg exception.Exception) {
	r.c.JSON(http.StatusOK, Result{
		Code: msg.Code(),
		Msg:  r.locale.Tr(msg.Msg(), r.params),
	})
}

func (r *Response) Error(statusCode int, msg exception.Exception) {
	r.c.JSON(statusCode, Result{
		Code: msg.Code(),
		Msg:  r.locale.Tr(msg.Msg(), r.params),
	})
}

func New(c *gin.Context, locale *i18n.I18N) *Response {
	return &Response{
		c:      c,
		locale: locale,
	}
}
