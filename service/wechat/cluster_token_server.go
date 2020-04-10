package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/micro/go-micro/broker"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/util"
)

//
// 参考信息： https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140183
//
// 公众平台的API调用所需的access_token的使用及生成方式说明：
//
// 1、建议公众号开发者使用中控服务器统一获取和刷新Access_token，其他业务逻辑服务器所使用的access_token均来自于该中控服务器，
// 	不应该各自去刷新，否则容易造成冲突，导致access_token覆盖而影响业务；
//
// 2、目前Access_token的有效期通过返回的expire_in来传达，目前是7200秒之内的值。
// 	中控服务器需要根据这个有效时间提前去刷新新access_token。在刷新过程中，中控服务器可对外继续输出的老access_token，
// 	此时公众平台后台会保证在5分钟内，新老access_token都可用，这保证了第三方业务的平滑过渡；
//
// 3、Access_token的有效时间可能会在未来有调整，所以中控服务器不仅需要内部定时主动刷新，还需要提供被动刷新access_token的接口，
// 	这样便于业务服务器在API调用获知access_token已超时的情况下，可以触发access_token的刷新流程。
//

// ClusterAccessTokenServer 用于集群环境的 AccessToken 中控 Server
type ClusterAccessTokenServer interface {
	// 当前节点的 Instance ID
	InstanceID() string
}

var _ ClusterAccessTokenServer = (*DefaultClusterAccessTokenServer)(nil)

// DefaultClusterAccessTokenServer 默认的用于集群环境的 AccessToken 中控 Server
type DefaultClusterAccessTokenServer struct {
	instanceID uuid.UUID

	appID      string
	appSecret  string
	httpClient *http.Client

	refreshTokenRequestChan  chan string             // chan currentToken
	refreshTokenResponseChan chan refreshTokenResult // chan {token, err}

	clusterRefreshTokenTopic    string
	refreshTokenFormClusterChan chan accessToken // chan {Token, ExpiresIn}
	broker                      broker.Broker    // pub/sub 使用的 go micro Broker

	tokenCache unsafe.Pointer // *accessToken

}

const (
	// DefaultClusterRefreshTokenTopic  默认的Cluster刷新token的topic
	DefaultClusterRefreshTokenTopic = "com.jinmuhealth.topic.jinmu-wx-refresh-access-token"
)

// NewDefaultClusterAccessTokenServer 创建一个 DefaultClusterAccessTokenServer. 如果 httpClient == nil 则默认使用 util.DefaultHttpClient.
func NewDefaultClusterAccessTokenServer(appID, appSecret string, httpClient *http.Client, clusterRefreshTokenTopic string, bk broker.Broker) (srv *DefaultClusterAccessTokenServer) {
	if httpClient == nil {
		httpClient = util.DefaultHttpClient
	}

	if clusterRefreshTokenTopic == "" {
		clusterRefreshTokenTopic = DefaultClusterRefreshTokenTopic
	}

	if bk == nil {
		bk = broker.DefaultBroker
	}

	srv = &DefaultClusterAccessTokenServer{
		instanceID:                  uuid.New(),
		appID:                       appID,
		appSecret:                   appSecret,
		httpClient:                  httpClient,
		refreshTokenRequestChan:     make(chan string),
		refreshTokenResponseChan:    make(chan refreshTokenResult),
		clusterRefreshTokenTopic:    clusterRefreshTokenTopic,
		broker:                      bk,
		refreshTokenFormClusterChan: make(chan accessToken),
	}

	go srv.subscribeTopics()
	go srv.tokenUpdateDaemon(time.Hour * 24 * time.Duration(100+rand.Int63n(200)))
	return
}

// InstanceID 当前节点的 Instance ID
func (srv *DefaultClusterAccessTokenServer) InstanceID() string {
	return srv.instanceID.String()
}

// IID01332E16DF5011E5A9D5A4DB30FED8E1 接口标识, 没有实际意义
func (srv *DefaultClusterAccessTokenServer) IID01332E16DF5011E5A9D5A4DB30FED8E1() {}

// Token 请求中控服务器返回缓存的 access_token
func (srv *DefaultClusterAccessTokenServer) Token() (token string, err error) {
	if p := (*accessToken)(atomic.LoadPointer(&srv.tokenCache)); p != nil {
		return p.Token, nil
	}
	return srv.RefreshToken("")
}

