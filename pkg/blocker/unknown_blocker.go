package blocker

// UnknownBlocker 未知Blocker
type UnknownBlocker struct {
}

// NewUnknownBlocker 创建未知Blocker
func NewUnknownBlocker() *UnknownBlocker {
	b := &UnknownBlocker{}
	return b
}

// IsMacBlocked zone或者mac是否被限制
func (f *UnknownBlocker) IsMacBlocked(mac, zone string) bool {
	return true
}

// IgnoreIPCheck 是否忽略对mac对应客户端进行ip过滤
func (f *UnknownBlocker) IgnoreIPCheck(ip string) bool {
	return false
}

// IsIPBlocked ip是否被限制
func (f *UnknownBlocker) IsIPBlocked(ip string) bool {
	return true
}
