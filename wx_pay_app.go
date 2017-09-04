package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

type WxAppOrderReq struct {
	Appid            string `xml:"appid"`
	Body             string `xml:"body"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Notify_url       string `xml:"notify_url"`
	Trade_type       string `xml:"trade_type"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
	Total_fee        int    `xml:"total_fee"`
	Out_trade_no     string `xml:"out_trade_no"`
	Sign             string `xml:"sign"`
}
type WxAppOrderResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
	Appid       string `xml:"appid"`
	Mch_id      string `xml:"mch_id"`
	Nonce_str   string `xml:"nonce_str"`
	Sign        string `xml:"sign"`
	Result_code string `xml:"result_code"`
	Prepay_id   string `xml:"prepay_id"`
	Trade_type  string `xml:"trade_type"`
}
type pay_wxpay_app struct {
	config map[string]interface{}
}

func (this *pay_wxpay_app) Init(conf map[string]interface{}) {
	this.config = conf

}

func (this *pay_wxpay_app) Create_Order(shop_name string, shop_body string, money int, order_no string) (map[string]interface{}, error) {
	var xmlReq WxAppOrderReq
	//fmt.Println(t.pay_index.config)
	fmt.Println(money)
	xmlReq.Appid = this.config["appid"].(string) //
	xmlReq.Mch_id = this.config["mchid"].(string)

	xmlReq.Body = shop_name

	xmlReq.Nonce_str = "12341234243"
	xmlReq.Notify_url = "pay.yaove.com/_callback_" + order_no + "-3"
	xmlReq.Trade_type = "APP"
	xmlReq.Spbill_create_ip = "103.37.160.115"
	xmlReq.Total_fee = money //单位是分
	xmlReq.Out_trade_no = order_no

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = xmlReq.Appid
	m["body"] = xmlReq.Body
	m["mch_id"] = xmlReq.Mch_id
	m["notify_url"] = xmlReq.Notify_url
	m["trade_type"] = xmlReq.Trade_type
	m["spbill_create_ip"] = xmlReq.Spbill_create_ip
	m["total_fee"] = xmlReq.Total_fee
	m["out_trade_no"] = xmlReq.Out_trade_no
	m["nonce_str"] = xmlReq.Nonce_str

	xmlReq.Sign = this.WxpayCalcSign(m, this.config["apikey"].(string)) //这个是计算wxpay签名的函数上面已贴出

	bytes_req, err := xml.Marshal(xmlReq)

	if err != nil {
		Log.Write("[微信APP支付]XMl编码错误,错误原因"+err.Error(), LogErr)

		return make(map[string]interface{}), errors.New("[微信APP支付]XMl编码错误,错误原因" + err.Error())
	}

	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "UnifyOrderReq", "xml", -1)
	bytes_req = []byte(str_req)
	fmt.Println(string(bytes_req))
	//发送unified order请求.
	req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", bytes.NewReader(bytes_req))
	if err != nil {
		Log.Write("[微信APP支付]Http Request失败,错误原因"+err.Error(), LogErr)
		return make(map[string]interface{}), errors.New("[微信APP支付]Http Request失败,错误原因" + err.Error())
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	c := http.Client{}
	resp, _err := c.Do(req)

	if _err != nil {
		Log.Write("[微信APP支付]请求微信支付统一下单接口发送错误,错误原因"+_err.Error(), LogErr)
		return make(map[string]interface{}), errors.New("[微信APP支付]请求微信支付统一下单接口发送错误,错误原因" + _err.Error())
	}
	bodyByte, _ := ioutil.ReadAll(resp.Body)
	body := bodyByte
	fmt.Println(string(body))
	var xmlResp WxAppOrderResp
	xml.Unmarshal(body, &xmlResp)
	//fmt.Println(string(body))
	if xmlResp.Return_code == "SUCCESS" {
		returnMap := make(map[string]interface{})
		fmt.Println(xmlResp.Prepay_id)
		returnMap["appid"] = this.config["appid"].(string)

		returnMap["timestamp"] = sys_intstr(int(time.Now().Unix()))
		returnMap["noncestr"] = sys_intstr(int(time.Now().Unix()))
		returnMap["prepayid"] = xmlResp.Prepay_id
		returnMap["partnerid"] = xmlResp.Mch_id
		returnMap["package"] = "Sign=WXPay"
		returnMap["paySign"] = this.WxpayCalcSign(returnMap, this.config["apikey"].(string))
		Log.Write("[微信APP支付] 订单"+order_no+"生成完毕", LogInfo)
		return returnMap, nil

		//
	}
	Log.Write("[微信APP支付]统一下单接口返回错误,错误原因"+xmlResp.Return_msg, LogErr)
	return make(map[string]interface{}), errors.New("[微信APP支付]统一下单接口返回错误,错误原因" + xmlResp.Return_msg)

	//defer resp.Body.Close()
}
func Refund_app() {

}

