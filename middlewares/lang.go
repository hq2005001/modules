package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/hq2005001/modules/i18n"
	"golang.org/x/text/language"
)

// I18N 设置语言中间件
func I18N(n *i18n.I18N) gin.HandlerFunc {
	return func(c *gin.Context) {
		n.SetLang(getAcceptLang(c, n))
		c.Next()
	}
}

func getAcceptLang(c *gin.Context, n *i18n.I18N) language.Tag {
	lang := c.Query("lang")
	if lang != "" {
		tag, err := language.Parse(lang)
		if err == nil {
			return tag
		}
	}
	tag, _ := language.MatchStrings(n.Matcher(), c.GetHeader("Accept-language"), language.Chinese.String())
	return tag
}
