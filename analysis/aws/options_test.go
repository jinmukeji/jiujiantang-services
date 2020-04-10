package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// OptionsTestSuite 测试 aws 连接配置
type OptionsTestSuite struct {
	suite.Suite
}

// SetupSuite 初始化测试
func (suite *OptionsTestSuite) SetupSuite() {

}

// TestOptionsSetup 测试配置 aws
func (suite *OptionsTestSuite) TestOptionsSetup() {
	t := suite.T()
	opts := newOptions(
		BucketName("name"),
		AccessKeyID("id"),
		SecretKey("key"),
		Region("us"),
		PulseTestRawDataEnvironmentS3KeyPrefix("testdata"),
	)
	assert.Equal(t, "name", opts.BucketName)
	assert.Equal(t, "id", opts.AccessKeyID)
	assert.Equal(t, "key", opts.SecretKey)
	assert.Equal(t, "us", opts.Region)
	assert.Equal(t, "testdata", opts.PulseTestRawDataEnvironmentS3KeyPrefix)
}

func TestOptionsTestSuite(t *testing.T) {
	suite.Run(t, new(OptionsTestSuite))
}
