package mysqldb

// DeviceTestSuite 是 device 单元测试
import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    context "golang.org/x/net/context"
)

type DeviceTestSuite struct {
	suite.Suite
	db *DbClient
}

// DeviceTestSuite 初始化测试
func (suite *DeviceTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.db, _ = newTestingDbClientFromEnvFile(envFilepath)
}

// TestGetOrganizationDeviceList 测试通过organizationID查询与Device的关联关系
func (suite *DeviceTestSuite) TestGetOrganizationDeviceList() {
	t := suite.T()
	ctx := context.Background()
	var organizationID int32 = 2
    num := 6
	deviceOrganizationBindingList, err := suite.db.GetOrganizationDeviceList(ctx, organizationID)
	assert.NoError(t, err)
	assert.Equal(t, num, len(deviceOrganizationBindingList))
}

func TestDeviceTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceTestSuite))
}
