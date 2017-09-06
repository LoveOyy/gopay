package AlipayApp

import (
	"errors"
	"strconv"
	"time"

	"../../ali_pay"
)

func (this *AlipayWeb) CreateOrder(shop_name string, shop_body string, money int, order_no string) (map[string]interface{}, error) {
	Req := new(Alipay.Req)
	Req.Set("app_id", this.config["app_id"].(string))
	Req.Set("sign_type", this.config["sign_type"].(string))
	Req.Set("method", "alipay.trade.page.pay")
	Req.Set("charset", "utf-8")
	Req.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	Req.Set("version", "1.0")
	Req.Set("notify_url", "http://www.baidu.com")
	Req.Set("biz_content", Alipay.BizContent(map[string]interface{}{"body": shop_body, "subject": shop_name, "out_trade_no": order_no, "product_code": "FAST_INSTANT_TRADE_PAY", "total_amount": strconv.FormatFloat(float64(money)/100, 'f', -1, 64)}))
	param := Req.Build(this.config["privateKey"].(string))

	if param["sign"] != nil {
		return param, nil
	}
	return make(map[string]interface{}), errors.New("[支付宝APP支付]签名失败")

}
