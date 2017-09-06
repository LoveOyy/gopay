package AlipayApp

import (
	"net/http"
)

func (this *AlipayWeb) Callback(w http.ResponseWriter, r *http.Request) bool {
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

	}
	if this.CheckPublic(md5Map, r.Form["sign"][0]) {

		return true
	}
	return false
}
