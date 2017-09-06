package Alipay

import (
	"errors"
	"strconv"
	"time"
)

func (this *Alipay) Refund(refund_amount int, order_no string, is_out_order bool) (map[string]interface{}, error) {
	_Req := new(Req)
	_Req.Set("app_id", this.config["app_id"].(string))
	_Req.Set("sign_type", this.config["sign_type"].(string))
	_Req.Set("method", "alipay.trade.page.refund")
	_Req.Set("charset", "utf-8")
	_Req.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	_Req.Set("version", "1.0")
	bizContent := make(map[string]interface{})
	if is_out_order {
		bizContent["out_trade_no"] = order_no
	} else {
		bizContent["trade_no"] = order_no
	}
	bizContent["refund_amount"] = strconv.FormatFloat(float64(refund_amount)/100, 'f', -1, 64)
	_Req.Set("biz_content", BizContent(bizContent))
	param := _Req.Build(this.config["privateKey"].(string))

	if param["sign"] != nil {
		return param, nil
	}
	return make(map[string]interface{}), errors.New("[支付宝APP支付]签名失败")

}
