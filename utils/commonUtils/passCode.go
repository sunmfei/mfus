package commonUtils

import (
	"math/rand"
	"regexp"
	"time"
)

var MAX = 90
var MIN = 0

func GenPass(LENGTH int64) []byte {
	startChar := "!"
	var code []byte
	var i int64
	for i = 0; i < LENGTH; i++ {
		anInt := random(MIN, MAX)
		newChar := startChar[0] + byte(anInt)
		if newChar == ' ' {
			i = i - i
			continue
		}
		code = append(code, newChar)
	}
	return code
}
func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// GetStrongPasswordString 随机生成指定位数的大写字母和数字的组合
func GetStrongPasswordString(l int) string {
	//~!@#$%^&*()_+{}":?><;.,
	str := "123456789ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz!@#$%&*"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	ok1, _ := regexp.MatchString(".[1|2|3|4|5|6|7|8|9]", string(result))
	ok2, _ := regexp.MatchString(".[Z|X|C|V|B|N|M|A|S|D|F|G|H|J|K|L|Q|W|E|R|T|Y|U|I|P]", string(result))
	ok3, _ := regexp.MatchString(".[z|x|c|v|b|n|m|a|s|d|f|g|h|j|k|l|q|w|e|r|t|y|u|i|p]", string(result))
	ok4, _ := regexp.MatchString(".[!|@|#|$|%|&|*]", string(result))
	if ok1 && ok2 && ok3 && ok4 {
		return string(result)
	} else {
		return GetStrongPasswordString(l)
	}

}

// GetRandCode 过去随机码只包含数字字母
func GetRandCode(l int) []byte {
	str := "123456789ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return result
}
