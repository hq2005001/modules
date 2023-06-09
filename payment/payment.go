package payment

import (
	"github.com/hq2005001/modules/payment/iap"
	"github.com/shopspring/decimal"
	"net/http"
)

// Type 支付类型
type Type string

// IPayment 支付接口
type IPayment interface {
	Create(params CreatePaymentParam) CreatePaymentResult

	Refund(id string, total, amount int64) interface{}

	Verify(id string, thirdID []string, extraData interface{}, price int) (bool, *[]iap.Response)
	// 是否是内购
	IsIAP() bool
	// SetIsAPP 设置是否是app
	SetIsAPP(isAPP bool) IPayment

	// 验证通知
	VerifyNotify(req *http.Request) (bool, *[]iap.Response)

	// Ack 确认
	Ack(writer http.ResponseWriter, isOK bool)

	Sync(req *http.Request) (bool, *[]iap.Response)
}

// CreatePaymentResult 创建订单的返回结果
type CreatePaymentResult struct {
	Status        bool        `json:"status"`         // 订单创建状态
	ThirdPayID    string      `json:"third_pay_id"`   // 三方支付订单
	PaymentID     string      `json:"payment_id"`     // 本地支付订单号
	ProductID     string      `json:"product_id"`     // 商品ID
	QrCode        string      `json:"qr_code"`        // 二维码链接
	PaymentURL    string      `json:"payment_url"`    // 支付链接
	PaymentParams string      `json:"payment_params"` // 参数
	ErrMsg        string      `json:"err_msg"`        // 错误信息
	ErrCode       string      `json:"err_code"`       // 错误码
	IsAsync       bool        `json:"is_async"`       // 是否为异步支付
	ProductName   string      `json:"product_name"`   // 商品名
	RawData       interface{} `json:"raw_data"`
	IsSubscribe   bool        `json:"is_subscribe"` // 是否是订阅
}

// CreatePaymentParam 创建支付请求参数
type CreatePaymentParam struct {
	PaymentID    string          `json:"payment_id" form:"payment_id"`
	Total        decimal.Decimal `json:"total" form:"total"`                 // 支付金额 单位(分)
	Title        string          `json:"title" form:"title"`                 // 显示的标题
	Body         string          `json:"body" form:"body"`                   // 显示的内容
	ClientIP     string          `json:"client_ip" form:"client_ip"`         // 客户端Ip
	IsAPP        bool            `json:"is_app" form:"is_app"`               // 是否是app
	SubPayType   string          `json:"sub_pay_type" form:"sub_pay_type"`   // 子支付类型
	OpenID       string          `json:"open_id" form:"open_id"`             // open_id
	Duration     int             `json:"duration" form:"duration"`           // 时长
	MemberID     string          `json:"member_id" form:"member_id"`         // 用户唯一ID
	AgreementNO  string          `json:"agreement_no" form:"agreement_no"`   // 协议号
	SubscribeMax int             `json:"subscribe_max" form:"subscribe_max"` // 订阅的金额上限
	ProductID    string          `json:"product_id" form:"product_id"`       // 商品ID
}
