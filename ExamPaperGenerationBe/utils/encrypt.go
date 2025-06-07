package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// EncPassword 使用bcrypt包的函数直接加密密码
func EncPassword(password string) (string, error) {
	cipher, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(cipher), nil
}

// DecPassword 比较密码是否正确
func DecPassword(plainText, cipher string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(cipher), []byte(plainText))
	return err == nil
}
