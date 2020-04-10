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

// OrganizationUserTestSuite 是 OrganizationUser 单元测试
type OrganizationUserTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 初始化测试
func (suite *OrganizationUserTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.db, _ = newTestingDbClientFromEnvFile(envFilepath)
}

// TestCreateOrganizationUsers 测试再 organization_user 新增多条记录
func (suite *OrganizationUserTestSuite) TestCreateOrganizationUser() {
	t := suite.T()
	ctx := context.Background()
	users := make([]*OrganizationUser, 0)
	// 创建一个组织
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
	// 加入组织
	now := time.Now()
	for userID := 1; userID < 5; userID++ {
		users = append(users, &OrganizationUser{
			OrganizationID: o.OrganizationID,
			UserID:         userID,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}
	assert.NoError(t, suite.db.CreateOrganizationUsers(ctx, users))
	ok, err := suite.db.CheckOrganizationUser(ctx, users[0].UserID, o.OrganizationID)
	assert.NoError(t, err)
	assert.True(t, ok)
}

// TestDeleteOrganizationUser 测试再 organization_user 删除一条记录
func (suite *OrganizationUserTestSuite) TestDeleteOrganizationUser() {
	t := suite.T()
	ctx := context.Background()
	const userID = 1
	// 创建一个组织
	organizationName := uuid.New().String()
	o := &Organization{
		Name:    organizationName,
		IsValid: 1,
	}
	now := time.Now()
	assert.NoError(t, suite.db.CreateOrganization(ctx, o))
	s := &Subscription{
		OrganizationID: o.OrganizationID,
		ActivatedAt:    now.UTC(),
		ExpiredAt:      now.UTC(),
		ContractYear:   0,
	}
	assert.NoError(t, suite.db.CreateSubscription(ctx, s))
	// 加入组织
	users := make([]*OrganizationUser, 0)
	now = time.Now()
	users = append(users, &OrganizationUser{
		OrganizationID: o.OrganizationID,
		UserID:         userID,
		CreatedAt:      now.UTC(),
		UpdatedAt:      now.UTC(),
	})
	assert.NoError(t, suite.db.CreateOrganizationUsers(ctx, users))
	// 从组织删除该用户
	errDeleteOrganizationUser := suite.db.DeleteOrganizationUser(ctx, userID, o.OrganizationID)
	assert.NoError(t, errDeleteOrganizationUser)
	ok, err := suite.db.CheckOrganizationUser(ctx, userID, o.OrganizationID)
	assert.NoError(t, err)
	assert.False(t, ok)
}

// TestDeleteOrganizationUsers 测试在 organization_user 删除多条记录
func (suite *OrganizationUserTestSuite) TestDeleteOrganizationUsers() {
	t := suite.T()
	ctx := context.Background()
	// 创建一个组织
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
	// 加入组织
	now := time.Now()
	uids := make([]int32, 0)
	users := make([]*OrganizationUser, 0)
	for userID := 1; userID < 5; userID++ {
		users = append(users, &OrganizationUser{
			OrganizationID: o.OrganizationID,
			UserID:         userID,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
		uids = append(uids, int32(userID))
	}
	assert.NoError(t, suite.db.CreateOrganizationUsers(ctx, users))
	// 从组织删除该用户
	errDeleteOrganizationUsers := suite.db.DeleteOrganizationUsers(ctx, uids, int32(o.OrganizationID))
	assert.NoError(t, errDeleteOrganizationUsers)
	for _, uid := range uids {
		ok, err := suite.db.CheckOrganizationUser(ctx, int(uid), o.OrganizationID)
		assert.NoError(t, err)
		assert.False(t, ok)
	}
}

// TestFindOrganizationUsers 测试在 organization_user 查找用户
func (suite *OrganizationUserTestSuite) TestFindOrganizationUsers() {
	t := suite.T()
	ctx := context.Background()
	const organizationID = 1
	const userID = 1
	const username = "xx"
	// 加入组织
	now := time.Now()
	errCreateOrganizationUsers := suite.db.CreateOrganizationUsers(ctx, []*OrganizationUser{&OrganizationUser{
		OrganizationID: organizationID,
		UserID:         userID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}})
	assert.NoError(t, errCreateOrganizationUsers)
	users, err := suite.db.FindOrganizationUsers(ctx, organizationID)
	assert.NoError(t, err)
	assert.Equal(t, username, users[0].Username)
}

// TestOrganizationUserTestSuite 启动测试
func TestOrganizationUserTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationUserTestSuite))
}