type WxAppNotifyReq struct {
	Return_code    string `xml:"return_code"`
	Return_msg     string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Nonce          string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Result_code    string `xml:"result_code"`
	Openid         string `xml:"openid"`
	Is_subscribe   string `xml:"is_subscribe"`
	Trade_type     string `xml:"trade_type"`
	Bank_type      string `xml:"bank_type"`
	Total_fee      int    `xml:"total_fee"`
	Fee_type       string `xml:"fee_type"`
	Cash_fee       int    `xml:"cash_fee"`
	Cash_fee_Type  string `xml:"cash_fee_type"`
	Transaction_id string `xml:"transaction_id"`
	Out_trade_no   string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	Time_end       string `xml:"time_end"`
}

type WxAppNotifyResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}

func (t *pay_wxpay_app) Callback(w http.ResponseWriter, r *http.Request) bool {
	// body
	body, read_err := ioutil.ReadAll(r.Body)
	if read_err != nil {
		fmt.Println("读取http body失败，原因!", read_err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}
	return_type := false
	var mr WxAppNotifyReq
	err := xml.Unmarshal(body, &mr)
	if err != nil {
		fmt.Println("解析HTTP Body格式到xml失败，原因!", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)

	reqMap["return_code"] = mr.Return_code
	reqMap["return_msg"] = mr.Return_msg
	reqMap["appid"] = mr.Appid
	reqMap["mch_id"] = mr.Mch_id
	reqMap["nonce_str"] = mr.Nonce
	reqMap["result_code"] = mr.Result_code
	reqMap["openid"] = mr.Openid
	reqMap["is_subscribe"] = mr.Is_subscribe
	reqMap["trade_type"] = mr.Trade_type
	reqMap["bank_type"] = mr.Bank_type
	reqMap["total_fee"] = mr.Total_fee
	reqMap["fee_type"] = mr.Fee_type
	reqMap["cash_fee"] = mr.Cash_fee
	reqMap["cash_fee_type"] = mr.Cash_fee_Type
	reqMap["transaction_id"] = mr.Transaction_id
	reqMap["out_trade_no"] = mr.Out_trade_no
	reqMap["attach"] = mr.Attach
	reqMap["time_end"] = mr.Time_end

	var resp WxAppNotifyResp
	//进行签名校验
	if t.WxpayVerifySign(reqMap, mr.Sign) {
		return_type = true
		//这里就可以更新我们的后台数据库了，其他业务逻辑同理。
		resp.Return_code = "SUCCESS"
		resp.Return_msg = "OK"
	} else {
		return_type = false
		resp.Return_code = "FAIL"
		resp.Return_msg = "failed to verify sign, please retry!"
	}

	//结果返回，微信要求如果成功需要返回return_code "SUCCESS"
	bytes, _err := xml.Marshal(resp)
	strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)
	if _err != nil {
		fmt.Println("xml编码失败，原因：", _err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	w.(http.ResponseWriter).WriteHeader(http.StatusOK)
	fmt.Fprint(w.(http.ResponseWriter), strResp)
	return return_type
}

func (this *pay_wxpay_app) WxpayVerifySign(needVerifyM map[string]interface{}, sign string) bool {
	signCalc := this.WxpayCalcSign(needVerifyM, this.config["apikey"].(string))
	if sign == signCalc {

		return true
	}
	return false
}
func (pay_wxpay_app) WxpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	//fmt.Println("微信支付签名计算, API KEY:", key)
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		//fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}

	}

	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}

	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}
