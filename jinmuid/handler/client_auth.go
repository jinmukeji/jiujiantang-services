package handler

import (
	"context"
	"errors"

	"fmt"

	"github.com/jinmukeji/go-pkg/crypto/hash"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// ClientAuth 客户端授权
func (j *JinmuIDService) ClientAuth(ctx context.Context, req *proto.ClientAuthRequest, resp *proto.ClientAuthResponse) error {
	client, err := j.datastore.FindClientByClientID(ctx, req.ClientId)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find client by clientID %s: %s", req.ClientId, err.Error()))
	}
	if req.SecretKeyHash != hash.HexString(hash.SHA256String(client.SecretKey+req.Seed)) {
		return NewError(ErrInvalidSecretKey, errors.New("secretkey not match"))
	}
	resp.Zone = client.Zone
	resp.Name = client.Name
	resp.CustomizedCode = client.CustomizedCode
	resp.ClientId = client.ClientID
	return nil
}
