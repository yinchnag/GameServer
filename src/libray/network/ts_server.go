package network

import (
	"net"

	"wgame_server/libray/manager"
)

type TsServer struct {
	maxConn int
	listen  net.Listener
	Alive   int32 // 是否运行中
}

func (that *TsServer) Init() {
}

func (that *TsServer) Start(addr string, maxConn int) {
	that.maxConn = maxConn
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	that.listen = listen
}

func (that *TsServer) Accept() {
	for {
		conn, err := that.listen.Accept()
		if err != nil {
			continue
		}

		manager.GetSignalManager().Notify(manager.SIGNAL_CONNECT_SERVER, conn)
	}
}
