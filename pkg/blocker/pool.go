package blocker

import (
	"github.com/jpillora/ipfilter"
)

// BlockerPool 过滤池
type BlockerPool map[string]Blocker

const (
	// 未知的客户端名称
	unknownClient = "unknownClient"
)

// NewBlockerPool 初始化配置Pool
func NewBlockerPool(configDoc *ConfigDoc, blockerDBConfigFile string) (*BlockerPool, error) {

	// 构造过滤映射关系
	blockers := make(BlockerPool)
	// 初始化
	for _, item := range *configDoc {

		// 初始化ip过滤选项
		opt := ipfilter.Options{
			AllowedIPs:       item.AllowedIPs,
			AllowedCountries: item.AllowedIPCountries,
			// IPDBPath:         blockerDBConfigFile,
			BlockByDefault: true,
		}
		ipFilter := ipfilter.New(opt)

		// 创建默认Blocker
		defaultBlocker := NewDefaultBlocker(item.ClientID, item.AllowedMacs, item.AllowedMacZones, item.IgnoreIPCheckingMacs, ipFilter)
		blockers[item.ClientID] = defaultBlocker
	}
	// 添加未知Blocker
	unknownBlocker := NewUnknownBlocker()
	blockers[unknownClient] = unknownBlocker
	return &blockers, nil
}

// GetBlocker 根据client_id获取对应的Blocker
func (p *BlockerPool) GetBlocker(clientID string) Blocker {
	if blocker, ok := (*p)[clientID]; ok {
		return blocker
	}
	// 找不到对应的Blocker时返回未知客户端对应的Blocker
	return (*p)[unknownClient]
}
