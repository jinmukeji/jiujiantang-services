package rest

import (
	"github.com/jinmukeji/gf-api2/pkg/rest"
	"github.com/kataras/iris/v12"
)

// BluetoothNamePrefix 蓝牙name前缀
type BluetoothNamePrefix struct {
	BluetoothNamePrefix []string `json:"localNamePrefixs"`
}

// GetBluetoothNamePrefixes 查找搜索蓝牙name前缀接口
func (h *v2Handler) GetBluetoothNamePrefixes(ctx iris.Context) {
	rest.WriteOkJSON(ctx, BluetoothNamePrefix{
		BluetoothNamePrefix: []string{"JinMu", "HJT", "KM"},
	})
}
