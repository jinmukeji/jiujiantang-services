package handler

import (
	"context"
	"errors"
	"regexp"

	"github.com/jinmukeji/go-pkg/v2/areacode"

	"fmt"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// UserValidateUsernameOrPhone 验证手机号码和用户名是否存在
func (j *JinmuIDService) UserValidateUsernameOrPhone(ctx context.Context, req *proto.UserValidateUsernameOrPhoneRequest, resp *proto.UserValidateUsernameOrPhoneResponse) error {
	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_UNKNOWN {
		return NewError(ErrInvalidSecureQuestionValidationMethod, errors.New("invalid secure queston validation type"))
	}
	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_PHONE {
		if req.Username != "" {
			return NewError(ErrInvalidValidationValue, errors.New("username should be empty when validation type is phone"))
		}

		errValidatePhoneFormat := validatePhoneFormat(req.Phone, req.NationCode)
		if errValidatePhoneFormat != nil {
			return errValidatePhoneFormat
		}

		existSignInPhone, _ := j.datastore.ExistSignInPhone(ctx, req.Phone, req.NationCode)
		if !existSignInPhone {
			return NewError(ErrNoneExistentPhone, fmt.Errorf("phone %s%s doesn't exist", req.NationCode, req.Phone))
		}

		// FindUserIDByPhone 通过电话号码找到userID
		userID, errFindUserIDByPhone := j.datastore.FindUserIDByPhone(ctx, req.Phone, req.NationCode)
		if errFindUserIDByPhone != nil {
			return NewError(ErrDatabase, fmt.Errorf("fail to find user by phone %s%s: %s", req.NationCode, req.Phone, errFindUserIDByPhone.Error()))
		}
		user, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, userID)
		if errFindUserByUserID != nil {
			return NewError(ErrDatabase, fmt.Errorf("fail to find user by userID %d: %s", userID, errFindUserByUserID.Error()))
		}
		resp.SecureEmail = user.SecureEmail
		resp.HasSetEmail = user.HasSetEmail
		return nil

	}
	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_USERNAME {
		if req.Phone != "" || req.NationCode != "" {
			return NewError(ErrInvalidValidationValue, errors.New("phone should be empty when validation type is username"))
		}
		existUsername, _ := j.datastore.ExistUsername(ctx, req.Username)
		if !existUsername {
			return NewError(ErrNonexistentUsername, fmt.Errorf("username %s doesn't exist", req.Username))
		}
		// 通过用户名找到User
		user, errFindUserByUsername := j.datastore.FindUserByUsername(ctx, req.Username)
		if errFindUserByUsername != nil {
			return NewError(ErrDatabase, fmt.Errorf("fail to find user by username %s: %s", req.Username, errFindUserByUsername.Error()))
		}
		user, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, user.UserID)
		if errFindUserByUserID != nil {
			return NewError(ErrDatabase, fmt.Errorf("fail to find user by userID %d: %s", user.UserID, errFindUserByUserID.Error()))
		}
		resp.SecureEmail = user.SecureEmail
		resp.HasSetEmail = user.HasSetEmail
		return nil
	}
	return nil
}

// 验证电话号码格式
func validatePhoneFormat(phoneNumber string, nationCode string) error {
	switch nationCode {
	case areacode.Mainland:
		// 11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.HongKong:
		// 8位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{8}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Macao:
		// 8位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{8}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Taiwan:
		// 9位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{9}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.UnitedStatesOrCanada:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.UnitedKingdom:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Japan:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Brunei:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Cambodia:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Indonesia:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Laos:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Malaysia:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Myanmar:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Philippines:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Singapore:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Thailand:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Vietnam:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	case areacode.Korea:
		// 7-11位数字
		if matchPhoneFormat, _ := regexp.MatchString("^[0-9]{7,11}$", phoneNumber); !matchPhoneFormat {
			return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
		}
		return nil
	}
	return NewError(ErrWrongFormatPhone, errors.New("wrong format of phone"))
}
