package upload

import (
	"errors"
	"fmt"
	"github.com/hq2005001/modules/exception"
	"github.com/hq2005001/modules/i18n"
	"github.com/hq2005001/modules/response"
	"github.com/hq2005001/modules/utils"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// ResponseFormatter 响应格式
type ResponseFormatter func(c *gin.Context, group, path, name string)
type ResponseFormatter1 func(c *gin.Context, group, path, name string) (string, error)
type Formatter interface {
	Format(c *gin.Context, locale *i18n.I18N, group string, domain string, fileInfo *UploadFileInfo)
	FormatMulti(c *gin.Context, locale *i18n.I18N, group string, domain string, filenames map[string]UploadFileInfo)
}

// Upload 上传
func Upload(
	key string,
	locale *i18n.I18N,
	formatter Formatter,
	domain string,
	allowTypes ...string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		group := c.PostForm("group")
		if group == "" {
			group = "upload"
		}
		form, err := c.MultipartForm()
		if err != nil {
			response.New(c, locale).Fail(exception.NewParamsError(nil))
			return
		}
		var uploadInfo *UploadFileInfo
		urlMap := make(uploadFileInfos, 0, len(form.File[key]))
		isMulti := false
		if len(form.File[key]) > 1 {
			isMulti = true
		}
		for _, file := range form.File[key] {
			uploadInfo, err = saveFile(c, file, group, allowTypes...)
			if err != nil {
				continue
			}
			urlMap = append(urlMap, *uploadInfo)
		}
		if len(urlMap) == 0 {
			response.New(c, locale).Fail(exception.NewParamsError(nil))
			return
		}

		if formatter == nil {
			if isMulti {
				response.New(c, locale).Data(urlMap.FilePaths())
				return
			}
			response.New(c, locale).Data(uploadInfo.FilePath)
		} else if isMulti {
			formatter.FormatMulti(c, locale, group, domain, urlMap.FileMap())
		} else {
			formatter.Format(c, locale, group, domain, uploadInfo)
		}
		return
	}
}

func getOriginFilename(f *multipart.FileHeader) string {
	filename := f.Filename
	suffix := path.Ext(filename)
	return strings.TrimSuffix(filename, suffix)
}

type UploadFileInfo struct {
	Filename   string
	FilePath   string
	OriginName string
}

type uploadFileInfos []UploadFileInfo

func (u uploadFileInfos) FilePaths() []string {
	var rs = make([]string, 0, len(u))
	for _, item := range u {
		rs = append(rs, item.FilePath)
	}
	return rs
}

func (u uploadFileInfos) FileMap() map[string]UploadFileInfo {
	var rs = make(map[string]UploadFileInfo)
	for _, item := range u {
		rs[item.Filename] = item
	}
	return rs
}

func saveFile(c *gin.Context, file *multipart.FileHeader, group string, allowTypes ...string) (*UploadFileInfo, error) {

	f, err := file.Open()
	defer f.Close()
	if err != nil {
		return nil, errors.New("params error")
	}
	byt, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.New("params error")
	}
	fileType := new(Type).GetType(byt[:10])
	var typAllow = false
	for _, typ := range allowTypes {
		typAllow = typAllow || typ == fileType
	}
	if !typAllow {
		return nil, errors.New("upload file type not allow")
	}
	id := utils.UniqueID()
	cwd, _ := os.Getwd()
	uploadDir := cwd + "/public/upload/" + group
	if !utils.IsDir(uploadDir) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	originName := getOriginFilename(file)
	fileName := fmt.Sprintf("%s/%d.%s", uploadDir, id, fileType)
	err = c.SaveUploadedFile(file, fileName)
	if err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("/upload/%s/%d.%s", group, id, fileType)
	filePath = strings.ReplaceAll(filePath, "//", "/")
	return &UploadFileInfo{
		Filename:   fileName,
		FilePath:   filePath,
		OriginName: originName,
	}, nil
}
