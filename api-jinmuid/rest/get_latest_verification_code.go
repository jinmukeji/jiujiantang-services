package rest

import (
	"fmt"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

const (
	// SendViaPhone 发送途径为手机号码
	SendViaPhone = "phone"
	// SendViaEmail 发送途径为邮箱
	SendViaEmail = "email"
)

// GetLatestVerificationCode 最新验证码请求
type GetLatestVerificationCode struct {
	SendVia    string `json:"send_via"`    // 发送途径
	Email      string `json:"email"`       // 邮箱
	Phone      string `json:"phone"`       // 手机号码
	NationCode string `json:"nation_code"` // 区号
}

// LatestVerificationCodeBody 获取最新验证码请求
type LatestVerificationCodeBody struct {
	SendInformation []GetLatestVerificationCode `json:"send_information"` // 发送信息
}

// LatestVerificationCode 最新验证码
type LatestVerificationCode struct {
	Email            string `json:"email"`             // 邮箱
	Phone            string `json:"phone"`             // 手机号码
	NationCode       string `json:"nation_code"`       // 区号
	VerificationCode string `json:"verification_code"` // 验证码
}

// GetLatestVerificationCodes 获取最新验证码
func (h *webHandler) GetLatestVerificationCodes(ctx iris.Context) {
	var reqGetLatestVerificationCodes LatestVerificationCodeBody
	errReadJSON := ctx.ReadJSON(&reqGetLatestVerificationCodes)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req := new(proto.GetLatestVerificationCodesRequest)
	latestVerificationCodes := make([]*proto.SingleGetLatestVerificationCode, len(reqGetLatestVerificationCodes.SendInformation))
	for idx, item := range reqGetLatestVerificationCodes.SendInformation {
		protoSendVia, errmapRestSendViaToProto := mapRestSendViaToProto(item.SendVia)
		if errmapRestSendViaToProto != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestSendViaToProto), false)
			return
		}
		latestVerificationCodes[idx] = &proto.SingleGetLatestVerificationCode{
			SendVia:    protoSendVia,
			Email:      item.Email,
			Phone:      item.Phone,
			NationCode: item.NationCode,
		}
	}
	req.SendTo = latestVerificationCodes
	resp, errGetLatestVerificationCodes := h.rpcSvc.GetLatestVerificationCodes(
		newRPCContext(ctx), req,
	)
	if errGetLatestVerificationCodes != nil {
		writeRpcInternalError(ctx, errGetLatestVerificationCodes, false)
		return
	}
	LatestVerificationCodes := make([]LatestVerificationCode, len(resp.LatestVerificationCodes))
	for idx, item := range resp.LatestVerificationCodes {
		LatestVerificationCodes[idx].Email = item.Email
		LatestVerificationCodes[idx].Phone = item.Phone
		LatestVerificationCodes[idx].NationCode = item.NationCode
		LatestVerificationCodes[idx].VerificationCode = item.VerificationCode
	}
	rest.WriteOkJSON(ctx, LatestVerificationCodes)
}

// mapRestSendViaToProto 将 rest 使用的 string 类型的 发送验证码的方式 send_via 转换为 proto 类型
func mapRestSendViaToProto(sendType string) (proto.SendVia, error) {
	switch sendType {
	case SendViaPhone:
		return proto.SendVia_SEND_VIA_PHONE_SEND_VIA, nil
	case SendViaEmail:
		return proto.SendVia_SEND_VIA_USERNAME_SEND_VIA, nil
	}
	return proto.SendVia_SEND_VIA_INVALID, fmt.Errorf("invalid string send via %s", sendType)
}
