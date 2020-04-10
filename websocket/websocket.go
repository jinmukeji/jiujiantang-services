package main

// FIX<E: websockrt 因为 iris 升级之后不兼容，需要修复

func main() {

}

// // Run first `go run main.go server`
// // and `go run main.go client` as many times as you want.
// // Originally written by: github.com/antlaw to describe an old issue.
// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	"syscall"
// 	"time"

// 	"github.com/jinmukeji/gf-api2/websocket/config"
// 	websocket "github.com/kataras/iris/v12/websocket"
// 	"github.com/kataras/iris/v12"
// 	"github.com/micro/cli"
// 	"github.com/micro/go-micro"
// 	"github.com/micro/go-micro/broker"
// )

// const (
// 	// CreatedUserTopic 创建用户的Topic
// 	CreatedUserTopic = "com.jinmuhealth.topic.jinmul-wx-created-user"
// )

// var (
// 	// ws websocket.Server
// 	ws *websocket.Server
// 	// port 端口
// 	port int
// )

// // Message 消息
// type Message struct {
// 	SceneID int `json:"scene_id"`
// 	UserID  int `json:"user_id"`
// }

// func subscribeTopics() {
// 	_, err := broker.Subscribe(CreatedUserTopic, SubscribeMessage)
// 	if err != nil {
// 		die(fmt.Errorf("sub error: %v", err))
// 	}
// }

// // SubscribeMessage 关注消息
// func SubscribeMessage(p broker.Event) error {

// 	var msg Message
// 	err := json.Unmarshal(p.Message().Body, &msg)
// 	if err != nil {
// 		return err
// 	}

// 	err = SendMessage(msg)
// 	if err != nil {
// 		fmt.Println("faild to send message:", string(p.Message().Body), err)
// 		return err
// 	}

// 	return nil
// }

// func main() {
// 	service := micro.NewService(
// 		micro.Name(config.FullServiceName()),
// 		micro.Flags(
// 			cli.IntFlag{
// 				Name:        "x_port",
// 				Value:       9100,
// 				Usage:       "WebSocket port",
// 				EnvVar:      "X_PORT",
// 				Destination: &port,
// 			},
// 			cli.BoolFlag{
// 				Name:  "version",
// 				Usage: "Show version information",
// 			},
// 		),
// 	)
// 	service.Init(
// 		micro.Action(func(c *cli.Context) {
// 			if c.Bool("version") {
// 				config.PrintFullVersionInfo()
// 				os.Exit(0)
// 			}
// 		}),
// 	)

// 	if err := broker.Init(); err != nil {
// 		die(fmt.Errorf("Broker Init error: %v", err))
// 	}
// 	if err := broker.Connect(); err != nil {
// 		die(fmt.Errorf("Broker Connect error: %v", err))
// 	}

// 	subscribeTopics()
// 	ServerLoop(port)

// 	// Go signal notification works by sending `os.Signal`
// 	// values on a channel. We'll create a channel to
// 	// receive these notifications (we'll also make one to
// 	// notify us when the program can exit).
// 	sigs := make(chan os.Signal, 1)
// 	done := make(chan bool, 1)

// 	// `signal.Notify` registers the given channel to
// 	// receive notifications of the specified signals.
// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

// 	// This goroutine executes a blocking receive for
// 	// signals. When it gets one it'll print it out
// 	// and then notify the program that it can finish.
// 	go func() {
// 		sig := <-sigs
// 		fmt.Println()
// 		fmt.Println(sig)
// 		done <- true
// 	}()

// 	// The program will wait here until it gets the
// 	// expected signal (as indicated by the goroutine
// 	// above sending a value on `done`) and then exit.
// 	fmt.Println("Running... ")
// 	fmt.Println("Press Ctrl+C to exit")
// 	<-done
// 	fmt.Println("bye!")

// }

// // server side

// // OnConnect handles incoming websocket connection
// func OnConnect(c websocket.Connection) {
// 	fmt.Println("socket.OnConnect()")
// 	c.On("join", func(message string) { OnJoin(message, c) })
// 	c.On("objectupdate", func(message string) { OnObjectUpdated(message, c) })
// 	c.OnDisconnect(func() { OnDisconnect(c) })
// 	c.OnMessage(func(msg []byte) {
// 		sceneID, err := strconv.Atoi(string(msg))
// 		if err != nil {
// 			fmt.Println("RECV MSG Error:", err)
// 		}

