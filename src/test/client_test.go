package test

import (
	"log"
	"net/http"
	"testing"
	"time"

	"wgame_server/libray/core"
	"wgame_server/libray/define/PB"
	"wgame_server/libray/network"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func TestClientConnect(t *testing.T) {
	dia := websocket.DefaultDialer
	header := make(http.Header)
	header.Add("Origin", "http://localhost/")
	conn, _, err := dia.Dial("ws://127.0.0.1:8000", header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	conn.SetPongHandler(func(appData string) error {
		core.Logger.Infof("心跳链接 msg:%s", appData)
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})
	// go func() {
	// 	for {
	// 		msgType, msgData, err := conn.ReadMessage()
	// 		if err != nil {
	// 			log.Println("read:", err)
	// 			continue
	// 		}
	// 		log.Printf("recv: %d-%s", msgType, msgData)
	// 	}
	// }()
	msg := &PB.C2S_Player_Login{}
	msgData, _ := proto.Marshal(msg)
	data := network.HF_EncodeMsgPB(1, msgData)
	conn.WriteMessage(websocket.BinaryMessage, data)
	time.Sleep(time.Second * 6)
}
