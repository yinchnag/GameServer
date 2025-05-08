package network

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"wgame_server/libray/core"

	"github.com/gorilla/websocket"
)

var (
	wsServer     *WsServer // 单例
	wsServerOnce sync.Once // 单次加锁
)

// 单例(多线程安全)
func GetWsServer() *WsServer {
	wsServerOnce.Do(func() {
		wsServer = &WsServer{}
	})
	return wsServer
}

type WsServer struct {
	Addr           string // 监听路径
	MaxConnections int
	ClientCount    int32               // 客户端数量
	protoMap       map[int]ILogicProto // 逻辑协议映射
	Upgrader       *websocket.Upgrader // 升级WS
	Server         *http.Server        // 服务
	alive          int32               // 是否运行中
	rejectConn     int32               // 是否拒绝连接

	storage  *WsSessionStorage // 会话存储
	sessions sync.Map          // 会话列表
}

func (that *WsServer) Init() {
	that.storage = newWsSessionStorage()
	that.protoMap = make(map[int]ILogicProto)
	for _, proto := range C2S_MSG_MAP {
		if proto != nil {
			protoID := proto.GetProtoID()
			if _, ok := that.protoMap[protoID]; !ok {
				that.protoMap[protoID] = proto
			} else {
				core.Logger.Errorf("协议冲突 %d", protoID)
			}
		}
	}
}

func (that *WsServer) Start(addr string, maxConn int) {
	if atomic.LoadInt32(&that.alive) != 0 { // 不具备启动服务的条件
		return
	}

	atomic.StoreInt32(&that.alive, 1)
	that.Addr = addr
	that.MaxConnections = maxConn
	that.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	that.Server = &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	that.Server.RegisterOnShutdown(func() {
		fmt.Print("WsServer Closed")
	})
	http.HandleFunc("/", that.ServerHTTP)
	err := that.Server.ListenAndServe()
	if err != nil {
		core.Logger.Error(err)
	}
}

func (that *WsServer) ServerHTTP(writer http.ResponseWriter, req *http.Request) {
	conn, err := that.Upgrader.Upgrade(writer, req, nil)
	if err != nil {
		return
	}
	if !that.IsAllowConn() {
		conn.Close() // 暂时中断服务
		return
	}
	atomic.AddInt32(&that.ClientCount, 1)
	sessionClient := &WsConnectSession{conn: conn, server: that}
	that.OnSessionConnecting(sessionClient)
	if sessionId := that.storage.RegisterRecv(sessionClient); sessionId != 0 {
		sessionClient.Init()
	}
}

func (that *WsServer) OnSessionConnecting(client *WsConnectSession) {
}

// 是否允许连接
func (that *WsServer) IsAllowConn() bool {
	return atomic.LoadInt32(&that.rejectConn) == 0
}

func (that *WsServer) OnSessionClosing(client *WsConnectSession) {
	that.ResetSession(client.GetSessionID())
}

func (that *WsServer) OnSessionReceived(client *WsConnectSession, msgType int, msgData []byte) {
	if !that.IsAllowConn() {
		that.ResetSession(client.GetSessionID())
		client.Close("not allow connection")
		return
	}
	headerSize := PACK_SIZE + MSGID_SIZE
	if len(msgData) < headerSize {
		core.Logger.Warnln("invalid messsage size", len(msgData))
		return
	}
	_, msgid, err := HF_DecodeMsgPB(msgData)
	if err != nil {
		core.Logger.Warnf("消息处理出错: %v", err)
		return
	}

	that.OnPlayerMessage(client, msgid, msgData[headerSize:])
}

func (that *WsServer) OnPlayerMessage(client *WsConnectSession, msgid int, msgData []byte) {
	proto, ok := that.protoMap[msgid]
	if !ok {
		core.Logger.Warnf("not found proto %d", msgid)
		return
	}
	meta := proto.GetMsg()
	proto.HandleProtocol(client, msgData, msgid, meta)
}

func (that *WsServer) ResetSession(sessionID uint32) {
	session := that.GetSession(sessionID)
	if session != nil {
		user := session.GetPlayer()
		if !core.IsNil(user) {
			user.SetOnline(false)
		}
		that.sessions.Delete(sessionID)
		session.sessionID = 0
	}
}

func (that *WsServer) GetSession(sessionID uint32) *WsSession {
	client, ok := that.sessions.Load(sessionID)
	if !ok {
		return nil
	}
	return client.(*WsSession)
}

// 设置会话
func (that *WsServer) SetSession(sessionID uint32, session *WsSession) {
	if session.sessionID != 0 && session.sessionID != sessionID {
		that.sessions.Delete(session.sessionID)
	}
	session.sessionID = sessionID
	if sessionID != 0 {
		that.sessions.Store(sessionID, session)
	}
}

// 发送裸数据
func (that *WsServer) SendBytes(sessionID uint32, body []byte) error {
	client := that.storage.FindSessionByID(sessionID)
	if client == nil || client.conn == nil {
		return fmt.Errorf("session not found")
	}
	return client.SendBytes(body)
}

// 获取连接客户端
func (that *WsServer) GetSessionClient(sessionID uint32) *WsConnectSession {
	return that.storage.FindSessionByID(sessionID)
}

// 设置运行中
func (that *WsServer) SetAlive(val bool) {
	if val {
		atomic.StoreInt32(&that.alive, 1)
	} else {
		atomic.StoreInt32(&that.alive, 0)
	}
}

// 是否运行中
func (that *WsServer) IsAlive() bool {
	return atomic.LoadInt32(&that.alive) == 1
}

// 设置是否允许连接
func (that *WsServer) SetAllowConn(val bool) {
	if val {
		atomic.StoreInt32(&that.rejectConn, 0)
	} else {
		atomic.StoreInt32(&that.rejectConn, 1)
	}
}

func (that *WsServer) Shutdown(ctx context.Context) {
	if !that.IsAlive() {
		core.Logger.Error("server is not alive")
		return
	}
	that.SetAlive(false)
	that.SetAllowConn(false)
	that.storage.ForEach(func(client *WsConnectSession) {
		client.Close("server shutdown")
	})
	if that.Server != nil {
		that.Server.Shutdown(ctx)
		that.Server = nil
	}
}
