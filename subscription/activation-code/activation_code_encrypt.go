package activationcode

import (
	"encoding/base64"
	"fmt"

	"github.com/jinmukeji/go-pkg/crypto/encrypt/legacy"
)

// DefaultActivationCodeCipherHelper 默认激活码加密算法帮助类
type DefaultActivationCodeCipherHelper struct {
}

// ActivationCodeCipherHelper 激活码帮助类
type ActivationCodeCipherHelper interface {
	Encrypt(code string, key string, contractYear, maxUserLimits int32) string
	Decrypt(encryptedCode, key string, contractYear, maxUserLimits int32) string
}

// NewActivationCodeCipherHelper 新建激活码帮助类
func NewActivationCodeCipherHelper() ActivationCodeCipherHelper {
	return DefaultActivationCodeCipherHelper{}
}

// Encrypt 加密算法
func (client DefaultActivationCodeCipherHelper) Encrypt(code string, key string, contractYear, maxUserLimits int32) string {
	encryptedKey := fmt.Sprintf("%d%s%d", contractYear, key, maxUserLimits)
	return base64.StdEncoding.EncodeToString(legacy.AESEncrypt([]byte(code), []byte(encryptedKey)))
}

// Decrypt 解密算法
func (client DefaultActivationCodeCipherHelper) Decrypt(encryptedCode, key string, contractYear, maxUserLimits int32) string {
	encryptedPassword, _ := base64.StdEncoding.DecodeString(encryptedCode)
	encryptedKey := fmt.Sprintf("%d%s%d", contractYear, key, maxUserLimits)
	return string(legacy.AESDecrypt([]byte(encryptedPassword), []byte(encryptedKey)))
}
