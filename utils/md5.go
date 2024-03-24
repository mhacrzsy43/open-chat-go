package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode((data)))
}

// 加密
func MakePassword(plainpwd, salt string) string {
	return MD5Encode(plainpwd + salt)
}

// 解密
func ValidPassword(plainpwd string, salt int32, password string) bool {
	// 将int32类型的salt转换为字符串
	saltStr := fmt.Sprintf("%06d", salt)

	// 拼接plainpwd和salt的字符串表示，然后进行MD5编码
	encodedPwd := MD5Encode(plainpwd + saltStr)
	fmt.Println("encodePwd:" + encodedPwd + "======saltStr:" + saltStr + "===========password:" + password)
	// 比较编码后的密码和提供的密码是否相等
	return encodedPwd == password
}
