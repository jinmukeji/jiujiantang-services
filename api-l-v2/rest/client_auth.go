package rest

import (
	"context"
	"net/http"
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/v2/metadata"
)

func init() {
	regErrorCode = regexp.MustCompile(`\[errcode:(\d+)\]`)
}

// TODO: 重构 REST 中的 JWT 授权验证与 AccessTokenKey 流程
const (
	// AccessTokenKey 用于从 context 中获取和设置用户登录凭证
	// nolint: gas
	AccessTokenKey       = "Access-Token"
	AccessTokenTypeKey   = "Access-Token-Type"
	AccessTokenTypeValue = "JinmuL"

	// ClientZoneKey 用于从 Context 的 Metadata中获取和设置Zone
	ClientZoneKey = "ClientZone"
	// ClientNameKey 用于从 Context 的 Metadata中获取和设置Name
	ClientNameKey = "ClientName"
	// ClientIDKey 用于从 Context 的 Metadata中获取和设置ClientID
	ClientIDKey = "ClientID"
	// ClientCustomizedCodeKey 用于从 Context 的 Metadata中获取和设置CustomizedCode
	ClientCustomizedCodeKey = "ClientCustomizedCode"
	// jwtDuration jwt过期时长
	jwtDuration = time.Hour * 12
	// RemoteClientIPKey 用于从 Context 的 Metadata中获取和设置Client的IP地址
	RemoteClientIPKey = "RemoteClientIP"
)

// ClientAuthReq 客户端授权的 Request
type ClientAuthReq struct {
	ClientID      string `json:"client_id"`
	SecretKeyHash string `json:"secret_key_hash"`
	Seed          string `json:"seed"`
}

// ClientAuth 授权
type ClientAuth struct {
	Authorization string    `json:"authorization"`
	ExpiredAt     time.Time `json:"expired_at"`
}

// newRPCContext 得到 RPC 的 Context
func newRPCContext(ctx iris.Context) context.Context {
	return metadata.NewContext(ctx.Request().Context(), map[string]string{
		// go 底层源码里面对 Key 传递的时候做了 CanonicalMIMEHeaderKey 处理
		http.CanonicalHeaderKey(AccessTokenKey):          ctx.GetHeader("X-Access-Token"),
		http.CanonicalHeaderKey(ClientIDKey):             ctx.Values().GetString(ClientIDKey),
		http.CanonicalHeaderKey(ClientZoneKey):           ctx.Values().GetString(ClientZoneKey),
		http.CanonicalHeaderKey(ClientCustomizedCodeKey): ctx.Values().GetString(ClientCustomizedCodeKey),
		http.CanonicalHeaderKey(RemoteClientIPKey):       ctx.Values().GetString(RemoteClientIPKey),
		http.CanonicalHeaderKey(ClientNameKey):           ctx.Values().GetString(ClientNameKey),
		http.CanonicalHeaderKey(AccessTokenTypeKey):      AccessTokenTypeValue,
	})
}

// 客户端授权
func (h *v2Handler) ClientAuth(ctx iris.Context) {
	var clientAuth ClientAuthReq
	err := ctx.ReadJSON(&clientAuth)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.ClientAuthRequest)
	req.ClientId = clientAuth.ClientID
	req.SecretKeyHash = clientAuth.SecretKeyHash
	req.Seed = clientAuth.Seed
	resp, err := h.rpcSvc.ClientAuth(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRpcInternalError(ctx, err, false)
		return
	}
	jwtToken, errBuildJwtToken, _ := h.jwtMiddleware.BuildJwtToken(jwt.SigningMethodHS256, resp.ClientId, resp.Zone, resp.Name, resp.CustomizedCode)
	if errBuildJwtToken != nil {
		writeError(ctx, wrapError(ErrBuildJwtToken, "", errBuildJwtToken), false)
	}
	rest.WriteOkJSON(ctx, ClientAuth{
		Authorization: jwtToken,
		ExpiredAt:     time.Now().Add(jwtDuration),
	})

}
