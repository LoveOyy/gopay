package AlipayApp

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"sort"
	"strings"
)

type AlipayApp struct {
	config map[string]interface{}
}

func (this *AlipayApp) Init(conf map[string]interface{}) {
	this.config = conf
}
func (this *AlipayApp) Alipay_sortParamPrivate(mReq map[string]interface{}) (sign string) {
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
	signStrings = strings.Trim(signStrings, "&")

	return string(this.RsaEncryptPrivate(signStrings))

}
func (this *AlipayApp) CheckPublic(mReq map[string]interface{}, sign string) bool {
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
	signStrings = strings.Trim(signStrings, "&")
	return this.RsaCheckPublic(signStrings, sign)
}

func (this *AlipayApp) RsaCheckPublic(str string, sign string) bool {
	block, _ := pem.Decode([]byte(this.config["publicKey"].(string)))
	if block == nil {
		//Log.Write("[支付宝APP支付]回调私匙加载失败1", LogErr)
		return false
	}
	public, err := x509.ParsePKIXPublicKey(block.Bytes)
	rsaPub, _ := public.(*rsa.PublicKey)
	if err != nil {
		//Log.Write("[支付宝APP支付]回调私匙加载失败2", LogErr)
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
func (t *AlipayApp) RsaEncryptPrivate(origData string) []byte {

	block, _ := pem.Decode([]byte(t.config["privateKey"].(string)))
	if block == nil {
		fmt.Println("[支付宝APP支付]私匙加载失败1")
		return nil
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("[支付宝APP支付]私匙加载失败2")
		return nil
	}
	h := sha1.New()
	h.Write([]byte(origData))
	hashed := h.Sum(nil)
	re, err2 := rsa.SignPKCS1v15(nil, private, crypto.SHA1, hashed)
	if err2 != nil {
		fmt.Println("[支付宝APP支付]SHA1签名失败")

		return nil
	}
	data := base64.StdEncoding.EncodeToString(re)
	return []byte(data)
}

func (t *AlipayApp) Refund() {

}
