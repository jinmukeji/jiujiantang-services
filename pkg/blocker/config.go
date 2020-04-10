package blocker

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// ConfigDoc 用于ip和mac过滤的条件
type ConfigDoc []ClientConfigDoc

// ClientConfigDoc 用于从配置文件加载过滤的配置文件
type ClientConfigDoc struct {
	ClientID             string   `yaml:"client_id"`               // 客户端id
	AllowedIPs           []string `yaml:"allowed_ips"`             // 允许通过的ip
	AllowedIPCountries   []string `yaml:"allowed_ip_countries"`    // 允许通过的ip所在国家,如果ip所在国家在此集合,则允许ip通过
	AllowedMacZones      []string `yaml:"allowed_mac_zones"`       // 允许通过的mac所在区域,如果mac所在区域在此集合,则允许mac通过
	AllowedMacs          []string `yaml:"allowed_macs"`            // 允许通过的mac
	IgnoreIPCheckingMacs []string `yaml:"ignore_ip_checking_macs"` // 忽略ip检查的mac
}

// LoadConfig 从某个路径中加载配置文件,blockerConfigFile是配置文件的位置，blockerDBConfigFile是ip数据库文件的位置
func LoadConfig(blockerConfigFile string) (*ConfigDoc, error) {
	data, _ := ioutil.ReadFile(blockerConfigFile)
	configDoc := ConfigDoc{}
	// 读取配置文件
	err := yaml.Unmarshal(data, &configDoc)
	if err != nil {
		return nil, err
	}
	return &configDoc, nil
}
