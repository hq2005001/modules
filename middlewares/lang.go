package middlewares

//import (
//	"email/i18n"
//
//	"github.com/gin-gonic/gin"
//	"golang.org/x/text/language"
//)
//
//// I18N 设置语言中间件
//func I18N() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		i18n.Locale().SetLang(getAcceptLang(c))
//		c.Next()
//	}
//}
//
//func getAcceptLang(c *gin.Context) language.Tag {
//	lang := c.Query("lang")
//	if lang != "" {
//		tag, err := language.Parse(lang)
//		if err == nil {
//			return tag
//		}
//	}
//	tag, _ := language.MatchStrings(i18n.Locale().Matcher(), c.GetHeader("Accept-language"), language.Chinese.String())
//	return tag
//}
