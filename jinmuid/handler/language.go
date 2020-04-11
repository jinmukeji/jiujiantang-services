package handler

import (
	"context"
	"errors"
	"fmt"

	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
)

const (
	// LanguageSimpleChinese 简体
	LanguageSimpleChinese = "zh-Hans"
	// LanguageTraditionalChinese 繁体
	LanguageTraditionalChinese = "zh-Hant"
	// LanguageEnglish 英文
	LanguageEnglish = "en"
)

// SetJinmuIDWebLanguage 设置jinmuID web语言
func (j *JinmuIDService) SetJinmuIDWebLanguage(ctx context.Context, req *jinmuidpb.SetJinmuIDWebLanguageRequest, resp *jinmuidpb.SetJinmuIDWebLanguageResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
	}
	langauge, errmapProtoLanguageToDB := mapProtoLanguageToDB(req.Language)
	if errmapProtoLanguageToDB != nil {
		return NewError(ErrInvalidUser, errmapProtoLanguageToDB)
	}
	errSetLanguageByUserID := j.datastore.SetLanguageByUserID(ctx, req.UserId, langauge)
	if errSetLanguageByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set the language of user %d to %s: %s", req.UserId, langauge, errSetLanguageByUserID.Error()))
	}
	return nil
}

// GetJinmuIDWebLanguage  得到jinmuID web语言
func (j *JinmuIDService) GetJinmuIDWebLanguage(ctx context.Context, req *jinmuidpb.GetJinmuIDWebLanguageRequest, resp *jinmuidpb.GetJinmuIDWebLanguageResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
	}
	user, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, req.UserId)
	if errFindUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by userId %d: %s", req.UserId, errFindUserByUserID.Error()))
	}
	if user.HasSetLanguage {
		protoLanguage, errMapDBLanguageToProto := mapDBLanguageToProto(string(user.Language))
		if errMapDBLanguageToProto != nil {
			return NewError(ErrInvalidUser, errMapDBLanguageToProto)
		}

		resp.Language = protoLanguage
	} else {
		resp.Language = generalpb.Language_LANGUAGE_UNSET
	}
	return nil
}

func mapProtoLanguageToDB(lanuage generalpb.Language) (string, error) {
	switch lanuage {
	case generalpb.Language_LANGUAGE_INVALID:
		return LanguageSimpleChinese, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_INVALID)
	case generalpb.Language_LANGUAGE_UNSET:
		return LanguageTraditionalChinese, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_UNSET)
	case generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE:
		return LanguageSimpleChinese, nil
	case generalpb.Language_LANGUAGE_TRADITIONAL_CHINESE:
		return LanguageTraditionalChinese, nil
	case generalpb.Language_LANGUAGE_ENGLISH:
		return LanguageEnglish, nil
	}
	return LanguageSimpleChinese, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_INVALID)
}

func mapDBLanguageToProto(language string) (generalpb.Language, error) {
	switch language {
	case LanguageSimpleChinese:
		return generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE, nil
	case LanguageTraditionalChinese:
		return generalpb.Language_LANGUAGE_TRADITIONAL_CHINESE, nil
	case LanguageEnglish:
		return generalpb.Language_LANGUAGE_ENGLISH, nil
	}
	return generalpb.Language_LANGUAGE_INVALID, fmt.Errorf("invalid string language %s", language)
}
