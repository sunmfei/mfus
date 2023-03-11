package cryption

import (
	"golang.org/x/crypto/bcrypt"
)

// EncryptStr EncryptPassword 将密码加密，需要传入密码返回的是加密后的密码
func EncryptStr(str string) (string, error) {
	// 加密密码，使用 cryption 包当中的 GenerateFromPassword 方法，cryption.DefaultCost 代表使用默认加密成本
	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		// 如果有错误则返回异常，加密后的空字符串返回为空字符串，因为加密失败
		return str, err
	} else {
		// 返回加密后的密码和空异常
		return string(encryptPassword), nil
	}
}

// EqualsStr EqualsPassword 对比密码是否正确
func EqualsStr(str, encryptStr string) bool {
	// 使用 cryption 当中的 CompareHashAndPassword 对比密码是否正确，第一个参数为加密后的密码，第二个参数为未加密的密码
	err := bcrypt.CompareHashAndPassword([]byte(encryptStr), []byte(str))
	// 对比密码是否正确会返回一个异常，按照官方的说法是只要异常是 nil 就证明密码正确
	return err == nil
}
