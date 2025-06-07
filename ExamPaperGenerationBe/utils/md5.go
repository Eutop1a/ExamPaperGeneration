package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// CalcStringMD5 计算字符串的 MD5 值
func CalcStringMD5(input string) string {
	if input == "" {
		return ""
	}
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
