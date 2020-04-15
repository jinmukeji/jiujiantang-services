package handler

import (
	"context"
	"net/http"

	"github.com/micro/go-micro/v2/metadata"
)

// AccessTokenKey 用于从 Context 的 Metadata 中获取和设置用户会话访问凭证
const AccessTokenKey = "Access-Token"

// ClientIDKey 用于从 Context 的 Metadata中获取ClientID
const ClientIDKey = "ClientID"

// TokenFromContext 从 Context 的 Metadata 获取 token
func TokenFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", false
	}
	token, ok := md[http.CanonicalHeaderKey(AccessTokenKey)]
	return token, ok
}

// ClientIDFromContext  从 Context 的 Metadata 获取 clientID
func ClientIDFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", false
	}
	clientID, ok := md[http.CanonicalHeaderKey(ClientIDKey)]
	return clientID, ok
}

// AddContextToken  把account放入 context 的 metadata
func AddContextToken(ctx context.Context, token string) context.Context {
	return metadata.NewContext(ctx, map[string]string{
		AccessTokenKey: token,
	})
}
