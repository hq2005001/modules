package alipay

import (
	"errors"
	"fmt"
	"github.com/hq2005001/modules/payment"
	"github.com/hq2005001/modules/payment/config"
	"github.com/hq2005001/modules/payment/iap"
	"github.com/shopspring/decimal"
	"github.com/smartwalle/alipay/v3"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// Name 名称
	Name                 = "alipay"
	SubscribePay         = "ALIPAYAPP"
	SubscribeQrcodePay   = "QRCODE"
	NotifyTypeUserSign   = "dut_user_sign"
	NotifyTypeUserUnSign = "dut_user_unsign"
)

// Alipay 支付宝
type Alipay struct {
	client *alipay.Client
	config *config.AlipayConfig
	isApp  bool
}

// SetIsAPP 设置是否是App
func (a *Alipay) SetIsAPP(isAPP bool) payment.IPayment {
	a.isApp = isAPP
	return a
}

// Create 创建
func (a *Alipay) Create(params payment.CreatePaymentParam) payment.CreatePaymentResult {
	if params.SubPayType == SubscribePay || params.SubPayType == SubscribeQrcodePay {
		return a.Subscribe(params)
	}
	if a.isApp {
		return a.APPCreate(params)
	}
	if params.AgreementNO != "" {
		return a.AssistPay(params)
	}
	p := alipay.TradePreCreate{}
	p.NotifyURL = a.config.NotifyUrl
	p.ReturnURL = a.config.ReturnUrl
	p.OutTradeNo = params.PaymentID
	p.Subject = params.Title
	p.Body = params.Body
	p.TimeoutExpress = "30m"
	p.TotalAmount = params.Total.RoundDown(2).String()
	response, err := a.client.TradePreCreate(p)
	if err != nil {
		return payment.CreatePaymentResult{
			Status: false,
			ErrMsg: err.Error(),
		}
	}
	return payment.CreatePaymentResult{
		Status:        response.IsSuccess(),
		PaymentID:     response.OutTradeNo,
		QrCode:        response.QRCode,
		PaymentURL:    response.QRCode,
		PaymentParams: "",
		ErrMsg:        response.SubMsg,
	}
}

// Refund 退款
func (a *Alipay) Refund(id string, total, amount decimal.Decimal) interface{} {
	rs, err := a.client.TradeRefund(alipay.TradeRefund{
		OutTradeNo:   id,
		OutRequestNo: id,
		RefundAmount: amount.RoundDown(2).String(),
	})
	if err != nil {
		return false
	}
	if !rs.IsSuccess() {

	}
	return rs.IsSuccess()
}

// Verify 校验
func (a *Alipay) Verify(id string, thirdID []string, extraData interface{}, price int) (bool, *[]iap.Response) {

	rs, err := a.client.TradeQuery(alipay.TradeQuery{OutTradeNo: id})
	if err != nil {
		return false, nil
	}
	totalAmount, _ := strconv.ParseFloat(rs.TotalAmount, 64)
	amount, _ := strconv.Atoi(fmt.Sprintf("%.0f", float64(totalAmount*100)))
	data := make([]iap.Response, 0)
	data = append(data, iap.Response{
		PaymentID:           rs.OutTradeNo,
		OriginTransactionID: rs.TradeNo,
		Price:               amount,
		PayType:             Name,
	})
	return rs.Code == alipay.CodeSuccess && rs.TradeStatus == alipay.TradeStatusSuccess, &data
}

// VerifyNotify 验证通知
func (a *Alipay) VerifyNotify(req *http.Request) (bool, *[]iap.Response) {
	rs, err := a.client.GetTradeNotification(req)
	if err != nil || rs == nil {
		return false, nil
	}
	agreementNO := ""
	agreementType := ""
	productID := ""
	principalID := ""
	switch rs.NotifyType {
	case NotifyTypeUserUnSign:
		agreementType = "unsign"
		productID = rs.ExternalAgreementNo
		agreementNO = rs.AgreementNo
	case NotifyTypeUserSign:
		agreementNO = rs.AgreementNo
		agreementType = "sign"
		productID = rs.ExternalAgreementNo
		// 签约成功后取出签约信息
		agreementInfo, err := a.QueryAgreement(agreementNO)
		if err != nil {
			return false, nil
		}
		principalID = agreementInfo.PrincipalId
	}
	totalAmount, _ := strconv.ParseFloat(rs.TotalAmount, 64)
	amount, _ := strconv.Atoi(fmt.Sprintf("%.2f", float64(totalAmount*100)))
	data := make([]iap.Response, 0)
	isRefund := false
	if rs.NotifyType == alipay.NotifyTypeTradeStatusSync && rs.TradeStatus == alipay.TradeStatusClosed {
		isRefund = true
	}
	data = append(data, iap.Response{
		OriginTransactionID: rs.TradeNo,
		PaymentID:           rs.OutTradeNo,
		ProductID:           productID,
		Price:               amount,
		PayType:             Name,
		AgreementNO:         agreementNO,
		AgreementType:       agreementType,
		IsRefund:            &isRefund,
		PrincipalID:         principalID,
	})

	return true, &data
}

// Ack 确认
func (a *Alipay) Ack(writer http.ResponseWriter, isOK bool) {
	a.client.AckNotification(writer)
}

// IsIAP 是否是内购
func (a *Alipay) IsIAP() bool { return false }

