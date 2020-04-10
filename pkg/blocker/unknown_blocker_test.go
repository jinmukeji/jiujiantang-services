package blocker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UnknownBlockerTestSuite 是 UnknownBlocker 的单元测试的 Test Suite
type UnknownBlockerTestSuite struct {
	suite.Suite
	blockerPool *BlockerPool
}

// SetupSuite 设置测试环境
func (suite *UnknownBlockerTestSuite) SetupSuite() {
	configDoc, _ := LoadConfig(blockerConfigPath)
	suite.blockerPool, _ = NewBlockerPool(configDoc, blockerDBFile)
}

// TestIsMacBlocked zone或者mac是否在白名单内
func (suite *UnknownBlockerTestSuite) TestIsMacBlocked() {
	t := suite.T()
	clientID := "A123456789"
	mac := "mac2"
	zone := "CN"
	blocker := suite.blockerPool.GetBlocker(clientID)
	ok := blocker.IsMacBlocked(mac, zone)
	assert.Equal(t, true, ok)
}

// TestIsIPBlocked ip是否在白名单内
func (suite *UnknownBlockerTestSuite) TestIsIPBlocked() {
	t := suite.T()
	ip := "121.0.0.10"
	clientID := "A123456789"
	blocker := suite.blockerPool.GetBlocker(clientID)
	ok := blocker.IsIPBlocked(ip)
	assert.Equal(t, true, ok)
}

// TestIgnoreIPCheck mac是否在免ip过滤白名单内
func (suite *UnknownBlockerTestSuite) TestIgnoreIPCheck() {
	t := suite.T()
	clientID := "A123456789"
	mac := "mac3"
	blocker := suite.blockerPool.GetBlocker(clientID)
	ok := blocker.IgnoreIPCheck(mac)
	assert.Equal(t, false, ok)
}

func TestUnknownBlockerTestSuite(t *testing.T) {
	suite.Run(t, new(UnknownBlockerTestSuite))
}
