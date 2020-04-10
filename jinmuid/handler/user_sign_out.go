package handler

import (
	"context"
	"errors"

	"fmt"

	"github.com/jinmukeji/gf-api2/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// UserSignOut 注销用户
func (j *JinmuIDService) UserSignOut(ctx context.Context, req *proto.UserSignOutRequest, resp *proto.UserSignOutResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrGetAccessTokenFailure, errors.New("signout fail: cannot find access token"))
	}
	userID, errFindUserIDByToken := j.datastore.FindUserIDByToken(ctx, token)
	if errFindUserIDByToken != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", errFindUserIDByToken.Error()))
	}
	// TODO: createAuditUserSignout，DeleteToken要在同一个事务
	errCreateAuditUserSignout := j.createAuditUserSignout(ctx, userID, req.Ip)
	if errCreateAuditUserSignout != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create audit user signout of user %d and ip %s: %s", userID, req.Ip, errCreateAuditUserSignout.Error()))
	}
	errDeleteToken := j.datastore.DeleteToken(ctx, token)
	if errDeleteToken != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to delete token: %s", errDeleteToken.Error()))
	}
	return nil
}

// createAuditUserSignout 创建登出审计记录
func (j *JinmuIDService) createAuditUserSignout(ctx context.Context, userID int32, ip string) error {
	clientID, _ := ClientIDFromContext(ctx)
	errCreateAuditUserSigninSignout := j.datastore.CreateAuditUserSigninSignout(ctx, &mysqldb.AuditUserSigninSignout{
		UserID:     userID,
		ClientID:   clientID,
		RecordType: mysqldb.SignoutRecordType,
		IP:         ip,
	})
	if errCreateAuditUserSigninSignout != nil {
		return NewError(ErrGetAccessTokenFailure, fmt.Errorf("failed to create audit user signin or signout: %s", errCreateAuditUserSigninSignout.Error()))
	}
	return nil
}