// 		addConnection(sceneID, c)
// 	})
// }

// var conns = NewIntStringMap()

// // addConnection 添加连接
// func addConnection(sceneID int, c websocket.Connection) {
// 	conns.Store(sceneID, c.ID())
// }

// // remoteConnection 移除连接
// func remoteConnection(sceneID int) {
// 	conns.Delete(sceneID)
// }

// // SendMessage 群发message
// func SendMessage(msg Message) error {
// 	connID, ok := conns.Load(msg.SceneID)
// 	if !ok {
// 		return fmt.Errorf("Connection ID is not found[%d]", msg.SceneID)
// 	}

// 	c := ws.GetConnection(connID)
// 	if c != nil {
// 		userID := strconv.Itoa(msg.UserID)
// 		if err := c.EmitMessage([]byte(userID)); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // ServerLoop listen and serve websocket requests
// func ServerLoop(p int) {
// 	app := iris.New().
// 		Configure(iris.WithRemoteAddrHeader("X-Forwarded-For"))

// 	ws = websocket.New(websocket.Config{})

// 	// register the server on an endpoint.
// 	// see the inline javascript code i the websockets.html, this endpoint is used to connect to the server.
// 	app.Get("/socket", ws.Handler())

// 	// Health Check
// 	app.Get("/check", func(ctx iris.Context) {
// 		_, errWriteString := ctx.WriteString("HEALTHY")
// 		if errWriteString != nil {
// 			return
// 		}
// 	})

// 	ws.OnConnection(OnConnect)
// 	errRun := app.Run(iris.Addr(fmt.Sprintf(":%d", p)))
// 	if errRun != nil {
// 		return
// 	}
// }

// // OnJoin handles Join broadcast group request
// func OnJoin(message string, c websocket.Connection) {
// 	t := time.Now()
// 	c.Join("server")
// 	fmt.Println("OnJoin() time taken:", time.Since(t))
// }

// // OnObjectUpdated broadcasts to all client an incoming message
// func OnObjectUpdated(message string, c websocket.Connection) {
// 	fmt.Println("OnObjectUpdated() invalid message format:")
// 	t := time.Now()
// 	s := strings.Split(message, ";")
// 	if len(s) != 3 {
// 		fmt.Println("OnObjectUpdated() invalid message format:" + message)
// 		return
// 	}
// 	serverID, _, objectID := s[0], s[1], s[2]
// 	err := c.To("server"+serverID).Emit("objectupdate", objectID)
// 	if err != nil {
// 		fmt.Println(err, "failed to broacast object")
// 		return
// 	}
// 	fmt.Println(fmt.Sprintf("OnObjectUpdated() message:%v, time taken: %v", message, time.Since(t)))
// }

// // OnDisconnect clean up things when a client is disconnected
// func OnDisconnect(c websocket.Connection) {
// 	for key, value := range conns.internal {
// 		if value == c.ID() {
// 			remoteConnection(key)
// 			break
// 		}
// 	}

// 	c.Leave("server")
// 	fmt.Println("OnDisconnect(): client disconnected!")
// }

// // die 有error会调用die
// func die(err error) {
// 	panic(err)
// }

// // IntStringMap 安全的map[int]string
// type IntStringMap struct {
// 	sync.RWMutex
// 	internal map[int]string // Key: Scene ID, Value: Connection ID
// }

// // NewIntStringMap 创建一个IntStringMap
// func NewIntStringMap() *IntStringMap {
// 	return &IntStringMap{
// 		internal: make(map[int]string),
// 	}
// }

// // Load 加载
// func (rm *IntStringMap) Load(key int) (value string, ok bool) {
// 	rm.RLock()
// 	result, ok := rm.internal[key]
// 	rm.RUnlock()
// 	return result, ok
// }

// // Delete 删除
// func (rm *IntStringMap) Delete(key int) {
// 	rm.Lock()
// 	delete(rm.internal, key)
// 	rm.Unlock()
// }

// // Store 存储
// func (rm *IntStringMap) Store(key int, value string) {
// 	rm.Lock()
// 	rm.internal[key] = value
// 	rm.Unlock()
// }
