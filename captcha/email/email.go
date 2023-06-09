package email

import (
	"bytes"
	"context"
	"github.com/hq2005001/modules/captcha/captcha"
	"github.com/hq2005001/modules/captcha/conf"
	"github.com/hq2005001/modules/drivers/nsq"
	"github.com/hq2005001/modules/drivers/redis"

	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// Name 名称
	Name = "email"
)

// Email 邮箱
type Email struct {
	captcha.BaseCaptcha
	Config      *conf.Config
	QueueConfig *nsq.NsqConfig
	RedisConf   *redis.Conf
}

// Send 发送
func (e *Email) Send(ip, account, subject, content string) error {
	//ctx := context.TODO()
	//if !redis.New(e.RedisConf).SetNX(ctx, "email_lock_"+account, "1", time.Minute*1).Val() {
	//	return errors.New("too frequent")
	//}
	//redis.New(e.RedisConf).Set(ctx, e.GetCaptchaCacheKey(account), content, time.Minute*captcha.CacheMinute)
	//rs, err := e.GenerateContent(account, subject, content)
	//if err != nil {
	//	return err
	//}
	//e.LogSendNum(ctx, ip, account)
	//tasks.SendEmail(rs.(*tasks.MailContent))
	return nil
}

// GenerateContent 生成内容
func (e *Email) GenerateContent(account, subject, content string) (interface{}, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	templateFile := wd + "/views/email.html"
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer

	err = t.Execute(&buf, gin.H{"email": account, "code": content, "date": time.Now().Format("2006-01-02")})
	if err != nil {
		return nil, err
	}
	return "", nil
	//
	//return &tasks.MailContent{
	//	Email:   account,
	//	Subject: subject,
	//	Content: buf.String(),
	//}, nil
}

// Validate 验证
func (e *Email) Validate(account string, code string) bool {
	ctx := context.TODO()
	key := e.GetCaptchaCacheKey(account)
	if redis.New(e.RedisConf).Get(ctx, key).Val() == code {
		redis.New(e.RedisConf).Del(ctx, key)
		redis.New(e.RedisConf).Set(ctx, e.GetValidateKey(account), "1", time.Minute*captcha.CacheMinute)
		return true
	}
	e.LogErrorNum(ctx, account)
	return false
}

// IsValidate 是否通过验证
func (e *Email) IsValidate(account string, clear bool) bool {
	ctx := context.TODO()
	key := e.GetValidateKey(account)
	if redis.New(e.RedisConf).Get(ctx, key).Val() == "1" {
		if clear {
			redis.New(e.RedisConf).Del(ctx, key)
			redis.New(e.RedisConf).Del(ctx, e.GetValidateKey(account))
		}
		return true
	}
	return false
}

// GetValidateKey 得到验证的键
func (e *Email) GetValidateKey(account string) string {
	return fmt.Sprintf("email_validate::%s", account)
}

// GetCaptchaCacheKey 得到验证码缓存键
func (e *Email) GetCaptchaCacheKey(account string) string {
	return fmt.Sprintf("email_captcha::%s", account)
}

//func init() {
//	captcha.RegCaptcha(Name, &Email{})
//}
