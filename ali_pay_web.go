package main

import (
	"fmt"
	//"net/url"
	"net/http"
	"sort"
)

type pay_alipay_web struct {
	config map[string]interface{}
}

func (this *pay_alipay_web) Init(conf map[string]interface{}) {
	this.config = conf
}

func (this *pay_alipay_web) Create_Order(shop_name string, shop_body string, money int, order_no string) (map[string]interface{}, error) {
	param := make(map[string]interface{})
	param["service"] = "create_direct_pay_by_user"
	param["partner"] = this.config["partner"].(string)
	param["seller_email"] = this.config["seller_email"].(string)
	param["payment_type"] = "1"
	param["notify_url"] = "http://pay.4000968114.com/_callback_" + order_no + "-1"
	param["return_url"] = ""
	param["out_trade_no"] = order_no
	param["subject"] = shop_name
	param["total_fee"] = float64(money) / 100
	param["body"] = shop_body
	param["_input_charset"] = "utf-8"
	param["sign"] = this.Alipay_sortParam(param, this.config["key"].(string))
	Log.Write("[支付宝网页支付] 订单"+order_no+"生成完毕", LogInfo)
	return param, nil

}

func (pay_alipay_web) Alipay_sortParam(mReq map[string]interface{}, key string) (sign string) {

	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	var signStrings string
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", mReq[k])
		if k == "sign" {
			continue
		}
		if k == "sign_type" {
			continue
		}
		if value != "" {

			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	signStrings = sys_substr(signStrings, 0, -1)
	return (sys_md5(signStrings + key))
}
func (this *pay_alipay_web) Callback(w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()

	postMap := make(map[string]string)
	for k, v := range r.Form {
		if v[0] == "" {
			continue
		}
		if k == "call_param" {
			continue
		}
		postMap[k] = v[0]
	}
	md5Map := make(map[string]interface{})
	for k, v := range postMap {

		md5Map[k] = v
		fmt.Println(k)
		fmt.Println(v)
	}

	if this.Alipay_sortParam(md5Map, this.config["key"].(string)) != postMap["sign"] {
		Log.Write("[支付宝网页支付]支付回失败："+"签名错误!", LogErr)
		return false
	}
	Log.Write("[支付宝网页支付]支付回调成功", LogInfo)
	return true

}

func (t *pay_alipay_web) Refund() {

}
func (t *pay_alipay_web) getSignVeryfy() {

}
