package main

import (
	"fmt"
	//"net/url"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

type pay_alipay_app struct {
	config map[string]interface{}
}

func (this *pay_alipay_app) Init(conf map[string]interface{}) {
	this.config = conf

}

func (this *pay_alipay_app) Create_Order(shop_name string, shop_body string, money int, order_no string) (map[string]interface{}, error) {
	param := make(map[string]interface{})
	param["app_id"] = this.config["app_id"].(string)
	param["method"] = "alipay.trade.app.pay"
	param["sign_type"] = "RSA"
	param["charset"] = "utf-8"
	param["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	param["version"] = `1.0`
	param["notify_url"] = "http://pay.yaove.com/_callback_" + order_no + "-4"
	param["biz_content"] = `{"body":"` + shop_body + `"` + `,"subject":"` + shop_name + `","out_trade_no":"` + order_no + `","total_amount":` + strconv.FormatFloat(float64(money)/100, 'f', -1, 64) + `}`
	param["sign"] = this.Alipay_sortParamPrivate(param)
	if param["sign"] != nil {
		Log.Write("[支付宝APP支付] 订单"+order_no+"生成完毕", LogInfo)
		return param, nil
	}
	return make(map[string]interface{}), errors.New("[支付宝APP支付]签名失败")

	//defer resp.Body.Close()
}

func (this *pay_alipay_app) Alipay_sortParamPrivate(mReq map[string]interface{}) (sign string) {

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

		if value != "" {

			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	signStrings = sys_substr(signStrings, 0, -1)

	return url.QueryEscape(string(this.RsaEncryptPrivate(signStrings)))

}

func (this *pay_alipay_app) CheckPublic(mReq map[string]interface{}, sign string) bool {
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

		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	signStrings = sys_substr(signStrings, 0, -1)
	return this.RsaCheckPublic(signStrings, sign)

}
func (this *pay_alipay_app) Callback(w http.ResponseWriter, r *http.Request) bool {
	//fmt.Println("-----")
	dstFile, _ := os.Create("url.txt")
	defer dstFile.Close()
	dstFile.WriteString(r.Form.Encode())
	r.ParseForm()

	postMap := make(map[string]string)
	for k, v := range r.Form {

		if k == "call_param" {
			continue
		}
		if k == "sign" {
			continue
		}
		if k == "sign_type" {
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

	fmt.Println("sign" + r.Form["sign"][0])
	//fmt.Println(this.CheckPublic(md5Map, r.Form["sign"][0]))
	if this.CheckPublic(md5Map, r.Form["sign"][0]) {
		Log.Write("[支付宝APP支付]支付回调成功", LogInfo)
		return true
	}
	Log.Write("[支付宝APP支付]回调签名校验失败", LogWarning)
	return false
}
func (this *pay_alipay_app) RsaCheckPublic(str string, sign string) bool {
	block, _ := pem.Decode([]byte(this.config["publicKey"].(string)))
	if block == nil {
		Log.Write("[支付宝APP支付]回调私匙加载失败1", LogErr)
		return false
	}
	public, err := x509.ParsePKIXPublicKey(block.Bytes)
	rsaPub, _ := public.(*rsa.PublicKey)

	if err != nil {
		Log.Write("[支付宝APP支付]回调私匙加载失败2", LogErr)
		return false
	}

	h := sha1.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	data, _ := base64.StdEncoding.DecodeString(sign)

	err2 := rsa.VerifyPKCS1v15(rsaPub, crypto.SHA1, hashed, data)
	if err2 != nil {

		return false
	} else {
		return true
	}

}

func (t *pay_alipay_app) RsaEncryptPrivate(origData string) []byte {

	block, _ := pem.Decode([]byte(t.config["privateKey"].(string)))
	if block == nil {
		Log.Write("[支付宝APP支付]私匙加载失败1", LogErr)
		return nil
	}

	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		Log.Write("[支付宝APP支付]私匙加载失败2", LogErr)
		return nil
	}
	h := sha1.New()
	h.Write([]byte(origData))
	hashed := h.Sum(nil)
	re, err2 := rsa.SignPKCS1v15(nil, private, crypto.SHA1, hashed)
	if err2 != nil {
		Log.Write("[支付宝APP支付]SHA1签名失败", LogErr)
		return nil
	}
	data := base64.StdEncoding.EncodeToString(re)
	return []byte(data)
}

func (t *pay_alipay_app) Refund() {

}
func (t *pay_alipay_app) getSignVeryfy() {
}
