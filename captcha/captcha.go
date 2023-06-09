package captcha

import (
	"github.com/hq2005001/modules/captcha/conf"
	"github.com/hq2005001/modules/captcha/email"
	"github.com/hq2005001/modules/captcha/mobile"
	"github.com/hq2005001/modules/drivers/nsq"
	"github.com/hq2005001/modules/drivers/redis"
)

type ICaptcha interface {
	Send(ip, account, subject, content string) error
	Validate(account string, code string) bool
	IsValidate(account string, clear bool) bool
	GenerateContent(account, subject, content string) (interface{}, error)
	GetCaptchaCacheKey(account string) string
	GetValidateKey(account string) string
}

// New 新建
func New(name string, config *conf.Config, nsqConfig *nsq.NsqConfig, redisConf *redis.Conf) ICaptcha {
	switch name {
	case email.Name:
		return &email.Email{Config: config, QueueConfig: nsqConfig, RedisConf: redisConf}
	case mobile.Name:
		return &mobile.Mobile{Config: config, QueueConfig: nsqConfig, RedisConf: redisConf}
	default:
		return &mobile.Mobile{Config: config, QueueConfig: nsqConfig, RedisConf: redisConf}
	}
}
