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

// OrganizationOwnerTestSuite 是 OrganizationOwner 单元测试
type OrganizationOwnerTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 初始化测试
func (suite *OrganizationOwnerTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.db, _ = newTestingDbClientFromEnvFile(envFilepath)
}

// TestCreateOrganizationOwner 测试 owner 成为一个组织的 owner
func (suite *OrganizationOwnerTestSuite) TestCreateOrganizationOwner() {
	t := suite.T()
	ctx := context.Background()
	const ownerID = 1
	organizationName := uuid.New().String()
	o := &Organization{
		Name:    organizationName,
		IsValid: 1,
	}
	assert.NoError(t, suite.db.GetDB(ctx).CreateOrganization(ctx, o))
	now := time.Now()
	s := &Subscription{
		OrganizationID: o.OrganizationID,
		ActivatedAt:    now.UTC(),
		ExpiredAt:      now.UTC(),
		ContractYear:   0,
	}
	assert.NoError(t, suite.db.GetDB(ctx).CreateSubscription(ctx, s))
	now := time.Now()
	assert.NoError(t, suite.db.GetDB(ctx).CreateOrganizationOwner(ctx, &OrganizationOwner{
		OrganizationID: o.OrganizationID,
		OwnerID:        ownerID,
		CreatedAt:      now.UTC(),
		UpdatedAt:      now.UTC(),
	}))
}

// TestCheckOrganizationOwnerSuccess 测试检查用户是否为组织的拥有者成功
func (suite *OrganizationOwnerTestSuite) TestCheckOrganizationOwnerSuccess() {
	const userID = 1
	const organizationID = 1
	t := suite.T()
	ctx := context.Background()
	ok, err := suite.db.GetDB(ctx).CheckOrganizationOwner(ctx, userID, organizationID)
	assert.NoError(t, err)
	assert.True(t, ok)
}

// TestCheckOrganizationOwnerFail 测试检查用户是否为组织的拥有者失败
func (suite *OrganizationOwnerTestSuite) TestCheckOrganizationOwnerFail() {
	const userID = 2
	const organizationID = 1
	t := suite.T()
	ctx := context.Background()
	ok, err := suite.db.GetDB(ctx).CheckOrganizationOwner(ctx, userID, organizationID)
	assert.NoError(t, err)
	assert.False(t, ok)
}

// TestOrganizationOwnerTestSuite 启动测试
func TestOrganizationOwnerTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationOwnerTestSuite))
}
