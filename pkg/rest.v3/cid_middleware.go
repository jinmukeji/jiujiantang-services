package rest

import (
	id "github.com/jinmukeji/go-pkg/id-gen"
	"github.com/kataras/iris/v12"
)

const (
	// ContextCidKey 上下文中注入的 cid 的键
	ContextCidKey = "cid"
)

// CidMiddleware cid中间件
func CidMiddleware(ctx iris.Context) {
	cid := id.NewXid()
	ctx.Values().Set(ContextCidKey, cid)
	ctx.Next()
}

// GetCidFromContext 得到cid从context
func GetCidFromContext(ctx iris.Context) string {
	return ctx.Values().GetString(ContextCidKey)
}
