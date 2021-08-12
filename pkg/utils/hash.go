/**
 * @Author: Anpw
 * @Description:
 * @File:  hash
 * @Version: 1.0.0
 * @Date: 2021/7/18 20:48
 */

package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// CompareHash 比对密码是否正确
func CompareHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// GeneratePassHash 密码加密
func GeneratePassHash(password string) (hash string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
