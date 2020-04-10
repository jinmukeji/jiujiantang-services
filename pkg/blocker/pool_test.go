package blocker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// NewBlockerPoolTestSuite 是 UnknownBlocker 的单元测试的 Test Suite
type NewBlockerPoolTestSuite struct {
	suite.Suite
	configDoc *ConfigDoc
}

// NewBlockerPoolTestSuite 设置测试环境
func (suite *NewBlockerPoolTestSuite) SetupSuite() {
	suite.configDoc, _ = LoadConfig(blockerConfigPath)
}

// TestIsMacBlocked zone或者mac是否在白名单内
func (suite *NewBlockerPoolTestSuite) TestNewBlockerPool() {
	t := suite.T()
	_, err := NewBlockerPool(suite.configDoc, blockerDBFile)
	assert.NoError(t, err)

}
func TestNewBlockerPoolTestSuite(t *testing.T) {
	suite.Run(t, new(NewBlockerPoolTestSuite))
}
