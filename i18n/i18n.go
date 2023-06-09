package i18n

import (
	"bytes"
	"embed"
	"github.com/hq2005001/modules/logger"
	"github.com/hq2005001/modules/utils"
	"go.uber.org/zap"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/language"
)

// Translator 翻译器
var (
	//Translator  *I18N
	DefaultLang = language.Chinese
)

//func init() {
//	absDir := utils.GetAbsDir()
//	Translator = New().Load(fmt.Sprintf("%s/public/i18n", absDir), yaml.Unmarshal)
//}

// Unmarshal 解析
type Unmarshal func(data []byte, v interface{}) error

// I18N 国际化
type I18N struct {
	Lang     language.Tag
	Messages map[language.Tag]map[string]string
	matcher  language.Matcher
	Tags     []language.Tag
	logger   *logger.Logger
}

// SetLang 设置国际化语言
func (i *I18N) SetLang(lang language.Tag) *I18N {
	i.Lang = lang
	return i
}

// Matcher  匹配器
func (i *I18N) Matcher() language.Matcher {
	return i.matcher
}

// Load 加载配置文件
func (i *I18N) Load(path string, unmarshal Unmarshal) *I18N {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		fileInfo := strings.Split(info.Name(), ".")
		filename := fileInfo[0]
		content, _ := ioutil.ReadFile(path)
		result := make(map[string]string)
		err = unmarshal(content, &result)

		if err != nil {
			return err
		}

		// 解析语言为tag
		t, err := language.Parse(filename)
		if err != nil {
			return err
		}
		i.Messages[t] = result
		i.Tags = append(i.Tags, t)
		return err
	})
	i.matcher = language.NewMatcher(i.Tags)
	return i
}

// LoadByFs 使用 fs加载
func (i *I18N) LoadByFs(dir embed.FS, unmarshal Unmarshal, allowExts []string) *I18N {
	_ = fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		fileInfo := strings.Split(d.Name(), ".")
		filename := fileInfo[0]
		ext := fileInfo[1]
		if utils.IsStrInArr(allowExts, ext) {
			return nil
		}
		content, _ := ioutil.ReadFile(path)
		result := make(map[string]string)
		err = unmarshal(content, &result)

		if err != nil {
			return err
		}

		// 解析语言为tag
		t, err := language.Parse(filename)
		if err != nil {
			return err
		}
		i.Messages[t] = result
		i.Tags = append(i.Tags, t)
		return err
	})
	i.matcher = language.NewMatcher(i.Tags)
	return i
}

// Tr 翻译
func (i *I18N) Tr(key string, params map[string]interface{}) string {
	var rs bytes.Buffer
	if tmpl, exist := i.Messages[i.Lang][key]; exist {
		t, err := template.New(key).Parse(tmpl)
		if err != nil {
			i.logger.Debug("解析多语言模板失败, ", zap.Error(err))
			return key
		}
		err = t.Execute(&rs, params)
		if err != nil {
			i.logger.Debug("执行解析多语言失败, ", zap.Error(err))
			return key
		}
		return rs.String()
	}
	return key
}

// New 创建新的翻译器
func New(logf *logger.Logger) *I18N {
	return &I18N{
		Lang:     DefaultLang,
		Messages: make(map[language.Tag]map[string]string),
		Tags:     make([]language.Tag, 0),
		logger:   logf,
	}
}