type refreshTokenResult struct {
	token string
	err   error
}

// RefreshToken 请求中控服务器刷新 access_token
func (srv *DefaultClusterAccessTokenServer) RefreshToken(currentToken string) (token string, err error) {
	srv.refreshTokenRequestChan <- currentToken
	result := <-srv.refreshTokenResponseChan
	return result.token, result.err
}

func (srv *DefaultClusterAccessTokenServer) tokenUpdateDaemon(initTickDuration time.Duration) {
	tickDuration := initTickDuration

NEW_TICK_DURATION:
	ticker := time.NewTicker(tickDuration)
	for {
		select {
		// 响应主动刷新 Token 请求
		case currentToken := <-srv.refreshTokenRequestChan:
			accessToken, cached, err := srv.updateToken(currentToken)
			if err != nil {
				srv.refreshTokenResponseChan <- refreshTokenResult{err: err}
				break
			}
			srv.refreshTokenResponseChan <- refreshTokenResult{token: accessToken.Token}
			if !cached {
				tickDuration = time.Duration(accessToken.ExpiresIn) * time.Second
				ticker.Stop()
				goto NEW_TICK_DURATION
			}

		// 响应内部定时器刷新 Token
		case <-ticker.C:
			accessToken, _, err := srv.updateToken("")
			if err != nil {
				break
			}
			newTickDuration := time.Duration(accessToken.ExpiresIn) * time.Second
			if abs(tickDuration-newTickDuration) > time.Second*5 {
				tickDuration = newTickDuration
				ticker.Stop()
				goto NEW_TICK_DURATION
			}

		case remoteToken := <-srv.refreshTokenFormClusterChan:
			accessToken := srv.updateTokenFromCluster(remoteToken)

			// 由于是集群更新，为了避免本地定时器更新与集群主控节点更新时间间隔太近，因此增加时延时间偏移量
			//
			// 场景说明：
			//
			// 1) 当顺序启动集群中3个节点的时候，各自内部定时器更新计划时刻如下：
			//  节点1定时器: t1_0 -----> t1_1
			//  节点2定时器:   t2_0 -----> t2_1
			//  节点3定时器:     t3_0 -----> t3_1
			//
			// 2) 当某个节点首次被调用 Token() 试图获取一个 AccessToken时 (假设被调用节点为节点3)
			//    此时，节点3会从微信服务器获取一个新的 Token，并且在成功后，向集群 Publish 这个新获取的 Access Token 的 Topic
			//    其它节点在获取这个 Topic 消息后，将会触发更新定时器。更新结果如下：
			//  	节点1定时器: t1_0 ---------------> t1_1(被替换，=t3_1+latency)
			//  	节点2定时器:   t2_0 ------------> t2_1(被替换，=t3_1+latency)
			//  	节点3定时器:     t3_0 -----> t3_1 -----> t3_2
			//	  其中 latency 是时延长度, 介于 2-10秒之间。
			//	  时延长度不能超过10秒，因为微信网络请求延迟最小修正值为10秒，参见 updateToken 方法代码内容
			//
			// 3.1) 正常情况下，节点3会在t31时刻触发定时器主动更新 Access Token，并再次发送 Topic。其它节点也随之更新，结果如下，如此往复：
			//  	节点1定时器: t1_0 --------------------------> t1_2(被替换，=t3_2+latency)
			//  	节点2定时器:   t2_0 -----------------------> t2_2(被替换，=t3_2+latency)
			//  	节点3定时器:     t3_0 -----> t3_1 -----> t3_2
			//
			// 3.2) 异常情况下，比较节点3被关闭或崩溃，比如节点2的定时器优先触发:
			//  	节点1定时器: t10 ------------------------------------------> t1_3(被替换，=t2_3+latency)
			//  	节点2定时器:   t20 -----------------------> t22 -----> t23
			//  	节点3定时器:     t30 -----> t31 -----> t32 (失效)
			//
			tickDuration = time.Duration(accessToken.ExpiresIn+remoteLatency()) * time.Second
			ticker.Stop()
			goto NEW_TICK_DURATION
		}

	}
}

// remoteLatency 返回远程更新时延，随机在 [2, 10) 秒
func remoteLatency() int64 {
	return 2 + rand.Int63n(8)
}

func abs(x time.Duration) time.Duration {
	if x >= 0 {
		return x
	}
	return -x
}

type accessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`

	InstanceID string `json:"instance_id,omitempty"`
}

// updateToken 从微信服务器获取新的 access_token 并存入缓存, 同时返回该 access_token.
func (srv *DefaultClusterAccessTokenServer) updateToken(currentToken string) (token *accessToken, cached bool, err error) {
	if currentToken != "" {
		if p := (*accessToken)(atomic.LoadPointer(&srv.tokenCache)); p != nil && currentToken != p.Token {
			return p, true, nil // 无需更改 p.ExpiresIn 参数值, cached == true 时用不到
		}
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", url.QueryEscape(srv.appID), url.QueryEscape(srv.appSecret))

	httpResp, err := srv.httpClient.Get(url)
	if err != nil {
		atomic.StorePointer(&srv.tokenCache, nil)
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		atomic.StorePointer(&srv.tokenCache, nil)
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	var result struct {
		core.Error
		accessToken
	}

	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		atomic.StorePointer(&srv.tokenCache, nil)
		return
	}
	if result.ErrCode != core.ErrCodeOK {
		atomic.StorePointer(&srv.tokenCache, nil)
		err = &result.Error
		return
	}

	// 由于网络的延时, access_token 过期时间留有一个缓冲区
	switch {
	case result.ExpiresIn > 31556952: // 60*60*24*365.2425
		atomic.StorePointer(&srv.tokenCache, nil)
		err = errors.New("expires_in too large: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	case result.ExpiresIn > 60*60:
		result.ExpiresIn -= 60 * 10
	case result.ExpiresIn > 60*30:
		result.ExpiresIn -= 60 * 5
	case result.ExpiresIn > 60*5:
		result.ExpiresIn -= 60
	case result.ExpiresIn > 60:
		result.ExpiresIn -= 10
	default:
		atomic.StorePointer(&srv.tokenCache, nil)
		err = errors.New("expires_in too small: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	}

	result.accessToken.InstanceID = srv.InstanceID()

	tokenCopy := result.accessToken
	atomic.StorePointer(&srv.tokenCache, unsafe.Pointer(&tokenCopy))
	token = &tokenCopy

	// 向 Cluster 发布消息
	errpublishRefreshAccessTokenMessage := srv.publishRefreshAccessTokenMessage(tokenCopy)
	if errpublishRefreshAccessTokenMessage != nil {
		return nil, false, errpublishRefreshAccessTokenMessage
	}

	return
}

func (srv *DefaultClusterAccessTokenServer) updateTokenFromCluster(token accessToken) accessToken {
	tokenCopy := token
	atomic.StorePointer(&srv.tokenCache, unsafe.Pointer(&tokenCopy))

	return tokenCopy
}

func (srv *DefaultClusterAccessTokenServer) subscribeTopics() {
	fmt.Println("Subscribe topic", srv.clusterRefreshTokenTopic)
	_, err := srv.broker.Subscribe(srv.clusterRefreshTokenTopic, srv.processClusterRefreshTokenMessage)
	if err != nil {
		die(fmt.Errorf("failed to subscribe topic %s: %v", srv.clusterRefreshTokenTopic, err))
	}
}

// publishMessage 广播消息
func (srv *DefaultClusterAccessTokenServer) publishRefreshAccessTokenMessage(token accessToken) error {
	data, err := json.Marshal(&token)
	if err != nil {
		return err
	}
	msg := &broker.Message{
		Body: data,
	}

	if err := srv.broker.Publish(srv.clusterRefreshTokenTopic, msg); err != nil {
		return fmt.Errorf("failed publish message on %s topic: %v", srv.clusterRefreshTokenTopic, err)
	}

	return nil
}

// processClusterRefreshTokenMessage 处理从 Cluster 中收到的 AccessToken 更新
func (srv *DefaultClusterAccessTokenServer) processClusterRefreshTokenMessage(p broker.Event) error {
	var msg accessToken
	err := json.Unmarshal(p.Message().Body, &msg)
	if err != nil {
		fmt.Printf("Received remote token error: %v", err)
		return err
	}

	// 忽略来自自身节点的消息，不触发更新 chan
	if msg.InstanceID == srv.InstanceID() {
		return nil
	}

	srv.refreshTokenFormClusterChan <- msg
	return nil
}

func die(err error) {
	panic(err)
}
