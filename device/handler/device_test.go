package handler

import (
	"context"
	"path/filepath"
	"testing"

	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	devicepb "github.com/jinmukeji/proto/gen/micro/idl/jm/device/v1"
	"github.com/micro/go-micro/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const rpcJinmuidServiceName = "com.xima.srv.svc-jinmuid"

type Account struct {
	account      string
	password     string
	userID       int32
	seed         string
	hashPassword string
	clientID     string
	deviceID     int32
}

// DeviceTestSuite 测试device
type DeviceTestSuite struct {
	suite.Suite
	deviceManagerService *DeviceManagerService
	account              *Account
	jinmuidSrv           jinmuidpb.UserManagerAPIService
}

// SetupSuite 初始化测试
func (suite *DeviceTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-device.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.account = newTestingAccountFromEnvFile(envFilepath)
	suite.deviceManagerService = NewDeviceManagerService(db)
	suite.jinmuidSrv = jinmuidpb.NewUserManagerAPIService(rpcJinmuidServiceName, client.DefaultClient)
}

// TestUserGetUsedDevices 测试UserGetUsedDevices
func (suite *DeviceTestSuite) TestUserGetUsedDevices() {
	t := suite.T()
	ctx := context.Background()
	ctx, err := mockSignin(ctx, suite.jinmuidSrv, suite.account.account, suite.account.hashPassword, suite.account.seed)
	assert.NoError(t, err)
	resp := new(devicepb.UserGetUsedDevicesResponse)
	err = suite.deviceManagerService.UserGetUsedDevices(ctx, &devicepb.UserGetUsedDevicesRequest{
		UserId: suite.account.userID,
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(resp.Devices))
}

// TestUserUseDevice 测试 UserUseDevice
func (suite *DeviceTestSuite) TestUserUseDevice() {
	t := suite.T()
	ctx := context.Background()
	ctx, err := mockSignin(ctx, suite.jinmuidSrv, suite.account.account, suite.account.hashPassword, suite.account.seed)
	assert.NoError(t, err)
	resp := new(devicepb.UserUseDeviceResponse)
	err = suite.deviceManagerService.UserUseDevice(ctx, &devicepb.UserUseDeviceRequest{
		UserId:   suite.account.userID,
		DeviceId: suite.account.deviceID,
		ClientId: suite.account.clientID,
	}, resp)
	assert.NoError(t, err)
}

func TestDeviceTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceTestSuite))
}

// mockSignin 模拟登录
func mockSignin(ctx context.Context, jinmuidSrv jinmuidpb.UserManagerAPIService, username string, passwordHash, seed string) (context.Context, error) {
	resp, err := jinmuidSrv.UserSignInByUsernamePassword(ctx, &jinmuidpb.UserSignInByUsernamePasswordRequest{
		Username:       username,
		HashedPassword: passwordHash,
		Seed:           seed,
	})
	if err != nil {
		return nil, err
	}
	return AddContextToken(ctx, resp.AccessToken), nil
}
