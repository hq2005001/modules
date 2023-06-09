package config

type Config struct {
	Alipay *AlipayConfig
	Wechat *WechatConfig
}

type AlipayConfig struct {
	AppID      string `mapstructure:"appid"`
	PublicKey  string `mapstructure:"public-key"`
	PrivateKey string `mapstructure:"private-key"`
	ReturnUrl  string `mapstructure:"return-url"`
	NotifyUrl  string `mapstructure:"notify-url"`
	Domain     string `mapstructure:"domain"`
	Product    bool   `mapstructure:"product"`
}

// WechatConfig 微信支付配置
type WechatConfig struct {
	AppID     string `mapstructure:"app_id"`
	ApiKey    string `mapstructure:"api_key"`
	MchID     string `mapstructure:"mch_id"`
	NotifyUrl string `mapstructure:"notify_url"`
	CertPem   string `mapstructure:"cert_pem"`
}
