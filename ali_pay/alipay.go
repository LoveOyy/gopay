package Alipay

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"sort"
	"strings"
)

type Alipay struct {
	config map[string]interface{}
}

func (this *Alipay) Init(conf map[string]interface{}) {
	this.config = conf
}

type Req struct {
	Param map[string]interface{}
}

func (this *Req) Set(name string, value interface{}) {
	if this.Param == nil {
		this.Param = make(map[string]interface{})
	}
	this.Param[name] = value
}
func (this *Req) Build(privateKey string) map[string]interface{} {
	str := SortParam(this.Param)

	this.Param["sign"] = RsaEncryptPrivate(str, this.Param["sign_type"].(string), privateKey)

	return this.Param
}

func BizContent(param map[string]interface{}) string {
	_json, _ := json.Marshal(param)
	return string(_json)
}

func RsaEncryptPrivate(strs string, signType string, privateKey string) string {
	block, _ := pem.Decode([]byte(privateKey))

	if block == nil {
		return ""
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}
	var err2 error
	var re []byte
	if signType == "RSA" {
		h := sha1.New()
		h.Write([]byte(strs))
		hashed := h.Sum(nil)
		re, err2 = rsa.SignPKCS1v15(nil, private, crypto.SHA1, hashed)
	} else {
		h := sha256.New()
		h.Write([]byte(strs))
		hashed := h.Sum(nil)
		re, err2 = rsa.SignPKCS1v15(nil, private, crypto.SHA256, hashed)
	}

	if err2 != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(re)

}

func SortParam(array map[string]interface{}) (sign string) {
	sorted_keys := make([]string, 0)
	for k, _ := range array {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	var signStrings string
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", array[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	signStrings = strings.Trim(signStrings, "&")

	return string(signStrings)

}
