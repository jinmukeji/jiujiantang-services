package rest

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
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

// 客户端授权
func (h *webHandler) ClientAuth(ctx iris.Context) {
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
	resp, errClientAuth := h.rpcSvc.ClientAuth(
		newRPCContext(ctx), req,
	)
	if errClientAuth != nil {
		writeRpcInternalError(ctx, errClientAuth, false)
		return
	}
	jwtToken, errBuildJwtToken, expiredAt := h.jwtMiddleware.BuildJwtToken(jwt.SigningMethodHS256, resp.ClientId, resp.Zone, resp.Name, resp.CustomizedCode)
	if errBuildJwtToken != nil {
		writeError(ctx, wrapError(ErrBuildJwtToken, "", fmt.Errorf("failed to build jwt token by client %s: %s", resp.ClientId, errBuildJwtToken.Error())), false)
	}
	rest.WriteOkJSON(ctx, ClientAuth{
		Authorization: jwtToken,
		ExpiredAt:     expiredAt.UTC(),
	})
}
