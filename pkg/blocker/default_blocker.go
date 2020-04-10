package blocker

import (
	"github.com/jinmukeji/go-pkg/stringset"
	"github.com/jpillora/ipfilter"
)

// DefaultBlocker 默认Blocker定义
type DefaultBlocker struct {
	clientID               string             // 客户端id
	allowedMacSet          stringset.Set      // 允许通过的mac
	IgnoreIPCheckingMacSet stringset.Set      // 忽略ip检查的mac地址
	allowedMacZoneSet      stringset.Set      // 允许通过的mac所在区域
	iPFilter               *ipfilter.IPFilter // ip过滤器
}

// NewDefaultBlocker 创建缺省过滤器
func NewDefaultBlocker(clientID string, allowedMacs []string, allowedMacZones []string, IgnoreIPCheckingMacs []string, iPFilter *ipfilter.IPFilter) *DefaultBlocker {
	allowedMacZoneSet := make(stringset.Set)
	allowedMacSet := make(stringset.Set)
	IgnoreIPCheckingMacSet := make(stringset.Set)

	for _, v := range allowedMacs {
		allowedMacSet.Put(v)
	}
	for _, v := range allowedMacZones {
		allowedMacZoneSet.Put(v)
	}
	for _, v := range IgnoreIPCheckingMacs {
		IgnoreIPCheckingMacSet.Put(v)
	}

	return &DefaultBlocker{
		clientID:               clientID,
		allowedMacSet:          allowedMacSet,
		IgnoreIPCheckingMacSet: IgnoreIPCheckingMacSet,
		allowedMacZoneSet:      allowedMacZoneSet,
		iPFilter:               iPFilter,
	}
}

// IsMacBlocked zone或者mac是否被限制
func (f *DefaultBlocker) IsMacBlocked(mac, zone string) bool {
	// zone是否在区域白名单内
	if f.allowedMacZoneSet.Exist(zone) {
		return false
	}
	// mac是否在mac白名单内
	if f.allowedMacSet.Exist(mac) {
		return false
	}
	return true
}

// IsIPBlocked ip是否被限制
func (f *DefaultBlocker) IsIPBlocked(ip string) bool {
	return !f.iPFilter.Allowed(ip)
}

// IgnoreIPCheck 是否忽略对mac对应客户端进行ip过滤
func (f *DefaultBlocker) IgnoreIPCheck(mac string) bool {
	// 是否不需要对mac对应客户端进行ip过滤
	return f.IgnoreIPCheckingMacSet.Exist(mac)
}
