package mysqldb

import (
	"context"
	"math/rand"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// DeleteSubscriptionUsersTestSuite 是 删除订阅下的用户 的 testSuite
type DeleteSubscriptionUsersTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *DeleteSubscriptionUsersTestSuite) SetupSuite() {
	envFilepath := filepath.Join("./testdata", "local.svc-subscription.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.db = db
}

// TestDeleteSubscriptionUsers 测试 删除订阅下的用户
func (suite *DeleteSubscriptionUsersTestSuite) TestDeleteSubscriptionUsers() {
	t := suite.T()
	var userID = rand.Int31n(100000000)
	var userIDList = []int32{userID}
	var subscriptionID = int32(77)
	ctx := context.Background()
	errCreateUserSubscriptionSharing := suite.db.GetDB(ctx).CreateUserSubscriptionSharing(ctx, &UserSubscriptionSharing{
		SubscriptionID: subscriptionID,
		UserID:         userID,
	})
	assert.NoError(t, errCreateUserSubscriptionSharing)
	errDeleteSubscriptionUsers := suite.db.GetDB(ctx).DeleteSubscriptionUsers(ctx, userIDList, subscriptionID)
	assert.NoError(t, errDeleteSubscriptionUsers)
}

func TestVerifyPhoneAndEmailTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteSubscriptionUsersTestSuite))
}
