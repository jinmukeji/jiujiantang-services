package auth

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
)

// GenerateToken 是用户凭证的生成函数
func GenerateToken() string {
	return uuid.New().String()
}

// contextKey 用于获取上下文环境
type contextKey string

// AccessTokenKey 用于从 Context 的 Metadata 中获取和设置用户会话访问凭证
const AccessTokenKey = "Access-Token"

// ZoneKey 用于从 Context 的 Metadata中获取Zone
const ZoneKey = "ClientZone"

// NameKey 用于从 Context 的 Metadata中获取Name
const NameKey = "ClientName"

// ClientIDKey 用于从 Context 的 Metadata中获取ClientID
const ClientIDKey = "ClientID"

// CustomizedCodeKey 用于从 Context 的 Metadata中获取CustomizedCode
const CustomizedCodeKey = "ClientCustomizedCode"

// RemoteClientIPKey 用于从 Context 中的 Metadata获取RemoteClientIP
const RemoteClientIPKey = "RemoteClientIP"

// AccessTokenType AccessToken类型
const AccessTokenType = "Access-Token-Type"

var (
	// userIDKey 用于从 context 中获取和设置用户ID
	userIDKey contextKey = "UserID"
	// accountKey 用于从 context 中获取和设置账户
	accountKey contextKey = "Account"
	// MachineUUID 机器UUID
	machineUUID contextKey = "MachineUUID"
)

// TokenFromContext 从 Context 的 Metadata 获取 token
func TokenFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", false
	}
	token, ok := md[http.CanonicalHeaderKey(AccessTokenKey)]
	return token, ok
}

// AccessTokenTypeFromContext 从 Context 的 Metadata 获取 tokenType
func AccessTokenTypeFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", false
	}
	accessTokenType, ok := md[http.CanonicalHeaderKey(AccessTokenType)]
	return accessTokenType, ok
}

// ZoneFromContext 从元数据获取 zone
func ZoneFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", false
	}
	zone, ok := md[http.CanonicalHeaderKey(ZoneKey)]
	return zone, ok
}

// NameFromContext 从context中获取用户名
func NameFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", false
	}
	name, ok := md[http.CanonicalHeaderKey(NameKey)]
	return name, ok
}

// ClientIDFromContext 从context中获取ClientID
func ClientIDFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", false
	}
	clientID, ok := md[http.CanonicalHeaderKey(ClientIDKey)]
	return clientID, ok
}

// CustomizedCodeFromContext 从context中获取CustomizedCode
func CustomizedCodeFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", false
	}
	customizedCode, ok := md[http.CanonicalHeaderKey(CustomizedCodeKey)]
	return customizedCode, ok
}

// RemoteClientIPFromContext 从context中获取远程客户端IP
func RemoteClientIPFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", false
	}
	ip, ok := md[http.CanonicalHeaderKey(RemoteClientIPKey)]
	return ip, ok
}

// UserIDFromContext 从 context 获取 userID
func UserIDFromContext(ctx context.Context) (int32, bool) {
	userID, ok := ctx.Value(userIDKey).(int32)
	return userID, ok
}

// AccountFromContext 从 context 获取 account
func AccountFromContext(ctx context.Context) (string, bool) {
	account, ok := ctx.Value(accountKey).(string)
	return account, ok
}

// AddContextUserID 把 userID 放入 context
func AddContextUserID(ctx context.Context, userID int32) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// AddContextAccount 把 account 放入 context
func AddContextAccount(ctx context.Context, account string) context.Context {
	return context.WithValue(ctx, accountKey, account)
}

// AddContextToken  把account放入 context 的 metadata
func AddContextToken(ctx context.Context, token string) context.Context {
	return metadata.NewContext(ctx, map[string]string{
		AccessTokenKey: token,
	})
}

// AddContextMachineUUID 把 MachineUUID 放入 context
func AddContextMachineUUID(ctx context.Context, uuid string) context.Context {
	return context.WithValue(ctx, machineUUID, uuid)
}

// MachineUUIDFromContext 从 context 获取 uuid
func MachineUUIDFromContext(ctx context.Context) (string, bool) {
	uuid, ok := ctx.Value(machineUUID).(string)
	return uuid, ok
}
