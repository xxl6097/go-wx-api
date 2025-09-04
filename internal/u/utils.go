package u

import (
	"crypto/sha1"
	"fmt"
	"runtime"
	"sort"
	"strings"
)

func IsMacOs() bool {
	if strings.Compare(runtime.GOOS, "darwin") == 0 {
		return true
	}
	return false
}

// CheckSignature 验证微信请求签名
// 参数: signature, timestamp, nonce 来自请求URL, token 为你在微信后台配置的令牌
// 返回值: 验证通过返回 true，否则返回 false
func CheckSignature(signature, timestamp, nonce, token string) bool {
	// 1. 将 token, timestamp, nonce 放入切片
	params := []string{token, timestamp, nonce}
	// 2. 按字典序排序
	sort.Strings(params)
	// 3. 拼接成一个字符串
	var combinedStr string
	for _, s := range params {
		combinedStr += s
	}
	// 4. 对拼接后的字符串进行 sha1 加密
	hasher := sha1.New()
	hasher.Write([]byte(combinedStr))
	calculatedSignature := fmt.Sprintf("%x", hasher.Sum(nil)) // %x 表示格式化为小写十六进制
	// 5. 将加密后的字符串与 signature 对比
	return calculatedSignature == signature
}
