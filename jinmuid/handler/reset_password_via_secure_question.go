package handler

import (
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	crypto "github.com/jinmukeji/go-pkg/crypto/encrypt/legacy"
	"github.com/jinmukeji/go-pkg/crypto/rand"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

const (
	// MinpasswordLength 密码的最小长度
	MinpasswordLength = 8
	// MaxpasswordLength 密码的最大长度
	MaxpasswordLength = 20
)

// UserResetPasswordViaSecureQuestions 通过密保问题重置密码
func (j *JinmuIDService) UserResetPasswordViaSecureQuestions(ctx context.Context, req *proto.UserResetPasswordViaSecureQuestionsRequest, resp *proto.UserResetPasswordViaSecureQuestionsResponse) error {
	clientID, _ := ClientIDFromContext(ctx)
	seed, _ := rand.RandomStringWithMask(rand.MaskLetterDigits, 4)
	helper := crypto.NewPasswordCipherHelper()
	encryptedPassword := helper.Encrypt(req.Password, seed, j.encryptKey)
	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_UNKNOWN {
		return NewError(ErrInvalidSecureQuestionValidationMethod, errors.New("invalid secure queston validation type"))
	}

	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_PHONE {
		if req.Username != "" {
			return NewError(ErrInvalidValidationValue, errors.New("non-empty username when getting secure questions by phone"))
		}

		errValidatePhoneFormat := validatePhoneFormat(req.Phone, req.NationCode)
		if errValidatePhoneFormat != nil {
			return errValidatePhoneFormat
		}

		existSignInPhone, _ := j.datastore.ExistSignInPhone(ctx, req.Phone, req.NationCode)
		if !existSignInPhone {
			return NewError(ErrNoneExistentPhone, fmt.Errorf("phone %s%s doesn't exist", req.NationCode, req.Phone))
		}
		questions, err := j.datastore.FindSecureQuestionByPhone(ctx, req.Phone, req.NationCode)
		if err != nil {
			return NewError(ErrCurrentSecureQuestionsNotSet, fmt.Errorf("failed to find secure questions by phone %s%s: %s", req.NationCode, req.Phone, err.Error()))
		}
		user, errFindUserByPhone := j.datastore.FindUserByPhone(ctx, req.Phone, req.NationCode)
		if errFindUserByPhone != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to find username by phone %s%s: %s", req.NationCode, req.Phone, errFindUserByPhone.Error()))
		}
		password := helper.Decrypt(user.EncryptedPassword, user.Seed, j.encryptKey)
		// 判断密码是否跟以前一样
		if password == req.Password {
			return NewError(ErrSamePassword, errors.New("new password cannot equals old password"))
		}
		// 比较密保是否正确
		wrongQuestions, errCompareSecureQuestion := compareSecureQuestion(req.SecureQuestions, questions)
		if errCompareSecureQuestion != nil {
			return errCompareSecureQuestion
		}
		if len(wrongQuestions) == 0 {
			// TODO: SetPasswordByPhone,DeleteTokenByUserID,CreateAuditUserCredentialUpdate要写在同一个事务中
			resp.Result = true
			errSetPasswordByPhone := j.datastore.SetPasswordByPhone(ctx, req.Phone, req.NationCode, encryptedPassword, seed)
			if errSetPasswordByPhone != nil {
				return NewError(ErrDatabase, fmt.Errorf("failed to set password by phone %s%s: %s", req.NationCode, req.Phone, errSetPasswordByPhone.Error()))
			}
			userID, errFindUserIDByPhone := j.datastore.FindUserIDByPhone(ctx, req.Phone, req.NationCode)
			if errFindUserIDByPhone != nil {
				return NewError(ErrDatabase, fmt.Errorf("failed to find userID by phone %s%s: %s", req.NationCode, req.Phone, errFindUserIDByPhone.Error()))
			}
			errDeleteTokenByUserID := j.datastore.DeleteTokenByUserID(ctx, userID)
			if errDeleteTokenByUserID != nil {
				return NewError(ErrDatabase, fmt.Errorf("failed to delete token by user %d: %s", userID, errDeleteTokenByUserID.Error()))
			}
			now := time.Now()
			record := &mysqldb.AuditUserCredentialUpdate{
				UserID:            userID,
				ClientID:          clientID,
				UpdatedRecordType: mysqldb.PasswordUpdated,
				CreatedAt:         now,
				UpdatedAt:         now,
			}
			errCreateAuditUserCredentialUpdate := j.datastore.CreateAuditUserCredentialUpdate(ctx, record)
			if errCreateAuditUserCredentialUpdate != nil {
				return NewError(ErrDatabase, fmt.Errorf("failed to create audit user credential update: %s", errCreateAuditUserCredentialUpdate.Error()))
			}
		} else {
			resp.WrongSecureQuestionKeys = wrongQuestions
		}
	}

	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_USERNAME {
		if req.Phone != "" || req.NationCode != "" {
			return NewError(ErrInvalidValidationValue, errors.New("phone should be empty when validation type is username"))
		}

		user, FindUserByUsername := j.datastore.FindUserByUsername(ctx, req.Username)
		if FindUserByUsername != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to find user by username %s: %s", req.Username, FindUserByUsername.Error()))
		}
		if !user.HasSetPassword {
			return nil
		}
		password := helper.Decrypt(user.EncryptedPassword, user.Seed, j.encryptKey)
		// 判断密码是否跟以前一样
		if password == req.Password {
			return NewError(ErrSamePassword, errors.New("new password cannot equals old password"))
		}
		existUsername, _ := j.datastore.ExistUsername(ctx, req.Username)
		if !existUsername {
			return NewError(ErrNonexistentUsername, fmt.Errorf("username %s doesn't exist", req.Username))
		}

		questions, errFindSecureQuestionByUsername := j.datastore.FindSecureQuestionByUsername(ctx, req.Username)
		if errFindSecureQuestionByUsername != nil {
			return NewError(ErrCurrentSecureQuestionsNotSet, fmt.Errorf("failed to find secure questions by username %s: %s", req.Username, errFindSecureQuestionByUsername.Error()))
		}
		// 比较密保是否正确
		wrongQuestions, err := compareSecureQuestion(req.SecureQuestions, questions)
		if err != nil {
			return err
		}
		if len(wrongQuestions) == 0 {
			// TODO: SetPasswordByUsername,CreateAuditUserCredentialUpdate,DeleteTokenByUserID要写在同一个事务中
			resp.Result = true
			errSetPasswordByUsername := j.datastore.SetPasswordByUsername(ctx, req.Username, encryptedPassword, seed)
			if errSetPasswordByUsername != nil {
				return NewError(ErrDatabase, fmt.Errorf("failed to set password by username %s: %s", req.Username, errSetPasswordByUsername.Error()))
			}
			userID, errFindUserIDByUsername := j.datastore.FindUserIDByUsername(ctx, req.Username)
			if errFindUserIDByUsername != nil {
				return NewError(ErrDatabase, fmt.Errorf("failed to find userID by username %s: %s", req.Username, errFindUserIDByUsername.Error()))
			}
			now := time.Now()
			record := &mysqldb.AuditUserCredentialUpdate{
				UserID:            userID,
				ClientID:          clientID,
				UpdatedRecordType: mysqldb.PasswordUpdated,
				CreatedAt:         now,
				UpdatedAt:         now,
			}
			errCreateAuditUserCredentialUpdate := j.datastore.CreateAuditUserCredentialUpdate(ctx, record)
			if errCreateAuditUserCredentialUpdate != nil {
				return NewError(ErrDatabase, errors.New("failed to create audit user credential"))
			}
			errDeleteTokenByUserID := j.datastore.DeleteTokenByUserID(ctx, userID)
			if errDeleteTokenByUserID != nil {
				return NewError(ErrDatabase, fmt.Errorf("failed to delete token by user %d: %s", userID, errDeleteTokenByUserID.Error()))
			}
		} else {
			resp.WrongSecureQuestionKeys = wrongQuestions
		}
	}

	return nil
}
