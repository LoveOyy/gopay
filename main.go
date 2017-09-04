package main

import (
	"fmt"
)

func main() {
	test_wx_pay_app()
	test_wx_pay_native()
	test_ali_pay_web()
	test_ali_pay_app()
}
func test_wx_pay_native() {
	var _pay = new(pay_wxpay_Native)
	_conf := make(map[string]interface{})
	_conf["appid"] = "wx43ba6b8465fb36d4"
	_conf["mchid"] = "1274622101"
	_conf["apikey"] = "dianying654321dianying654321cmcm"
	_pay.Init(_conf)
	_pay.Create_Order("1", "2", 10, "1112545")
}
func test_wx_pay_app() {
	var _pay = new(pay_wxpay_app)
	_conf := make(map[string]interface{})
	_conf["appid"] = "wx43ba6b8465fb36d4"
	_conf["mchid"] = "1274622101"
	_conf["apikey"] = "dianying654321dianying654321cmcm"
	_pay.Init(_conf)
	_pay.Create_Order("1", "2", 10, "1112545")
}
func test_ali_pay_web() {
	var _pay = new(pay_alipay_web)
	_conf := make(map[string]interface{})
	_conf["partner"] = "2088621448590388"
	_conf["seller_email"] = "2088621448590388"
	_conf["key"] = "lred7dl6ka6rjawfkdba25jwij1r17k1"
	_pay.Init(_conf)
	fmt.Println(_pay.Create_Order("1", "2", 10, "1112545"))
}
func test_ali_pay_app() {
	var _pay = new(pay_alipay_app)
	_conf := make(map[string]interface{})
	_conf["app_id"] = "2088621448590388"
	_conf["privateKey"] = ""
	_conf["publicKey"] = ""

	_pay.Init(_conf)
	fmt.Println(_pay.Create_Order("1", "2", 10, "1112545"))
}
