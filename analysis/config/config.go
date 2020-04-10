package config

import (
	"fmt"
	"time"
)

const (
	// ServiceName 是本微服务的名称
	ServiceName = "svc-analysis"

	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.himalife.srv"

	// DefaultRegisterTTL specifies how long a registration should exist in
	// discovery after which it expires and is removed
	DefaultRegisterTTL = 30 * time.Second

	// DefaultRegisterInterval is the time at which a service should re-register
	// to preserve it’s registration in service discovery.
	DefaultRegisterInterval = 15 * time.Second
)

// FullServiceName 返回微服务的全名
func FullServiceName() string {
	return fmt.Sprintf("%s.%s", ServiceNamespace, ServiceName)
}
