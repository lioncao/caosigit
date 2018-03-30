package network

import (
	"lioncao/util/tools"
	"net"
)

type Listener struct {
	id            int32
	l             *net.TCPListener
	index         int32
	connectMap    map[int32]*Connector
	MsgHandlerMgr ConnectHandler
}

func (this *Listener) Init(id int32, addr string, m ConnectHandler) {
	this.index = 0
	this.connectMap = make(map[int32]*Connector)
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", addr)
	this.l, _ = net.ListenTCP("tcp", tcpAddr)
	this.MsgHandlerMgr = m
}

func (this *Listener) Start() {
	for {
		conn, err := this.l.AcceptTCP()
		if err != nil {
			continue
		}
		c := new(Connector)
		this.index++
		c.Init(conn, this.index, this)
		this.connectMap[this.index] = c
		c.RegisterHanlder(this.MsgHandlerMgr)
		go c.StartHandler()
	}
}

func (this *Listener) Disconnect(id int32) {
	tools.GetLog().Log("connect id:%d disconnect", id)
	delete(this.connectMap, id)
}
