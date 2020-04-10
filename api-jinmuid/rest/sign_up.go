package rest

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// SignUpByPhoneVerificationNumberBody 手机号验证码注册body
type SignUpByPhoneVerificationNumberBody struct {
	Phone              string      `json:"phone"`               // 电话
	VerificationNumber string      `json:"verification_number"` // 验证号
	NationCode         string      `json:"nation_code"`         // 区号
	Language           string      `json:"language"`            // 常用语言
	PlainPassword      string      `json:"plain_password"`      // 密码
	UserProfile        UserProfile `json:"user_profile"`        // 用户档案
}

// SignUpByPhoneVerificationNumber 手机号验证号注册
func (h *webHandler) SignUpByPhoneVerificationNumber(ctx iris.Context) {
	var body SignUpByPhoneVerificationNumberBody
	err := ctx.ReadJSON(&body)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	if !checkNationCode(body.NationCode) {
		writeError(ctx, wrapError(ErrNationCode, "", fmt.Errorf("nation code %s is wrong", body.NationCode)), false)
		return
	}
	req := new(proto.UserSignUpByPhoneRequest)
	req.Phone = body.Phone
	req.VerificationNumber = body.VerificationNumber
	req.NationCode = body.NationCode
	protoLanguage, errmapRestLanguageToProto := mapRestLanguageToProto(body.Language)
	if errmapRestLanguageToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestLanguageToProto), false)
		return
	}
	req.Language = protoLanguage
	req.PlainPassword = body.PlainPassword
	birthday, _ := ptypes.TimestampProto(body.UserProfile.Birthday)
	protoGender, errmapRestGenderToProto := mapRestGenderToProto(body.UserProfile.Gender)
	if errmapRestGenderToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestGenderToProto), false)
		return
	}
	req.Profile = &proto.UserProfile{
		Nickname:     body.UserProfile.Nickname,
		Gender:       protoGender,
		BirthdayTime: birthday,
		Weight:       body.UserProfile.Weight,
		Height:       body.UserProfile.Height,
	}
	resp, errUserSignUpByPhone := h.rpcSvc.UserSignUpByPhone(newRPCContext(ctx), req)
	if errUserSignUpByPhone != nil {
		writeRpcInternalError(ctx, errUserSignUpByPhone, false)
		return
	}
	rest.WriteOkJSON(ctx, resp)
}
