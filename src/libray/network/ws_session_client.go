package network

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"wgame_server/libray/core"

	"github.com/gorilla/websocket"
)

var _ core.ISocketConnect = (*WsConnectSession)(nil)

type WsConnectSession struct {
	sessionId uint32          // 会话ID
	conn      *websocket.Conn // 套接字连接
	server    *WsServer       // 服务器
	destory   int32           // 是否销毁 0未销毁 1销毁
	clientIp  string          // 来源ip地址
	addr      string          // 监听地址
	sync.Mutex
}

// 获取会话ID
func (that *WsConnectSession) GetSessionID() uint32 {
	return that.sessionId
}

// 来源IP地址
func (that *WsConnectSession) GetClientIp() string {
	return that.clientIp
}

func (that *WsConnectSession) Init() {
	that.addr = that.conn.RemoteAddr().String()
	that.clientIp = HF_GetWsConnIP(that.conn)
	that.conn.SetReadLimit(MAX_MESSAGE_SIZE)
	that.conn.SetReadDeadline(time.Now().Add(PONG_WAIT)) // 心跳间隔
	that.conn.SetPingHandler(func(appData string) error {
		that.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return that.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})
	that.runRecvLoop()
}
func (that *WsConnectSession) DeferClose(src string) {}
func (that *WsConnectSession) Close(src string) {
	if err := recover(); err != nil {
		core.Logger.Warnf("%s拦截错误:%s", src, err)
	}
	if atomic.LoadInt32(&that.destory) != 1 {
		atomic.StoreInt32(&that.destory, 1)
		that.server.OnSessionClosing(that)
		that.server.storage.Unregister(that)
		that.conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(time.Second))
		that.conn.Close()
	}
}
func (that *WsConnectSession) ShutDown() {}

func (that *WsConnectSession) GetPlayer() core.IPlayer {
	session := that.server.GetSession(that.sessionId)
	if session == nil {
		return nil
	}
	return session.GetPlayer()
}

func (that *WsConnectSession) runRecvLoop() {
	defer that.Close("runRecvLoop")
	for {
		msgType, msgData, err := that.conn.ReadMessage()
		if err != nil {
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				continue
			}
			that.Close("normal")
			break
		}
		that.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
		that.server.OnSessionReceived(that, msgType, msgData)
	}
}

func (that *WsConnectSession) SendBytes(body []byte) error {
	that.Lock()
	defer that.Unlock()
	that.conn.SetWriteDeadline(time.Now().Add(SEND_WAIT)) // 发送间隔
	err := that.conn.WriteMessage(websocket.BinaryMessage, body)
	if err != nil {
		core.Logger.Debugf("SendBytes err: %v", err)
	}
	return err
}

func (that *WsConnectSession) SendError(errCode int) {
}

func (that *WsConnectSession) OnMessage(msgData []byte) bool {
	return true
}

func (that *WsConnectSession) GetSocketRecv() core.ISocketRecv {
	return that.server.GetSession(that.sessionId)
}
