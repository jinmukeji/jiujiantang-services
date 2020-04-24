package rest

import (
	"io/ioutil"

	"fmt"

	yaml "gopkg.in/yaml.v2"
)

// Resource 资源
type Resource struct {
	Entry Env `json:"entry"`
}

// Env 环境
type Env map[string]Content

// Content 内容
type Content map[string]string

// LoadResourceFile 加载资源链接文件
func LoadResourceFile(filepath string) (*Resource, error) {
	res := Resource{}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from path %s: %s", filepath, err.Error())
	}
	if err = yaml.Unmarshal(content, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %s", content, err.Error())
	}
	return &res, nil
}
