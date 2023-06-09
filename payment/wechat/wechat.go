package wechat

import (
	"context"
	"crypto/rsa"
	"github.com/hq2005001/modules/payment/config"
	jsoniter "github.com/json-iterator/go"

	"github.com/hq2005001/modules/payment"
	"github.com/hq2005001/modules/payment/iap"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/wechat"
)

const (
	// Name 名称
	Name = "wechat"
)

// Wechat 微信支付
type Wechat struct {
	//clientV3
	client     *wechat.Client
	conf       *config.WechatConfig
	ctx        context.Context
	privateKey *rsa.PrivateKey
	isAPP      bool
}

type AppParam struct {
	Appid     string `json:"appid"`
	Partnerid string `json:"partnerid"`
	Prepayid  string `json:"prepayid"`
	Noncestr  string `json:"noncestr"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

// SetIsAPP 设置是不是app
func (w *Wechat) SetIsAPP(isAPP bool) payment.IPayment {
	w.isAPP = isAPP
	return w
}

// Create 创建订单
func (w *Wechat) Create(params payment.CreatePaymentParam) payment.CreatePaymentResult {
	bm := make(gopay.BodyMap)
	noncestr := util.RandomString(32)
	bm.Set("nonce_str", noncestr).
		Set("out_trade_no", params.PaymentID).
		Set("total_fee", params.Total).
		Set("spbill_create_ip", params.ClientIP).
		Set("notify_url", w.conf.NotifyUrl).
		Set("trade_type", params.SubPayType).
		Set("sign_type", wechat.SignType_MD5).
		Set("limit_pay", "").
		Set("body", params.Body)
	if params.SubPayType == "JSAPI" {
		bm.Set("openid", params.OpenID)
	}

	//请求支付下单，成功后得到结果
	wxRsp, err := w.client.UnifiedOrder(context.Background(), bm)
	if err != nil {
		return payment.CreatePaymentResult{
			Status:    false,
			PaymentID: params.PaymentID,
		}
	}
	if wxRsp.ReturnCode != "SUCCESS" || wxRsp.ResultCode != "SUCCESS" {
		return payment.CreatePaymentResult{
			Status:    false,
			PaymentID: params.PaymentID,
		}
	}
	payParams := wxRsp.PrepayId
	if w.isAPP {
		timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
		sign := wechat.GetAppPaySign(w.conf.AppID, w.conf.MchID, wxRsp.NonceStr, wxRsp.PrepayId, wechat.SignType_MD5, timeStamp, w.conf.ApiKey)
		payParams, _ = jsoniter.MarshalToString(AppParam{
			Appid:     w.conf.AppID,
			Partnerid: w.conf.MchID,
			Prepayid:  wxRsp.PrepayId,
			Noncestr:  wxRsp.NonceStr,
			Timestamp: timeStamp,
			Sign:      sign,
		})
	}

	codeURL := wxRsp.CodeUrl
	if wxRsp.TradeType == "MWEB" {
		codeURL = wxRsp.MwebUrl
	}

	return payment.CreatePaymentResult{
		Status:        true,
		PaymentID:     params.PaymentID,
		QrCode:        wxRsp.CodeUrl,
		PaymentURL:    codeURL,
		PaymentParams: payParams,
		IsAsync:       true,
		RawData:       wxRsp,
	}
}

// Refund 退款
func (w *Wechat) Refund(id string, total, amount int64) interface{} {

	if err := w.client.AddCertPkcs12FilePath(w.conf.CertPem); err != nil {
		return false
	}
	bm := make(gopay.BodyMap)
	noncestr := util.RandomString(32)
	bm.Set("nonce_str", noncestr).
		Set("out_trade_no", id).
		Set("out_refund_no", id).
		Set("total_fee", total).
		Set("refund_fee", amount).
		Set("sign_type", wechat.SignType_MD5)

	////请求支付下单，成功后得到结果
	wxRsp, _, err := w.client.Refund(context.Background(), bm)
	if err != nil {
		return false
	}
	if wxRsp.ReturnCode != "SUCCESS" || wxRsp.ResultCode != "SUCCESS" {
		return false
	}
	return true
}

// Verify 验证
func (w *Wechat) Verify(id string, thirdID []string, extraData interface{}, price int) (bool, *[]iap.Response) {
	// 初始化参数结构体
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", id).
		Set("nonce_str", util.RandomString(32)).
		Set("sign_type", wechat.SignType_MD5)

	// 请求订单查询，成功后得到结果
	wxRsp, _, err := w.client.QueryOrder(context.Background(), bm)
	if err != nil {
		return false, nil
	}
	if wxRsp.ReturnCode == "SUCCESS" && wxRsp.ResultCode == "SUCCESS" && wxRsp.TradeState == "SUCCESS" {
		price, _ := strconv.Atoi(wxRsp.CashFee)
		return true, &[]iap.Response{
			{
				PaymentID:           id,
				Price:               price,
				OriginTransactionID: wxRsp.TransactionId,
			},
		}
	}
	return false, nil
}

// IsIAP 是否是内购
func (w *Wechat) IsIAP() bool {
	return false
}

// VerifyNotify 验证通知
func (w *Wechat) VerifyNotify(req *http.Request) (bool, *[]iap.Response) {
	notifyReq, err := wechat.ParseNotifyToBodyMap(req)
	if err != nil {
		return false, nil
	}
	ok, err := wechat.VerifySign(w.conf.ApiKey, wechat.SignType_MD5, notifyReq)
	if err != nil {
		return false, nil
	}
	if ok {
		resultCode := notifyReq.GetString("result_code")
		returnCode := notifyReq.GetString("return_code")
		if resultCode == "SUCCESS" && returnCode == "SUCCESS" {
			transactionID := notifyReq.GetString("transaction_id")
			priceStr := notifyReq.GetString("cash_fee")
			price, _ := strconv.Atoi(priceStr)
			paymentID := notifyReq.GetString("out_trade_no")
			return true, &[]iap.Response{
				{
					PaymentID:           paymentID,
					OriginTransactionID: transactionID,
					Price:               price,
				},
			}
		}

	}
	return false, nil
}

// Ack 确认
func (w *Wechat) Ack(writer http.ResponseWriter, isOK bool) {
	if isOK {
		rsp := new(wechat.NotifyResponse)
		rsp.ReturnCode = gopay.SUCCESS
		rsp.ReturnMsg = gopay.OK
		writer.Write([]byte(rsp.ToXmlString()))
		return
	}
	writer.WriteHeader(http.StatusInternalServerError)
}

// Sync 同步检查
func (w *Wechat) Sync(req *http.Request) (bool, *[]iap.Response) {
	return false, nil
}

func New(conf *config.WechatConfig) payment.IPayment {
	client := wechat.NewClient(conf.AppID, conf.MchID, conf.ApiKey, true)
	return &Wechat{
		client: client,
		conf:   conf,
		ctx:    context.TODO(),
	}
}
