package mysqldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// OptionsTestSuite 是 Options 的单元测试的 Test Suite
type OptionsTestSuite struct {
	suite.Suite
}

// TestDefaultOptions 测试 defaultOptions 方法成功返回 Options 记录
func (suite *OptionsTestSuite) TestDefaultOptions() {
	t := suite.T()
	opts := defaultOptions()

	assert.EqualValues(t, Options{
		Address:        "localhost:3306",
		EnableLog:      false,
		MaxConnections: 1,
		Charset:        "utf8mb4",
		ParseTime:      true,
		Locale:         "UTC",
	}, opts)
}

// TestNewOptions 测试 newOptions 方法成功返回 Options 记录
func (suite *OptionsTestSuite) TestNewOptions() {
	t := suite.T()
	opts := newOptions()

	assert.EqualValues(t, Options{
		Address:        "localhost:3306",
		EnableLog:      false,
		MaxConnections: 1,
		Charset:        "utf8mb4",
		ParseTime:      true,
		Locale:         "UTC",
	}, opts)
}

// TestNewOptions 测试带参的 newOptions 方法成功返回 Options 记录
func (suite *OptionsTestSuite) TestNewOptionsWithParameters() {
	t := suite.T()
	opts := newOptions(
		Address("0.0.0.0:6606"),
		EnableLog(true),
		MaxConnections(2),
		Charset("utf8"),
		ParseTime(false),
		Locale("en-US"),
	)

	assert.EqualValues(t, Options{
		Address:        "0.0.0.0:6606",
		EnableLog:      true,
		MaxConnections: 2,
		Charset:        "utf8",
		ParseTime:      false,
		Locale:         "en-US",
	}, opts)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOptionsTestSuite(t *testing.T) {
	suite.Run(t, new(OptionsTestSuite))
}
