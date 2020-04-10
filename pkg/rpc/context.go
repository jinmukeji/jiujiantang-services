package rpc

import (
	"context"
	"net/http"

	"github.com/micro/go-micro/metadata"
)

const (
	// CidKey cid的key
	CidKey = "cid"
)

// ContextGetCid 从 Context 中获取 cid 的值
func ContextGetCid(ctx context.Context) string {
	cid := ""
	if md, ok := metadata.FromContext(ctx); ok {
		cid = md[http.CanonicalHeaderKey(CidKey)]
	}
	return cid
}
