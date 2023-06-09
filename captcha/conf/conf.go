package conf

type CustomEmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Config struct {
	Sms         SmsConfig         `mapstructure:"sms"`
	Email       EmailConfig       `mapstructure:"email"`
	CustomEmail CustomEmailConfig `mapstructure:"custom_email"`
}

type EmailConfig struct {
	AppKey        string `mapstructure:"app-key"`
	AppSecret     string `mapstructure:"app-secret"`
	Host          string `mapstructure:"host"`
	Username      string `mapstructure:"username"`
	UsernameAlias string `mapstructure:"username-alias"`
	Product       bool   `mapstructure:"product"`
}

type SmsConfig struct {
	GlobalMsg string `mapstructure:"global-msg"`
	CNMsg     string `mapstructure:"cn-msg"`
	Sign      string `mapstructure:"sign"`
	AppKey    string `mapstructure:"app-key"`
	AppSecret string `mapstructure:"app-secret"`
	Product   bool   `mapstructure:"product"`
}
