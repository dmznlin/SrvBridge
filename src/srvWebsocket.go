/*Package srvbridge ***********************************************************
  作者: dmzn@163.com 2022-06-14 18:36:02
  描述: 提供web socket接口
******************************************************************************/
package srvbridge

import (
	"errors"
	"fmt"
	. "github.com/dmznlin/znlib-go/znlib"
	inifile "github.com/go-ini/ini"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strconv"
	"sync"
)

var (
	websocket_init sync.Once
	//wsWorker websocket服务对象
	websocketWorker srvWebSocket

	//websocketServer ws服务对象
	websocketServer *http.Server = nil
)

func init() {
	RegistWorker(websocketWorker)
}

func init_websocket() {
	websocket_init.Do(func() {
		if len(Application.HostIP) > 0 {
			websocketWorker.srvIP = Application.HostIP[0]
		} else {
			websocketWorker.srvIP = ""
		}

		websocketWorker.enable = true
		websocketWorker.srvPort = 8081
		websocketWorker.webPath = "/srv"
		websocketWorker.loadConfig() //apply config file
	})
}

//--------------------------------------------------------------------------------

type srvWebSocket struct {
	enable  bool
	srvIP   string
	srvPort int
	srvAddr string
	webPath string
}

func (ws *srvWebSocket) loadConfig() {
	if FileExists(Application.ConfigFile, false) {
		ini, err := inifile.Load(Application.ConfigFile)
		if err != nil {
			Warn("websocket.loadconfig: " + err.Error())
			return
		}

		sec := ini.Section("websocket")
		ws.enable = sec.Key("enable").In("true", []string{"true", "false"}) == "true"
		//是否启动服务

		vs := StrTrim(sec.Key("srvIP").String())
		if vs != "" {
			ws.srvIP = vs
		}

		vs = StrTrim(sec.Key("path").String())
		if vs != "" {
			ws.webPath = vs
		}

		vi := sec.Key("srvPort").MustInt(0)
		if vi > 0 {
			ws.srvPort = vi
		}
		ws.srvAddr = net.JoinHostPort(ws.srvIP, strconv.Itoa(ws.srvPort))
	}
}

func (ws srvWebSocket) WorkName() string {
	return "websocket-service"
}

func (ws srvWebSocket) Start() {
	init_websocket()
	if !websocketWorker.enable {
		return //不启动服务
	}

	go wsHub.run()
	http.HandleFunc(websocketWorker.webPath, wsHandle) //将请求交给wsHandle处理

	websocketServer = &http.Server{Addr: websocketWorker.srvAddr, Handler: nil}
	Info("websocket service on: " + websocketWorker.srvAddr)
	websocketServer.ListenAndServe()
}

func (ws srvWebSocket) Stop() (err error) {
	defer ErrorHandle(false, func(e any) {
		if e != nil {
			err = errors.New(fmt.Sprintf("stop %s failure.", ws.WorkName()))
		}
	})

	if websocketServer != nil {
		websocketServer.Shutdown(ServiceContext)
		Info(ws.WorkName() + "closed.")
	}
	return nil
}

//将普通的http连接升级为websocket连接
var wsUpgrader = &websocket.Upgrader{
	//定义读写缓冲区大小
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	//校验请求
	CheckOrigin: func(r *http.Request) bool {
		//如果不是get请求，返回错误
		if r.Method != "GET" {
			Info(fmt.Sprintf("%s: Host [%s] 请求方式错误", websocketWorker.WorkName(), r.Host))
			return false
		}
		//如果路径中不包括chat，返回错误
		if r.URL.Path != websocketWorker.webPath {
			Info(fmt.Sprintf("%s: Host [%s] 请求路径错误", websocketWorker.WorkName(), r.Host))
			return false
		}
		//还可以根据其他需求定制校验规则
		return true
	},
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	//通过升级后的升级器得到链接
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		Info(fmt.Sprintf("%s: Host [%s] 获取连接失败[%s]", websocketWorker.WorkName(), r.Host, err.Error()))
		return
	}

	//连接成功后注册用户
	client := &wsClient{
		conn: conn,
		msg:  make(chan []byte),
	}

	wsHub.register <- client
	defer func() {
		wsHub.unregister <- client
	}()

	//读写数据
	go read(client)
	write(client)
}

//--------------------------------------------------------------------------------

//websocket处理器
type websocketHub struct {
	clientList map[*wsClient]bool //客户端列表
	register   chan *wsClient     //客户端注册
	unregister chan *wsClient     //客户端注销
	broadcast  chan []byte        //待广播数据
}

//初始化处理中心
var wsHub = &websocketHub{
	clientList: make(map[*wsClient]bool),
	register:   make(chan *wsClient),
	unregister: make(chan *wsClient),
	broadcast:  make(chan []byte),
}

//处理中心处理获取到的信息
func (hub *websocketHub) run() {
loop:
	for {
		select {
		case <-ServiceContext.Done(): //退出服务
			break loop
		case client := <-hub.register: //注册客户端
			hub.clientList[client] = true
		case client := <-hub.unregister: //清理客户端
			if _, ok := hub.clientList[client]; ok {
				delete(hub.clientList, client)
			}
		case data := <-hub.broadcast: //广播数据
			for client := range hub.clientList {
				select {
				case client.msg <- data:
				default:
					delete(hub.clientList, client)
					close(client.msg)
				}
			}
		}
	}
}

//--------------------------------------------------------------------------------

//定义一个websocket连接对象，连接中包含每个连接的信息
type wsClient struct {
	conn *websocket.Conn
	msg  chan []byte
}

func read(user *wsClient) {
	//从连接中循环读取信息
	for {
		_, msg, err := user.conn.ReadMessage()
		if err != nil {
			Info(fmt.Sprintf("%s: Host [%s] 用户退出[%s]", websocketWorker.WorkName(), user.conn.RemoteAddr().String(), err.Error()))
			wsHub.unregister <- user
			break
		}
		//将读取到的信息传入websocket处理器中的broadcast中，
		wsHub.broadcast <- msg
	}
}

func write(user *wsClient) {
	for data := range user.msg {
		err := user.conn.WriteMessage(1, data)
		if err != nil {
			Info(fmt.Sprintf("%s: Host [%s] 写入错误[%s]", websocketWorker.WorkName(), user.conn.RemoteAddr().String(), err.Error()))
			break
		}
	}
}
