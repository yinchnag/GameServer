package network

import (
	"fmt"
	"net"
	"reflect"
	"sync/atomic"
	"time"

	"wgame_server/libray/core"
	"wgame_server/libray/define/PB"
	"wgame_server/libray/manager"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const (
	PONG_WAIT = 60 * time.Second

	SOCKET_TRY_NUM = 10              // 网络消息重发次数
	SEND_WAIT      = 5 * time.Second // 发送超时时间

	SEND_CHAN_SIZE   = 2000
	RECV_CHAN_SIZE   = 2000
	MAX_MESSAGE_SIZE = 1024
)

var _ core.ISocketRecv = (*WsSession)(nil)

type WsSession struct {
	sessionID uint32          // 会话ID
	Conn      *websocket.Conn // 套接字连接
	server    *WsServer

	SendChan     chan []byte // 服务器发送消息管道
	tryNum       int         // 重试次数
	ShutDownFlag int32       // 是否关闭 0未关闭 1关闭
	ShutTime     int64       // 关闭时间

	TryNum   int                 //
	protoMap map[int]ILogicProto // 协议列表
	user     core.IPlayer        // 玩家对象
	clientIp string              // 来源ip地址
}

// 是否允许连接
func (that *WsSession) isAllowConn() bool {
	if that.server == nil {
		return false
	}
	if atomic.LoadInt32(&that.ShutDownFlag) == 1 {
		return false
	}
	if !that.server.IsAllowConn() {
		return false
	}
	return true
}

func (that *WsSession) Init(user core.IPlayer, protos []ILogicProto) {
	that.server = GetWsServer()
	that.user = user
	for _, v := range protos {
		that.protoMap[v.GetProtoID()] = reflect.New(reflect.ValueOf(v).Elem().Type()).Interface().(ILogicProto)
		// that.protoMap[v.GetProtoID()].SetSession(that)
	}
	that.SendChan = make(chan []byte, SEND_CHAN_SIZE)

	go that.runSendLoop()
	manager.GetSignalManager().Notify(manager.SIGNAL_CONNECT_SERVER, that)
}

func (that *WsSession) runSendLoop() {
	for msg := range that.SendChan {
		if !that.isAllowConn() {
			break
		}
		for {
			err := that.server.SendBytes(that.sessionID, msg)
			if err == nil {
				_, msgid, _ := HF_DecodeMsgPB(msg)
				core.Logger.Debugf("Socket SendBytes: msgid=%d", msgid)
				that.tryNum = 0
				break
			}
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				that.tryNum += 1
				if that.tryNum < SOCKET_TRY_NUM {
					continue
				}
			}
			client := that.server.GetSessionClient(that.sessionID)
			if client != nil {
				client.Close("send_error")
			}
			break
		}
	}
}

func (that *WsSession) SendByte(msgId int, msg []byte) {
	if that.SendChan == nil {
		return
	}
	if len(that.SendChan) >= SEND_CHAN_SIZE {
		core.Logger.Warnln("send chan is full")
		return
	}
	that.SendChan <- msg
}

// 发送错误消息
func (that *WsSession) SendError(errCode int) {
	if that.SendChan == nil || that.server == nil {
		return
	}
	if !that.isAllowConn() {
		return
	}
	msg := &PB.S2C_Player_Error{}
	msg.ErrorCode = uint32(errCode)
	msgData, err := proto.Marshal(msg)
	if err != nil {
		core.Logger.Errorln("SendError: marshaling error: ", err)
		return
	}
	that.SendBytes(HF_EncodeMsgPB(uint16(PB.S2C_Player_s2c_player_error), msgData))
}

// 发送业务消息
func (that *WsSession) SendBytes(body []byte) error {
	if that.SendChan == nil || !that.isAllowConn() {
		return fmt.Errorf("SendBytes: SendChan is nil or not allow conn")
	}
	if len(that.SendChan) >= SEND_CHAN_SIZE-100 {
		core.Logger.Errorf("SendBytes: Chan overflow from %s", that.Conn.RemoteAddr())
		that.SendChan <- []byte("")
		that.ShutDown()
		that.ShutTime = time.Now().Unix()
		return fmt.Errorf("SendBytes: Chan overflow")
	}
	that.SendChan <- body
	return nil
}

// 关闭
func (that *WsSession) ShutDown() {
	atomic.StoreInt32(&that.ShutDownFlag, 1)
	that.SendChan <- []byte{}
}

func (that *WsSession) DeferClose(src string) {
}

func (that *WsSession) GetPlayer() core.IPlayer {
	return that.user
}

// 获取IP地址
func (that *WsSession) GetClientIp() string {
	return that.clientIp
}

func (that *WsSession) OnMessage(msgData []byte) bool {
	_, msgId, err := HF_DecodeMsgPB(msgData)
	if err != nil {
		core.Logger.Warnf("消息处理错误: %v", err)
		return false
	}
	proto, ok := that.protoMap[msgId]
	if !ok {
		core.Logger.Warnf("没有找到对应的协议: %v", msgId)
		return false
	}
	proto.HandleProtocol(that, msgData, msgId, proto.GetMsg())
	return false
}

func (that *WsSession) Close(str string) {
}

func (that *WsSession) SetSessionId(sessionID uint32) {
	that.server.SetSession(sessionID, that)
}