// APPCreate app中创建订单
func (a *Alipay) APPCreate(params payment.CreatePaymentParam) payment.CreatePaymentResult {
	p := alipay.TradeAppPay{}
	p.NotifyURL = a.config.NotifyUrl
	p.ReturnURL = a.config.ReturnUrl
	p.OutTradeNo = params.PaymentID
	p.Subject = params.Title
	p.Body = params.Body
	p.TimeoutExpress = "2h"
	p.TotalAmount = fmt.Sprintf("%.2f", params.Total.InexactFloat64())
	appSign, err := a.client.TradeAppPay(p)
	if err != nil {
		return payment.CreatePaymentResult{}
	}
	return payment.CreatePaymentResult{
		Status:        true,
		PaymentID:     params.PaymentID,
		QrCode:        "",
		PaymentURL:    "",
		PaymentParams: appSign,
		ErrMsg:        "",
	}
}

// Sync 同步检查
func (a *Alipay) Sync(req *http.Request) (bool, *[]iap.Response) {

	return false, nil
}

// Subscribe 订阅
func (a *Alipay) Subscribe(params payment.CreatePaymentParam) payment.CreatePaymentResult {

	rs, err := a.client.AgreementPageSign(alipay.AgreementPageSign{
		NotifyURL: a.config.NotifyUrl,
		AccessParams: &alipay.AccessParams{
			Channel: params.SubPayType,
		},
		ProductCode:         "CYCLE_PAY_AUTH",
		PersonalProductCode: "CYCLE_PAY_AUTH_P",
		PeriodRuleParams: &alipay.PeriodRuleParams{
			PeriodType:   "DAY",
			Period:       strconv.Itoa(params.Duration),
			ExecuteTime:  time.Now().Format("2006-01-02"),
			SingleAmount: fmt.Sprintf("%.2f", float64(params.SubscribeMax)/100),
		},
		ExternalAgreementNo: fmt.Sprintf("%s:%s", params.MemberID, params.ProductID),
		//ExternalLogonId:     params.MemberID,
		SignScene: "INDUSTRY|APPSTORE",
	})
	if err != nil {
		return payment.CreatePaymentResult{ErrMsg: err.Error()}
	}
	qrcode := rs.String()
	if params.IsAPP {
		//qrcode = fmt.Sprintf("alipays://platformapi/startapp?appId=60000157&appClearTop=false&startMultApp=YES&sign_params=%s", url.QueryEscape(rs.RawQuery))
		qrcode = fmt.Sprintf("alipays://platformapi/startapp?appId=20000067&url=%s", url.QueryEscape(rs.String()))
	}

	return payment.CreatePaymentResult{
		Status:        true,
		QrCode:        qrcode,
		PaymentParams: qrcode,
		PaymentID:     params.PaymentID,
	}
}

// AssistPay 代扣
func (a *Alipay) AssistPay(params payment.CreatePaymentParam) payment.CreatePaymentResult {
	p := alipay.TradePay{}
	p.NotifyURL = a.config.NotifyUrl
	p.ReturnURL = a.config.ReturnUrl
	p.OutTradeNo = params.PaymentID
	p.Subject = params.Title
	p.Body = params.Body
	p.TimeoutExpress = "30m"
	p.TotalAmount = fmt.Sprintf("%.2f", params.Total.InexactFloat64())
	p.Scene = "deduct_pay"
	p.ProductCode = "CYCLE_PAY_AUTH"
	p.AuthCode = params.AgreementNO
	//a.client.Client.Timeout = time.Second * 10
	//log.New().Warn("取出签约信息 start request")
	response, err := a.client.TradePay(p)
	//log.New().Warn("取出签约信息after request ")
	if err != nil {
		return payment.CreatePaymentResult{}
	}
	return payment.CreatePaymentResult{
		Status:        response.IsSuccess(),
		PaymentID:     response.OutTradeNo,
		PaymentParams: "",
		ErrMsg:        response.SubMsg,
		ErrCode:       string(response.Code),
	}
}

// QueryAgreement 查询协议
func (a *Alipay) QueryAgreement(agreementNo string) (*alipay.AgreementQueryRsp, error) {
	rs, err := a.client.AgreementQuery(alipay.AgreementQuery{
		AgreementNo: agreementNo,
	})
	return rs, err
}

//// ModifyExecutionPlan 修改执行计划
//func (a *Alipay) ModifyExecutionPlan(agreementNo string, deductTime string) (*alipay.AgreementExecutionPlanRsp, error) {
//	return a.client.AgreementExecutionPlanModify(alipay.AgreementExecutionPlan{
//		AgreementNo: agreementNo,
//		DeductTime:  deductTime,
//	})
//}

func (a *Alipay) UnSign(agreementNo string) error {
	agreementInfo, err := a.QueryAgreement(agreementNo)
	if err != nil || agreementInfo.Code != "10000" {
		return errors.New("查询协议失败")
	}

	rs, err := a.client.AgreementUnsign(alipay.AgreementUnsign{
		PersonalProductCode: "CYCLE_PAY_AUTH_P",
		AgreementNo:         agreementNo,
		SignScene:           "INDUSTRY|APPSTORE",
		ExternalAgreementNo: agreementInfo.ExternalAgreementNo,
	})
	if err != nil {
		return err
	}
	if rs.Code == "" && rs.SubCode == "" {
		return nil
	}
	return err
}

func New(conf *config.AlipayConfig) *Alipay {
	client, err := alipay.New(conf.AppID, conf.PrivateKey, conf.Product, alipay.WithSandboxGateway("https://openapi-sandbox.dl.alipaydev.com/gateway.do"))
	if err != nil {
		panic("初始化支付宝失败")
	}
	if err = client.LoadAliPayPublicKey(conf.PublicKey); err != nil {
		panic("加载支付宝公钥失败")
	}
	return &Alipay{client: client, config: conf}
}
