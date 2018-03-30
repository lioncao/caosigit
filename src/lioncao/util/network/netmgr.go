package network

import (
	"fmt"
	"lioncao/util/tools"
)

type Netmgr struct {
	listenAddrs map[int32]string
	listenMap   map[int32]*Listener
	connectMap  map[int32]*Connector
	HandlerMgr  map[int32]ConnectHandler
}

func (this *Netmgr) Init() {
	this.listenAddrs = make(map[int32]string)
	this.listenMap = make(map[int32]*Listener)
	this.connectMap = make(map[int32]*Connector)
	this.HandlerMgr = make(map[int32]ConnectHandler)
}

//开始网络层处理
func (this *Netmgr) Start() {
	//开始监听
	for listenId, addr := range this.listenAddrs {
		this.StartListen(listenId, addr)
	}

	for _, c := range this.connectMap {
		go c.StartHandler()
	}
}

func (this *Netmgr) DisConnect(id int32) {
	tools.GetLog().Log("zudong connect dis:%d", id)
	delete(this.connectMap, id)
}

func (this *Netmgr) AddListener(id int32, ip, port string, mgr ConnectHandler) {
	addr := fmt.Sprintf("%s:%s", ip, port)
	this.listenAddrs[id] = addr
	this.HandlerMgr[id] = mgr

}

func (this *Netmgr) AddConnector(id int32, ip, port string, mgr ConnectHandler) {
	c := new(Connector)
	c.Connect(id, ip, port, this, mgr)
	this.connectMap[id] = c
}

//开始监听
func (this *Netmgr) StartListen(id int32, addr string) {
	l := new(Listener)
	l.Init(id, addr, this.HandlerMgr[id])
	go l.Start()
	this.listenMap[id] = l
}

func (this *Netmgr) GetConnector(id int32) *Connector {
	c, ok := this.connectMap[id]
	if ok {
		return c
	} else {
		return nil
	}
}
