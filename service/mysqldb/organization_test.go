package mysqldb

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	context "golang.org/x/net/context"
)

// OrganizationTestSuite 是 Organization 单元测试
type OrganizationTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 初始化测试
func (suite *OrganizationTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.db, _ = newTestingDbClientFromEnvFile(envFilepath)
}

// TestFindFirstOrganizationByOwner 测试查找指定 owner 拥有的第一个组织
func (suite *OrganizationTestSuite) TestFindFirstOrganizationByOwner() {
	const testUsername = "xx"
	t := suite.T()
	ctx := context.Background()
	var testOwnerID = 1
	o, err := suite.db.FindFirstOrganizationByOwner(ctx, testOwnerID)
	assert.NoError(t, err)
	assert.Equal(t, testUsername, o.Name)
}

// TestCreateOrganization 测试查找指定 owner 拥有的第一个组织
func (suite *OrganizationTestSuite) TestCreateOrganization() {
	t := suite.T()
	ctx := context.Background()
	organizationName := uuid.New().String()
	o := &Organization{
		Name:    organizationName,
		IsValid: 1,
	}
	assert.NoError(t, suite.db.CreateOrganization(ctx, o))
	now := time.Now()
	s := &Subscription{
		OrganizationID: o.OrganizationID,
		ActivatedAt:    now.UTC(),
		ExpiredAt:      now.UTC(),
		ContractYear:   0,
	}
	assert.NoError(t, suite.db.CreateSubscription(ctx, s))
}

// TestFindOrganizationsByOwner 测试查找指定 owner 拥有的第一个组织
func (suite *OrganizationTestSuite) TestFindOrganizationsByOwner() {
	const testUsername = "xx"
	t := suite.T()
	ctx := context.Background()
	var testOwnerID = 1
	os, err := suite.db.FindOrganizationsByOwner(ctx, testOwnerID)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(os))
	o := os[0]
	assert.Equal(t, testUsername, o.Name)
}

// TestFindOrganizationByID 测试从 id 查找一个组织
func (suite *OrganizationTestSuite) TestFindOrganizationByID() {
	t := suite.T()
	ctx := context.Background()
	const organizationID = 1
	const organizationName = "xx"
	o, err := suite.db.FindOrganizationByID(ctx, organizationID)
	assert.NoError(t, err)
	assert.Equal(t, organizationName, o.Name)
}

// TestDeleteOrganization 测试删除组织
func (suite *OrganizationTestSuite) TestDeleteOrganization() {
	const ownerID = 1
	t := suite.T()
	ctx := context.Background()
	organizationName := uuid.New().String()
	o := &Organization{
		Name:    organizationName,
		IsValid: 1,
	}
	assert.NoError(t, suite.db.CreateOrganization(ctx, o))
	now := time.Now()
	s := &Subscription{
		OrganizationID: o.OrganizationID,
		ActivatedAt:    now.UTC(),
		ExpiredAt:      now.UTC(),
		ContractYear:   0,
	}
	assert.NoError(t, suite.db.CreateSubscription(ctx, s))
	errCreateOrganizationOwner := suite.db.CreateOrganizationOwner(ctx, &OrganizationOwner{
		OrganizationID: o.OrganizationID,
		OwnerID:        ownerID,
		CreatedAt:      now.UTC(),
		UpdatedAt:      now.UTC(),
	})
	assert.NoError(t, errCreateOrganizationOwner)
	ok, err := suite.db.CheckOrganizationOwner(ctx, ownerID, o.OrganizationID)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.NoError(t, suite.db.DeleteOrganizationByID(ctx, o.OrganizationID))
	ok, err = suite.db.CheckOrganizationOwner(ctx, ownerID, o.OrganizationID)
	assert.NoError(t, err)
	assert.False(t, ok)
}

// TestFindUserSizeByOrganizationID 测试通过组织ID查看 user的数量
func (suite *OrganizationTestSuite) TestFindUserSizeByOrganizationID() {
	t := suite.T()
	ctx := context.Background()
	const organizationID = 2
	o, err := suite.db.GetExistingUserCountByOrganizationID(ctx, organizationID)
	assert.NoError(t, err)
	assert.Equal(t, organizationID, o)

}

// TestFindOrganizationSizeByOwnerID 测试通过OwnID查看组织的数量
func (suite *OrganizationTestSuite) TestFindOrganizationSizeByOwnerID() {
	t := suite.T()
	const ownerID = 1
	const size = 6 // 通过查看数据库该用户发现有6个组织
	ctx := context.Background()
	o, err := suite.db.GetOrganizationCountByOwnerID(ctx, ownerID)
	assert.NoError(t, err)
	assert.Equal(t, size, o)
}

// TestOrganizationTestSuite 启动测试
func TestOrganizationTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationTestSuite))
}
