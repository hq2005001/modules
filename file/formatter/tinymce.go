package formatter

import (
	"fmt"
	"github.com/hq2005001/modules/file/upload"
	"github.com/hq2005001/modules/i18n"
	"github.com/hq2005001/modules/response"

	"github.com/gin-gonic/gin"
)

var Tinymce = new(tinymceFormatter)

type tinymceFormatter struct {
}

type tinymceResponse struct {
	Name string `json:"name"`
	URL  string `json:"location"`
}

func (c2 tinymceFormatter) Format(
	c *gin.Context,
	locale *i18n.I18N,
	group string,
	domain string,
	fileInfo *upload.UploadFileInfo,
) {
	response.New(c, locale).Json(tinymceResponse{
		Name: fileInfo.OriginName,
		URL:  fmt.Sprintf("%s%s", domain, fileInfo.FilePath),
	})
}

func (c2 tinymceFormatter) DryFormat(
	c *gin.Context,
	locale *i18n.I18N,
	group string,
	domain string,
	fileInfo *upload.UploadFileInfo,
) string {
	return fmt.Sprintf("%s%s", domain, fileInfo.FilePath)
}

func (c2 tinymceFormatter) FormatMulti(
	c *gin.Context,
	locale *i18n.I18N,
	group string,
	domain string,
	filenames map[string]upload.UploadFileInfo,
) {
	var rs = make([]tinymceResponse, 0, len(filenames))
	count := 0
	for k, v := range filenames {
		count++
		rs = append(rs, tinymceResponse{
			Name: k,
			URL:  fmt.Sprintf("%s%s", domain, v.FilePath),
		})
	}
	response.New(c, locale).Json(rs)
}
