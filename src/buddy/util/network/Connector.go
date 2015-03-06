package network

import (
	"buddy/util/tools"
	"code.google.com/p/goprotobuf/proto"
	"common/msg"
	"fmt"
	"net"
	"time"
)

type ConnectHandler interface {
	OnClose(sid int32)
	Init()
	GetMessageHandler(cmd int32) (func(c *Connector, data *msg.S2S) error, bool)
}

type Connector struct {
	id            int32
	sid           int32
	conn          *net.TCPConn
	buffer        LockBuffer
	owner         *Listener
	mgr           *Netmgr
	MsgHandlerMgr ConnectHandler
	status        int32
	t             *time.Timer
	tickTime      int64
}

func (this *Connector) SetSid(Id int32) {
	this.sid = Id
}

func (this *Connector) Init(c *net.TCPConn, id int32, l *Listener) {
	this.id = id
	this.conn = c
	this.buffer.Init(1024, 1024)
	this.owner = l
	this.mgr = nil
	this.status = 0
	//BagMountHandlerMgr(this.MsgHandlerMgr)

}

func (this *Connector) RegisterHanlder(h ConnectHandler) {
	this.MsgHandlerMgr = h
}

func (this *Connector) Send(data []byte) error {
	return this.buffer.PushOutputData(data)
}

func (this *Connector) Read() {
	request := make([]byte, 1024)
	for {
		if this.status == 1 {
			break
		}
		read_len, err := this.conn.Read(request)
		if err != nil {
			tools.GetLog().LogError("read err:%s", err)
			this.status = 1
			break
		}

		if read_len == 0 { // 在gprs时数据不能通过这个判断是否断开连接,要通过心跳包
			tools.GetLog().LogError("read len is 0")
			this.status = 1
			break
		} else {
			// request[:read_len]处理
			err = this.buffer.PushInputData(request[:read_len])
			if err != nil {
				tools.GetLog().LogError("this.buffer.PushInputData() err")
				this.status = 1
				break
			}
		}
	}
}

func (this *Connector) Connect(id int32, ip, port string, m *Netmgr, h ConnectHandler) {
	this.id = id
	this.buffer.Init(1024, 1024)
	this.mgr = m
	this.owner = nil
	this.MsgHandlerMgr = h
	addr := fmt.Sprintf("%s:%s", ip, port)
	//BagMountHandlerMgr(this.MsgHandlerMgr)
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", addr)
	this.conn, _ = net.DialTCP("tcp", nil, tcpAddr)
}

func (this *Connector) StartHandler() {
	go this.Read()
	this.tickTime = time.Now().Unix()
	this.t = time.AfterFunc(5*time.Second, this.Ping)
	for {
		if this.status == 1 {
			break
		}

		now := time.Now().Unix()
		if now-this.tickTime > 7 {
			tools.GetLog().LogError("heart time out id:%d", this.id)
			this.status = 1
			break
		}

		//发送数据
		for {
			data, _ := this.buffer.PopOutputData()
			if data == nil {
				break
			}
			this.conn.Write(data)
		}

		data1, err1 := this.buffer.PopInputData()
		if err1 != nil {
			tools.GetLog().LogError("this.buffer.PopInputData() err:%s", err1)
			this.status = 1
			break
		}
		if data1 == nil {
			continue
		}
		var tempp msg.S2S
		err1 = proto.Unmarshal([]byte(data1), &tempp)
		if err1 != nil {
			tools.GetLog().LogError("proto.Unmarshal failed err:%s", err1)
			this.status = 1
			break
		}
		Cmd := tempp.GetCmd()
		if Cmd == int32(msg.MSG_GS_PING) {
			this.tickTime = time.Now().Unix()
			tools.GetLog().Log("hello Heart id:%d", this.id)
		} else {
			handler, ok := this.MsgHandlerMgr.GetMessageHandler(Cmd)
			if ok {
				handler(this, &tempp)
			} else {
				tools.GetLog().LogError("process cmd err:%s", err1)
			}
		}
	}
	//断开连接
	this.t.Stop()
	this.Close()
}

func (this *Connector) Ping() {
	now := time.Now().Unix()
	var temp msg.S2S
	temp = msg.S2S{
		Cmd:      proto.Int32(int32(msg.MSG_GS_PING)),
		TickTime: proto.Int64(now),
	}
	buffer, _ := proto.Marshal(&temp)
	this.Send(buffer)
	time.AfterFunc(time.Second*5, this.Ping)
}

func (this *Connector) Close() {
	defer this.conn.Close()
	if this.owner != nil {
		this.MsgHandlerMgr.OnClose(this.sid)
		this.owner.Disconnect(this.id)
	} else {
		this.mgr.DisConnect(this.id)
	}
}
