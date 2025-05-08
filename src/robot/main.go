package main

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	connectServer()
}

func connectServer() {
	dialer := websocket.DefaultDialer
	header := make(http.Header)
	header.Add("Origin", "http://localhost/")
	conn, _, err := dialer.Dial("ws://127.0.0.1:8000", header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close() // 关闭连接
	time.Sleep(time.Second * 1)
}
