package rest

import (
	"errors"
	"fmt"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

const (
	// SimpleChinese 简体中文
	SimpleChinese = "zh-Hans"
	// TraditionalChinese 繁体中文
	TraditionalChinese = "zh-Hant"
	// English 英文
	English = "en"
)

// Language 语言
type Language struct {
	Language string `json:"language"`
}

// SetWebLanguage 设置语言
func (h *webHandler) SetWebLanguage(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(jinmuidpb.SetJinmuIDWebLanguageRequest)
	req.UserId = int32(userID)
	var body Language
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	if body.Language == "" {
		writeError(ctx, wrapError(ErrEmptyLanguage, "", errors.New("language is empty")), false)
		return
	}
	protoLanguage, errmapRestLanguageToProto := mapRestLanguageToProto(body.Language)
	if errmapRestLanguageToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestLanguageToProto), false)
		return
	}
	req.Language = protoLanguage
	_, errSetJinmuIDWebLanguage := h.rpcSvc.SetJinmuIDWebLanguage(
		newRPCContext(ctx), req,
	)
	if errSetJinmuIDWebLanguage != nil {
		writeRpcInternalError(ctx, errSetJinmuIDWebLanguage, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}

// GetWebLanguage 得到语言
func (h *webHandler) GetWebLanguage(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(jinmuidpb.GetJinmuIDWebLanguageRequest)
	req.UserId = int32(userID)
	resp, errGetJinmuIDWebLanguage := h.rpcSvc.GetJinmuIDWebLanguage(
		newRPCContext(ctx), req,
	)
	if errGetJinmuIDWebLanguage != nil {
		writeRpcInternalError(ctx, errGetJinmuIDWebLanguage, false)
		return
	}
	stringLanguage, ermapProtoLanguageToRest := mapProtoLanguageToRest(resp.Language)
	if ermapProtoLanguageToRest != nil {
		// 默认使用简体中文
		stringLanguage = SimpleChinese
	}
	rest.WriteOkJSON(ctx, Language{
		Language: stringLanguage,
	})
}

func mapProtoLanguageToRest(lanuage generalpb.Language) (string, error) {
	switch lanuage {
	case generalpb.Language_LANGUAGE_INVALID:
		return "", fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_INVALID)
	case generalpb.Language_LANGUAGE_UNSET:
		return "", fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_UNSET)
	case generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE:
		return SimpleChinese, nil
	case generalpb.Language_LANGUAGE_TRADITIONAL_CHINESE:
		return TraditionalChinese, nil
	case generalpb.Language_LANGUAGE_ENGLISH:
		return English, nil
	}
	return SimpleChinese, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_INVALID)
}

func mapRestLanguageToProto(language string) (generalpb.Language, error) {
	switch language {
	case SimpleChinese:
		return generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE, nil
	case TraditionalChinese:
		return generalpb.Language_LANGUAGE_TRADITIONAL_CHINESE, nil
	case English:
		return generalpb.Language_LANGUAGE_ENGLISH, nil
	}
	return generalpb.Language_LANGUAGE_INVALID, fmt.Errorf("invalid string language %s", language)
}
