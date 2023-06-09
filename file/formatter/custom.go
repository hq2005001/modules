package formatter

import (
	"fmt"
	"github.com/hq2005001/modules/file/upload"
	"github.com/hq2005001/modules/i18n"
	"github.com/hq2005001/modules/response"

	"github.com/gin-gonic/gin"
)

var Custom = new(customFormatter)

type customFormatter struct {
}

type customResponse struct {
	UID    int    `json:"uid"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

func (c2 customFormatter) Format(
	c *gin.Context,
	locale *i18n.I18N,
	group string,
	domain string,
	fileInfo *upload.UploadFileInfo,
) {
	response.New(c, locale).Data(customResponse{
		UID:    0,
		Name:   fileInfo.OriginName,
		URL:    fmt.Sprintf("%s%s", domain, fileInfo.FilePath),
		Status: "done",
	})
}

func (c2 customFormatter) FormatMulti(
	c *gin.Context,
	locale *i18n.I18N,
	group string,
	domain string,
	filenames map[string]upload.UploadFileInfo,
) {
	//TODO implement me
	var rs = make([]customResponse, 0, len(filenames))
	count := 0
	for k, v := range filenames {
		count++
		rs = append(rs, customResponse{
			UID:    count,
			Name:   k,
			URL:    fmt.Sprintf("%s%s", domain, v.FilePath),
			Status: "done",
		})
	}
	response.New(c, locale).Data(rs)
}
