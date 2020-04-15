package rest

import (
	"time"

	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/kataras/iris/v12"
)

const (
	// PhoneMvc 手机号验证码
	PhoneMvc = "phone_mvc"
	// UsernamePassword 用户名密码
	UsernamePassword = "username_password"
	// PhonePassword 手机号密码
	PhonePassword = "phone_password"
)

// SignInBody 登录的body
type SignInBody struct {
	SignInMethod   string `json:"sign_in_method"`  // 登陆方式
	Username       string `json:"username"`        // 用户名
	Phone          string `json:"phone"`           // 手机号
	MVC            string `json:"mvc"`             // 验证码
	HashedPassword string `json:"hashed_password"` // hash后的密码
	Seed           string `json:"seed"`            // 随机种子
	SerialNumber   string `json:"serial_number"`   // 序列号
	NationCode     string `json:"nation_code"`     // 国际代码
	SignInMachine  string `json:"sign_in_machine"` // 登陆机器
}

// SignInResp 登陆的返回
type SignInResp struct {
	UserID             int32     `json:"user_id"`
	AccessToken        string    `json:"access_token"`
	HasSetUserProfile  bool      `json:"has_set_user_profile"`
	HasSetRegion       bool      `json:"has_set_region"`
	HasSetLanguage     bool      `json:"has_set_language"`
	HasSetPassword     bool      `json:"has_set_password"`
	HasSetPhone        bool      `json:"has_set_phone"`
	IsProfileCompleted bool      `json:"is_profile_completed"`
	ExpiredAt          time.Time `json:"expired_at"`
}

// SignInMachine 登陆的机器
type SignInMachine struct {
	SignInMachine string    `json:"sign_in_machine"`
	SignInTime    time.Time `json:"sign_in_time"`
}

// SignIn 登录
func (h *webHandler) SignIn(ctx iris.Context) {
	var signInBody SignInBody
	err := ctx.ReadJSON(&signInBody)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	switch signInBody.SignInMethod {
	case PhoneMvc:
		if !checkNationCode(signInBody.NationCode) {
			writeError(ctx, wrapError(ErrNationCode, "", fmt.Errorf("nation code %s is wrong", signInBody.NationCode)), false)
			return
		}
		req := new(proto.UserSignInByPhoneVCRequest)
		req.Phone = signInBody.Phone
		req.Mvc = signInBody.MVC
		req.SerialNumber = signInBody.SerialNumber
		req.NationCode = signInBody.NationCode
		req.SignInMachine = signInBody.SignInMachine
		req.Ip = ctx.RemoteAddr()
		resp, errUserSignInByPhoneVC := h.rpcSvc.UserSignInByPhoneVC(newRPCContext(ctx), req)
		if errUserSignInByPhoneVC != nil {
			writeRpcInternalError(ctx, errUserSignInByPhoneVC, false)
			return
		}
		expiredAt, _ := ptypes.Timestamp(resp.ExpiredTime)
		rest.WriteOkJSON(ctx, SignInResp{
			UserID:             resp.UserId,
			AccessToken:        resp.AccessToken,
			HasSetLanguage:     resp.HasSetLanguage,
			HasSetRegion:       resp.HasSetRegion,
			HasSetUserProfile:  resp.HasSetUserProfile,
			HasSetPassword:     resp.HasSetPassword,
			HasSetPhone:        resp.HasSetPhone,
			IsProfileCompleted: resp.IsProfileCompleted,
			ExpiredAt:          expiredAt.UTC(),
		})
		return
	case UsernamePassword:
		req := new(proto.UserSignInByUsernamePasswordRequest)
		req.Username = signInBody.Username
		req.HashedPassword = signInBody.HashedPassword
		req.Seed = signInBody.Seed
		req.SignInMachine = signInBody.SignInMachine
		req.Ip = ctx.RemoteAddr()
		resp, errUserSignInByUsernamePassword := h.rpcSvc.UserSignInByUsernamePassword(newRPCContext(ctx), req)
		if errUserSignInByUsernamePassword != nil {
			writeRpcInternalError(ctx, errUserSignInByUsernamePassword, false)
			return
		}
		expiredAt, _ := ptypes.Timestamp(resp.ExpiredTime)
		rest.WriteOkJSON(ctx, SignInResp{
			UserID:             resp.UserId,
			AccessToken:        resp.AccessToken,
			HasSetLanguage:     resp.HasSetLanguage,
			HasSetRegion:       resp.HasSetRegion,
			HasSetUserProfile:  resp.HasSetUserProfile,
			HasSetPassword:     resp.HasSetPassword,
			HasSetPhone:        resp.HasSetPhone,
			IsProfileCompleted: resp.IsProfileCompleted,
			ExpiredAt:          expiredAt.UTC(),
		})
		return
	case PhonePassword:
		if !checkNationCode(signInBody.NationCode) {
			writeError(ctx, wrapError(ErrNationCode, "", fmt.Errorf("nation code %s is wrong", signInBody.NationCode)), false)
			return
		}
		req := new(proto.UserSignInByPhonePasswordRequest)
		req.Phone = signInBody.Phone
		req.HashedPassword = signInBody.HashedPassword
		req.Seed = signInBody.Seed
		req.NationCode = signInBody.NationCode
		req.SignInMachine = signInBody.SignInMachine
		req.Ip = ctx.RemoteAddr()
		resp, errUserSignInByPhonePassword := h.rpcSvc.UserSignInByPhonePassword(newRPCContext(ctx), req)
		if errUserSignInByPhonePassword != nil {
			writeRpcInternalError(ctx, errUserSignInByPhonePassword, false)
			return
		}
		expiredAt, _ := ptypes.Timestamp(resp.ExpiredTime)
		rest.WriteOkJSON(ctx, SignInResp{
			UserID:             resp.UserId,
			AccessToken:        resp.AccessToken,
			HasSetLanguage:     resp.HasSetLanguage,
			HasSetRegion:       resp.HasSetRegion,
			HasSetUserProfile:  resp.HasSetUserProfile,
			HasSetPassword:     resp.HasSetPassword,
			HasSetPhone:        resp.HasSetPhone,
			IsProfileCompleted: resp.IsProfileCompleted,
			ExpiredAt:          expiredAt.UTC(),
		})
		return
	}
}

// UserGetSignInMachines 得到登录的设备
func (h *webHandler) UserGetSignInMachines(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.UserGetSignInMachinesRequest)
	req.UserId = int32(userID)
	resp, errUserGetSignInMachines := h.rpcSvc.UserGetSignInMachines(newRPCContext(ctx), req)
	if errUserGetSignInMachines != nil {
		writeRpcInternalError(ctx, errUserGetSignInMachines, false)
		return
	}
	signInMachines := make([]SignInMachine, len(resp.Machines))
	for idx, machine := range resp.Machines {
		signInTime, _ := ptypes.Timestamp(machine.SignInTime)
		signInMachines[idx] = SignInMachine{
			SignInMachine: machine.SignInMachine,
			SignInTime:    signInTime,
		}
	}
	rest.WriteOkJSON(ctx, signInMachines)
}
