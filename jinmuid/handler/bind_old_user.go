package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinmukeji/go-pkg/crypto/encrypt/legacy"
	"github.com/jinmukeji/go-pkg/crypto/hash"
	. "github.com/jinmukeji/go-pkg/with"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// BindOldUser 绑定老用户
func (j *JinmuIDService) BindOldUser(ctx context.Context, req *proto.BindOldUserRequest, resp *proto.BindOldUserResponse) error {
	// 验证 token 里面的的新用户的 UserID 和传入的新用户的 UserID 是否一致
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get userID by token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
	}

	// 验证新用户已经设置手机号
	newUser, err := j.datastore.FindUserByUserID(ctx, userID)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by userID %d: %s", userID, err.Error()))
	}
	if !newUser.HasSetPhone {
		return NewError(ErrSignPhoneNotSet, fmt.Errorf("signin phone has not been set for userId %d", req.UserId))
	}

	// 验证输入的老用户的用户名密码是否正确
	oldUser, err := j.datastore.FindUserByUsername(ctx, req.Username)
	if oldUser == nil || err != nil {
		return NewError(ErrNonexistentUsername, fmt.Errorf("failed to find user by username %s", req.Username))
	}
	if !oldUser.IsActivated {
		return NewError(ErrDeactivatedUser, fmt.Errorf("user %d is deactivated", oldUser.UserID))
	}
	if !oldUser.HasSetPassword {
		return NewError(ErrNonexistentPassword, fmt.Errorf("password of user %d does not exist", oldUser.UserID))
	}
	helper := legacy.NewPasswordCipherHelper()
	password := helper.Decrypt(oldUser.EncryptedPassword, oldUser.Seed, j.encryptKey)
	if hash.HexString(hash.SHA256String(password+req.Seed)) != req.PasswordHash {
		return NewError(ErrUsernamePasswordNotMatch, fmt.Errorf("username and password of user %d does not match", oldUser.UserID))
	}
	// 验证老用户没有设置手机号
	if oldUser.HasSetPhone {
		return NewError(ErrExsitSignInPhone, fmt.Errorf("signin phone has been set of user %d", req.UserId))
	}
	var tokenOfOldUser string
	// 事务方式进行写操作
	txErr := With(func() error {
		tx := j.datastore
		c := tx.BeginTx(ctx)
		if err := tx.GetError(c); err != nil {
			return err
		}

		// 生成老用户的 token
		tokenOfOldUser = uuid.New().String()

		_, err = j.datastore.CreateToken(c, tokenOfOldUser, int32(oldUser.UserID), TokenAvailableDuration)
		if err != nil {
			return NewError(ErrGetAccessTokenFailure, fmt.Errorf("failed to create access token of user %d", oldUser.UserID))
		}

		// 设置老用户的手机号
		errSetSigninPhoneByUserID := j.datastore.SetSigninPhoneByUserID(c, oldUser.UserID, newUser.SigninPhone, newUser.NationCode)
		if errSetSigninPhoneByUserID != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed set signin phone %s%s to UserId %d: %s", newUser.NationCode, newUser.SigninPhone, oldUser.UserID, errSetSigninPhoneByUserID.Error()))
		}
		// 删除新用户的token
		errDeleteToken := j.datastore.DeleteToken(c, token)
		if errDeleteToken != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to delete token: %s", errDeleteToken.Error()))
		}

		// 伪删除用户误创建的新用户
		errDeleteUser := j.datastore.DeleteUser(c, req.UserId)
		if errDeleteUser != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to delete user %d: %s", req.UserId, errDeleteUser.Error()))
		}

		tx.CommitTx(c)
		return tx.GetError(c)

	})
	if txErr != nil {
		return txErr
	}

	resp.UserId = oldUser.UserID
	resp.AccessToken = tokenOfOldUser
	return nil
}
