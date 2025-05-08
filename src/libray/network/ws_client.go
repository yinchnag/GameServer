package network

import (
	"net/http"
	"time"

	"wgame_server/libray/core"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	Url  string
	Conn *websocket.Conn
}

func (that *WsClient) Init() {
}

func (that *WsClient) Start(url string) {
	that.connectServer(url)
	// session := NewWsSession(that.Conn, nil)
	// go session.runSendLoop()
}

func (that *WsClient) connectServer(url string) {
	dialer := websocket.DefaultDialer
	header := make(http.Header)
	header.Add("Origin", "http://localhost/")
	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		return
	}
	conn.SetPongHandler(func(appData string) error {
		core.Logger.Infof("心跳链接 msg:%s", appData)
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})
	that.Conn = conn
}
