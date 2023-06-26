package upload

import (
	"crypto/md5"
	"encoding/hex"
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
	DryFormat(c *gin.Context, locale *i18n.I18N, group string, domain string, fileInfo *UploadFileInfo) string
	FormatMulti(c *gin.Context, locale *i18n.I18N, group string, domain string, filenames map[string]UploadFileInfo)
}

type SaveHandler func(fileInfo *UploadFileInfo, checkOnly bool) (*UploadFileInfo, bool, error)

// Upload 上传
func Upload(
	key string,
	locale *i18n.I18N,
	saveHandler SaveHandler,
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
			uploadInfo, err = saveFile(c, locale, domain, file, group, saveHandler, formatter, allowTypes...)
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
	Digest     string
	Ext        string
	OriginName string
	URL        string
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

func saveFile(
	c *gin.Context,
	locale *i18n.I18N,
	domain string,
	file *multipart.FileHeader,
	group string,
	saveHandler SaveHandler,
	formatter Formatter,
	allowTypes ...string,
) (*UploadFileInfo, error) {

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
	d := md5.New()
	d.Write(byt)
	digest := hex.EncodeToString(d.Sum(nil))
	originName := getOriginFilename(file)
	fileName := fmt.Sprintf("%s/%d.%s", uploadDir, id, fileType)
	filePath := fmt.Sprintf("/upload/%s/%d.%s", group, id, fileType)
	filePath = strings.ReplaceAll(filePath, "//", "/")

	info := &UploadFileInfo{
		Filename:   fileName,
		Digest:     digest,
		Ext:        fileType,
		FilePath:   filePath,
		OriginName: originName,
	}
	needSave := true
	if saveHandler != nil {
		url := formatter.DryFormat(c, locale, group, domain, info)
		info.URL = url
		info, needSave, err = saveHandler(info, true)
		if err != nil {
			return nil, err
		}
	}
	if needSave {
		err = c.SaveUploadedFile(file, fileName)
		if err != nil {
			return nil, err
		}
		_, _, _ = saveHandler(info, false)
	}

	return info, nil
}
