package blocker

// Blocker 针对ip和mac过滤的接口定义
type Blocker interface {
	// zone或者mac是否被限制
	IsMacBlocked(mac, zone string) bool

	// 是否忽略对mac对应客户端进行ip过滤
	IgnoreIPCheck(mac string) bool

	// ip是否被限制
	IsIPBlocked(ip string) bool
}
