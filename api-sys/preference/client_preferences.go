package preference

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Client 客户端ID到Version的映射
type Client map[string]Version

// Version 版本到Environment的映射
type Version map[string]Environment

// Environment 环境到ClientPreference的映射
type Environment map[string]ClientPreference

// ClientPreference 客户端的资源配置
type ClientPreference struct {
	ApiURL       string `yaml:"api_url"`
	AppLoginURL  string `yaml:"app_login_url"`
	AppEntryURL  string `yaml:"app_entry_url"`
	AppFaqURL    string `yaml:"app_faq_url"`
	AppReportURL string `yaml:"app_report_url"`
}

// ClientPreferences 客户端配置信息
type ClientPreferences struct {
	mapClientConfig Client
}

// NewClientPreferences 建立ClientPreferences
func NewClientPreferences(configFile string) ClientPreferences {
	data, _ := os.ReadFile(configFile)
	configDoc := Client{}
	// 读取配置文件
	err := yaml.Unmarshal(data, &configDoc)
	if err != nil {
		log.Fatal(err)
	}
	client := ClientPreferences{
		mapClientConfig: configDoc,
	}
	return client
}
