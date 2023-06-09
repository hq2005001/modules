package payment

import (
	"github.com/hq2005001/modules/payment"
	"github.com/hq2005001/modules/payment/alipay"
	"github.com/hq2005001/modules/payment/config"
	"github.com/hq2005001/modules/payment/wechat"
)

// New 新建一个支付驱动
func New(conf *config.Config, paymentType payment.Type) payment.IPayment {
	switch paymentType {
	case alipay.Name:
		return alipay.New(conf.Alipay)
	case wechat.Name:
		return wechat.New(conf.Wechat)
	}
	return nil
}
