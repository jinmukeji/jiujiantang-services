package blocker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// DefaultBlockerTestSuite 是 DefaultBlocker 的单元测试的 Test Suite
type DefaultBlockerTestSuite struct {
	suite.Suite
	blockerPool *BlockerPool
}

// 改进单元测试
const (
	blockerConfigPath = "./testdata/config_doc.yml"
	blockerDBFile     = "./data/GeoLite2-Country.mmdb.gz"
)

// SetupSuite 设置测试环境
func (suite *DefaultBlockerTestSuite) SetupSuite() {
	configDoc, err := LoadConfig(blockerConfigPath)
	if err != nil {
		panic(err)
	}
	suite.blockerPool, err = NewBlockerPool(configDoc, blockerDBFile)
	if err != nil {
		panic(err)
	}
}

// TestIsMacBlocked zone或者mac是否在白名单内
func (suite *DefaultBlockerTestSuite) TestIsMacBlocked() {
	t := suite.T()
	clientID := "jm-10002"
	// mac := "30451143FAEE"
	mac := "30451143D99A"
	zone := "CN"
	blocker := suite.blockerPool.GetBlocker(clientID)
	ok := blocker.IsMacBlocked(mac, zone)
	assert.Equal(t, false, ok)
}

// TestIsIPBlocked ip是否在白名单内
func (suite *DefaultBlockerTestSuite) TestIsIPBlocked() {
	t := suite.T()
	clientID := "jm-10002"
	blocker := suite.blockerPool.GetBlocker(clientID)

	ok := blocker.IsIPBlocked("114.236.8.103") // CN
	assert.Equal(t, false, ok)

	ok = blocker.IsIPBlocked("8.8.8.8") // US
	assert.Equal(t, true, ok)
}

// TestIgnoreIPCheck mac是否在免ip过滤白名单内
func (suite *DefaultBlockerTestSuite) TestIgnoreIPCheck() {
	t := suite.T()
	clientID := "jm-10004" 
	mac := "BBC123456789"
	blocker := suite.blockerPool.GetBlocker(clientID)
	ok := blocker.IgnoreIPCheck(mac)
	assert.Equal(t, true, ok)
}

func TestDefaultBlockerTestSuite(t *testing.T) {
	suite.Run(t, new(DefaultBlockerTestSuite))
}
