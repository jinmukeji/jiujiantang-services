package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/micro/go-micro/v2/broker"

	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"

	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
)

// 创建默认微信用户的默认属性
const (
	defaultBirthday     = "1998-01-01"
	defaultGender       = generalpb.Gender_GENDER_FEMALE
	defaultHeight       = 170
	defaultWeight       = 60
	defaultZone         = "CN"
	defaultClientID     = "jm-10002"
	defaultRegisterType = "wechat"
	createdUserTopic    = "com.jinmuhealth.topic.jinmul-wx-created-user"
)

const (
	// 日期格式
	dateFormat = "2006-01-02"
)

// Message 消息
type Message struct {
	SceneID int32 `json:"scene_id"`
	UserID  int32 `json:"user_id"`
}

// CreateWxUser 创建微信User
func (j *JinmuHealth) CreateWxUser(ctx context.Context, req *corepb.CreateWxUserRequest, resp *corepb.CreateWxUserResponse) error {
	exist, errExistWXUser := j.datastore.ExistWXUser(ctx, req.UnionId)
	if errExistWXUser != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check wx user existence by unionId %s: %s", req.UnionId, errExistWXUser.Error()))
	}
	qrcode, errFindQRCodeBySceneID := j.datastore.FindQRCodeBySceneID(ctx, req.SceneId)
	if errFindQRCodeBySceneID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find QR code by sceneID %d: %s", req.SceneId, errFindQRCodeBySceneID.Error()))
	}
	if exist {
		wxUser, err := j.datastore.FindWXUserByUnionID(ctx, req.UnionId)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to wx code by unionID %d: %s", req.SceneId, err.Error()))
		}
		return publishMessage(Message{
			qrcode.SceneID, wxUser.UserID,
		})
	}
	now := time.Now()
	birthday, _ := time.Parse(dateFormat, defaultBirthday)
	reqSignUpUserViaUsernamePassword := new(jinmuidpb.SignUpUserViaUsernamePasswordRequest)
	reqSignUpUserViaUsernamePassword.ClientId = defaultClientID
	reqSignUpUserViaUsernamePassword.RegisterType = defaultRegisterType
	reqSignUpUserViaUsernamePassword.RegisterTime, _ = ptypes.TimestampProto(now)
	reqSignUpUserViaUsernamePassword.IsSkipVerifyToken = true
	reqSignUpUserViaUsernamePassword.IsNeedSetProfile = true
	birthdayProto, _ := ptypes.TimestampProto(birthday)
	reqSignUpUserViaUsernamePassword.Profile = &jinmuidpb.UserProfile{
		Nickname:     req.Nickname,
		BirthdayTime: birthdayProto,
		Gender:       defaultGender,
		Height:       defaultHeight,
		Weight:       defaultWeight,
	}
	reqSignUpUserViaUsernamePassword.Zone = defaultZone
	respSignUpUserViaUsernamePassword, errSignUpUserViaUsernamePassword := j.jinmuidSvc.SignUpUserViaUsernamePassword(ctx, reqSignUpUserViaUsernamePassword)
	if errSignUpUserViaUsernamePassword != nil {
		return errSignUpUserViaUsernamePassword
	}
	wxUser := &mysqldb.WXUser{
		UnionID:        req.UnionId,
		OpenID:         req.OpenId,
		Nickname:       req.Nickname,
		AvatarImageURL: req.AvatarImageUrl,
		OriginID:       j.wechat.GetOriginID(),
		UserID:         respSignUpUserViaUsernamePassword.UserId,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	errCreateWXUser := j.datastore.CreateWXUser(ctx, wxUser)
	if errCreateWXUser != nil {
		return errCreateWXUser
	}
	return publishMessage(Message{
		qrcode.SceneID, wxUser.UserID,
	})
}

// WechatGetUserInfo 得到微信用户信息
func (j *JinmuHealth) WechatGetUserInfo(ctx context.Context, req *corepb.WechatGetUserInfoRequest, resp *corepb.WechatGetUserInfoResponse) error {
	w := j.wechat
	user, err := w.GetUserInfoByOauthCode(req.Code)
	if err != nil {
		return err
	}

	resp.OpenId = user.OpenId
	resp.Nickname = user.Nickname
	resp.Sex = int32(user.Sex)
	resp.City = user.City
	resp.Province = user.Province
	resp.Country = user.Country
	resp.HeadImageUrl = user.HeadImageURL
	resp.Privilege = user.Privilege
	resp.UnionId = user.UnionId

	return nil
}

// WechatGetWxUserByOpenID 通过openID获取微信用户信息
func (j *JinmuHealth) WechatGetWxUserByOpenID(ctx context.Context, req *corepb.WechatGetWxUserByOpenIDRequest, resp *corepb.WechatGetWxUserByOpenIDResponse) error {
	w := j.wechat
	user, err := w.GetUserInfo(req.OpenId)
	if err != nil {
		return err
	}

	resp.UserInfo = new(corepb.WxUserInfo)
	resp.UserInfo.IsSubscriber = int64(user.IsSubscriber)
	resp.UserInfo.OpenId = user.OpenId
	resp.UserInfo.Nickname = user.Nickname
	resp.UserInfo.Sex = int64(user.Sex)
	resp.UserInfo.City = user.City
	resp.UserInfo.Country = user.Country
	resp.UserInfo.Province = user.Province
	resp.UserInfo.Language = user.Language
	resp.UserInfo.HeadImageUrl = user.HeadImageURL
	resp.UserInfo.SubscribeTime = user.SubscribeTime
	resp.UserInfo.UnionId = user.UnionId
	resp.UserInfo.Remark = user.Remark
	resp.UserInfo.GroupId = user.GroupId
	resp.UserInfo.TagIdList = intToInt64Array(user.TagIdList)

	return nil
}

func intToInt64Array(in []int) []int64 {
	out := make([]int64, len(in))
	for i := range in {
		out[i] = int64(in[i])
	}

	return out
}

// publishMessage 广播消息
func publishMessage(message Message) error {
	data, err := json.Marshal(&message)
	if err != nil {
		return err
	}
	msg := &broker.Message{
		Body: data,
	}

	if err := broker.Publish(createdUserTopic, msg); err != nil {
		return fmt.Errorf("[pub] failed: %v", err)
	}
	return nil
}

// GetWechatUser 得到微信user
func (j *JinmuHealth) GetWechatUser(ctx context.Context, req *corepb.GetWechatUserRequest, resp *corepb.GetWechatUserResponse) error {
	wxUser, err := j.datastore.FindWXUserByOpenID(ctx, req.OpenId)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find WX user by openId %s: %s", req.OpenId, err.Error()))
	}

	resp.OpenId = wxUser.OpenID
	resp.UnionId = wxUser.UnionID
	resp.UserId = wxUser.UserID
	resp.Nickname = wxUser.Nickname
	resp.AvatarImageUrl = wxUser.AvatarImageURL

	return nil
}
