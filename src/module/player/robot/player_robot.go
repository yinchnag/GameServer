package robot

import (
	"net"
	"net/http"
	"time"

	"wgame_server/libray/core"
	"wgame_server/libray/module"

	"github.com/gorilla/websocket"
)

type PlayerRobot struct {
	TryNum   int
	SendChan chan []byte // 服务器发送消息管道
	conn     *websocket.Conn
	module.ModObj
}

func (that *PlayerRobot) Init(host interface{}) module.IModule {
	that.conn = connectServer()
	that.SendChan = make(chan []byte)
	if that.conn != nil {
		go that.runSendLoop()
	}
	that.ModObj.Init(that)
	that.SetInvokerAll(that)
	return that
}

func (that *PlayerRobot) OnDestory() {
	that.conn.Close() // 关闭连接
}

func connectServer() *websocket.Conn {
	dialer := websocket.DefaultDialer
	header := make(http.Header)
	header.Add("Origin", "http://localhost/")
	conn, _, err := dialer.Dial("ws://127.0.0.1:8000", header)
	if err != nil {
		return nil
	}
	conn.SetPongHandler(func(appData string) error {
		core.Logger.Infof("心跳链接 msg:%s", appData)
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})
	return conn
}

func (that *PlayerRobot) runSendLoop() {
	for msg := range that.SendChan {
		exit := false
		for {
			if that.conn == nil {
				exit = true
				break
			}
			that.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			err := that.conn.WriteMessage(websocket.BinaryMessage, msg)
			if err == nil {
				that.TryNum = 0
				break
			}
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				that.TryNum += 1
				if that.TryNum < 10 {
					continue
				}
			}
			exit = true
			break
		}
		if exit {
			break
		}
	}
}
