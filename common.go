package main

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

func sys_intstr(tint int) string {
	return strconv.Itoa(tint)
}
func sys_substr(str string, start int, length int) string {

	rs := []rune(str)
	if start < 0 {
		start = len(rs) + start

	}
	if length == 0 {
		return string(rs[start:])
	} else if length < 0 {
		return string(rs[start : start+len(rs)+length])
	} else {
		return string(rs[start : start+length])
	}

}

func sys_md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
