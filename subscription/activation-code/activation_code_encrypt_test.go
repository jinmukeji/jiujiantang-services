package activationcode

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type ActivationCodeEncrypltTestTestSuite struct {
	suite.Suite
}

// TestEncrypt测试加密算法
func (suite *ActivationCodeEncrypltTestTestSuite) TestEncrypt() {
	t := suite.T()
	h := NewActivationCodeCipherHelper()
	encryptedText := h.Encrypt("123456", "jinmu", 1, 200)
	suite.T().Log(encryptedText)
	assert.Equal(t, "123456", h.Decrypt(encryptedText, "jinmu", 1, 200))
}

//TestDecrypt测试解密算法
func (suite *ActivationCodeEncrypltTestTestSuite) TestDecrypt() {
	t := suite.T()
	h := NewActivationCodeCipherHelper()
	plainText := h.Decrypt("fLqIXQy7DkNU5S8xIpshFBy2WPQWlg==", "jinmu", 1, 200)
	suite.T().Log(plainText)
	assert.Equal(t, "123456", plainText)
}

func TestActivationCodeEncrypltTestTestSuite(t *testing.T) {
	suite.Run(t, new(ActivationCodeEncrypltTestTestSuite))
}
