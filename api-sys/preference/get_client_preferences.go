package preference

import (
	"fmt"

	"github.com/blang/semver"
)

// GetClientPreferences 获取客户端对应资源配置
func (f ClientPreferences) GetClientPreferences(clientID, clientVersion, clientEnvironment string) (*ClientPreference, error) {
	version, err := semver.ParseTolerant(clientVersion)
	if err != nil {
		return nil, err
	}
	// 取版本号的前两位数字
	clientVersion = fmt.Sprintf("%d.%d", version.Major, version.Minor)

	if _, ok := f.mapClientConfig[clientID]; !ok {
		return nil, fmt.Errorf("clientID: %s ,clientVersion : %s , clientEnvironment : %s , clientID is invalid", clientID, clientVersion, clientEnvironment)
	}

	if _, ok := f.mapClientConfig[clientID][clientVersion]; !ok {
		return nil, fmt.Errorf("clientID: %s ,clientVersion : %s , clientEnvironment : %s , clientVersion is invalid", clientID, clientVersion, clientEnvironment)
	}

	if _, ok := f.mapClientConfig[clientID][clientVersion][clientEnvironment]; !ok {
		return nil, fmt.Errorf("clientID: %s ,clientVersion : %s , clientEnvironment : %s , clientEnvironment is invalid", clientID, clientVersion, clientEnvironment)
	}

	clientPreference := f.mapClientConfig[clientID][clientVersion][clientEnvironment]
	return &clientPreference, nil
}
