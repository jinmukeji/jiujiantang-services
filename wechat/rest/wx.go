package rest

import (
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
)

// Message 微信返回的信息
type Message struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgID        string `xml:"MsgId"`
	Event        string `xml:"Event"`
	EventKey     string `xml:"EventKey"`
}

// WxUser 微信用户
type WxUser struct {
	Subscribe      int32   `json:"subscribe"`       // 用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。
	OpenID         string  `json:"openid"`          // 用户的标识，对当前公众号唯一
	Nickname       string  `json:"nickname"`        // 昵称
	Sex            int32   `json:"sex"`             // 值为1时是男性，值为2时是女性
	City           string  `json:"city"`            // 用户所在城市
	Country        string  `json:"country"`         // 城市
	Province       string  `json:"province"`        // 省份
	Language       string  `json:"language"`        // 语言
	HeadImgUrl     string  `json:"headimgurl"`      // 头像
	SubscribeTime  int32   `json:"subscribe_time"`  // 关注的时间
	UnionID        string  `json:"unionid"`         // unionid
	Remark         string  `json:"remark"`          // 公众号运营者对粉丝的备注
	GroupID        int32   `json:"groupid"`         // 用户所在的分组ID
	TagIDList      []int32 `json:"tagid_list"`      // 用户被打上的标签ID列表
	SubscribeScene string  `json:"subscribe_scene"` // 来源
	QRScene        int32   `json:"qr_scene"`        // 二维码扫码场景
}

const (
	// MessageEventScan 扫码
	MessageEventScan = "SCAN"

	// MessageEventSubscribe 关注
	MessageEventSubscribe = "subscribe"

	// MessageEventClick 菜单文字消息
	MessageEventClick = "CLICK"
)

// WeChatAuth 微信公众号服务器接入配置验证回调
func (h *handler) WeChatAuth(ctx iris.Context) {
	signature := ctx.URLParam("signature")
	timestamp := ctx.URLParam("timestamp")
	nonce := ctx.URLParam("nonce")
	echostr := ctx.URLParam("echostr")

	req := new(proto.WechatCheckWxSignatureRequest)
	req.Signature = signature
	req.Timestamp = timestamp
	req.Nonce = nonce

	resp, err := h.rpcSvc.WechatCheckWxSignature(
		newRPCContext(ctx), req,
	)
	if err != nil {
		_, errText := ctx.Text("failed to check wx")
		if errText != nil {
			return
		}
		return
	}

	if resp.Ok {
		_, errText := ctx.Text(echostr)
		if errText != nil {
			return
		}
	} else {
		_, errText := ctx.Text("")
		if errText != nil {
			return
		}
	}
}

// receiveMessage 接受消息
func (h *handler) receiveMessage(ctx iris.Context) {
	var message Message
	var err error
	errReadXML := ctx.ReadXML(&message)
	if errReadXML != nil {
		return
	}
	switch message.Event {
	case MessageEventScan:
		logMsg(&message)
		err = h.SendScanOrSubscribeMessage(ctx, message, MessageEventScan)
		if err != nil {
			return
		}
	case MessageEventSubscribe:
		logMsg(&message)
		err = h.SendScanOrSubscribeMessage(ctx, message, MessageEventSubscribe)
		if err != nil {
			return
		}
	case MessageEventClick:
		logMsg(&message)
		err = h.SendTextMessage(ctx, message.FromUserName, message.EventKey)
		if err != nil {
			return
		}
	default:
	}
	// 直接回复空串，告诉微信服务器已正常处理事件推送
	_, errWriteString := ctx.WriteString("")
	if errWriteString != nil {
		return
	}
}

func logMsg(msg *Message) {
	log.Println("Content=", msg.Content,
		"ToUserName=", msg.ToUserName,
		"FromUserName=", msg.FromUserName,
		"MsgType=", msg.MsgType,
		"CreateTime=", msg.CreateTime,
		"MsgId=", msg.MsgID,
		"Event=", msg.Event,
		"EventKey", msg.EventKey)
}

// 发生文字消息
func (h *handler) SendTextMessage(ctx iris.Context, OpenID string, context string) error {
	req := new(proto.WechatSendTextMessageRequest)
	req.OpenId = OpenID
	req.Content = context
	_, err := h.rpcSvc.WechatSendTextMessage(
		newRPCContext(ctx), req,
	)
	if err != nil {
		return err
	}
	return nil
}

// SendScanOrSubscribeMessage 扫描二维码
func (h *handler) SendScanOrSubscribeMessage(ctx iris.Context, message Message, messageType string) error {
	reqScanQRCodeResquest := new(proto.ScanQRCodeRequest)
	sceneID, _ := strconv.Atoi(message.EventKey)
	reqScanQRCodeResquest.SceneId = int32(sceneID)
	date, _ := strconv.Atoi(message.CreateTime)
	createdAt, _ := ptypes.TimestampProto(time.Unix(int64(date), 0))
	reqScanQRCodeResquest.CreatedTime = createdAt
	rpcCtx := newRPCContext(ctx)
	_, err := h.rpcSvc.ScanQRCode(
		rpcCtx, reqScanQRCodeResquest,
	)
	if err != nil {
		return err
	}
	reqWechatGetWxUserByOpenIDRequest := new(proto.WechatGetWxUserByOpenIDRequest)
	reqWechatGetWxUserByOpenIDRequest.OpenId = message.FromUserName
	respWechatGetWxUserByOpenID, err := h.rpcSvc.WechatGetWxUserByOpenID(
		rpcCtx, reqWechatGetWxUserByOpenIDRequest,
	)
	if err != nil {

		return err
	}
	wxUser := respWechatGetWxUserByOpenID.UserInfo
	reqCreateWxUser := new(proto.CreateWxUserRequest)
	reqCreateWxUser.OpenId = wxUser.OpenId
	reqCreateWxUser.AvatarImageUrl = wxUser.HeadImageUrl
	reqCreateWxUser.Nickname = wxUser.Nickname
	reqCreateWxUser.UnionId = wxUser.UnionId
	if messageType == MessageEventScan {
		sceneID, _ := strconv.Atoi(message.EventKey)
		reqCreateWxUser.SceneId = int32(sceneID)
	} else {
		sceneID, _ := strconv.Atoi(strings.Split(message.EventKey, "_")[1])
		reqCreateWxUser.SceneId = int32(sceneID)
	}
	_, err = h.rpcSvc.CreateWxUser(
		rpcCtx, reqCreateWxUser,
	)
	if err != nil {
		return err
	}

	return nil
}
