/*
	socket消息处理句柄,目前尝试的处理流程为:
	1. 网络连接获取数据
	2. 将数据强行规定为 length(int32) +  data 的形式
	3. 将消息交给doMsg函数处理
	4. 返回消息
*/

package socket

import (
	"buddy/util/tools"
	"net"
	"time"
)

const (
	DEFAULT_RW_BUFF_SIZE = 1024 * 1024 // 默认的消息读写buff长度  1M
	TIMEOUT_ERROR_INFO   = "i/o timeout"
)

// 回调函数相关设置
// const (
// 	CB_ON_RECEIVE     = "onReceiveCallback"
// 	CB_ON_DO_MESSAGE = "onDoMessageCallback"
// 	CB_ON_WRITE      = "onWriteCallback"
// )

type OnReceiveFunc func(handler *MessageHandler, buf []byte, bufSize int) error
type DoMessageFunc func(handler *MessageHandler) error
type OnWriteFunc func(handler *MessageHandler) error

type TcpCallbacks interface {
	SetSessionData(data interface{})
	GetSessionData() interface{}
	TcpOnReceiveCallback(handler *MessageHandler, buf []byte, bufSize int) error // 网络数据到达后的回调
	TcpDoMessageCallback(handler *MessageHandler) error                          // 处理函数
	TcpOnWriteCallback(handler *MessageHandler) error                            // 数据写入回调
}

func CheckTimeout(err error) bool {
	if err != nil {
		e, ok := err.(net.Error)
		if ok && e.Timeout() {
			return true
		}
	}
	return false
}

// 读写buffer数据结构体
type MessageHandler struct {
	handlerType string // 消息句柄的连接类型(http/tcp/udp)
	tcpAddr     string // 连接地址
	httpAddr    string
	needClose   bool

	// 逻辑处理函数句柄
	Conn      net.Conn // 连接接口
	callbacks TcpCallbacks

	// 缓存的数据句柄
	SessionData interface{}
}

// 消息处理句柄初始化
func (this *MessageHandler) Init(conn net.Conn) error {

	if conn == nil {
		return tools.Error("Init MessageHandler err conn=%v", conn)
	}

	// 句柄设置
	this.Conn = conn
	return nil
}

func (this *MessageHandler) SetCallback(callbacks TcpCallbacks) {

	if callbacks == nil {
		tools.ShowError("MessageHandler.SetCallback err: callbacks=", callbacks)
		return
	}
	this.callbacks = callbacks
}

//// 消息数据到达
//func (this *MessageHandler) OnReceived(buf []byte, dataSize int) error {
//	if this.onReceiveCallback != nil {
//		this.onReceiveCallback(buf, dataSize)
//	}
//	return nil
//}

//// 消息处理逻辑
//func (this *MessageHandler) DoMessage() error {

//	doMsgFunc := this.doMsgCallback

//	if doMsgFunc == nil {
//		return errors.New("doMsgFunc is nil")
//	}

//	err := doMsgFunc(this, this.conn)
//	return err
//}
func (this *MessageHandler) SetSessionData(sessionData interface{}) {
	this.SessionData = sessionData
}
func (this *MessageHandler) RunTcp(heartBeat time.Duration) error {
	tools.CaoSiShowDebug("message handler RunTcp: enter", heartBeat)

	// 设置网络连接参数
	var now, timeout time.Time
	var lastVisit time.Time
	var deltaTime time.Duration
	var err error
	var length int

	var OnReceive OnReceiveFunc
	var DoMessage DoMessageFunc
	var OnWrite OnWriteFunc

	conn := this.Conn
	defer conn.Close()

	OnReceive = this.callbacks.TcpOnReceiveCallback
	DoMessage = this.callbacks.TcpDoMessageCallback
	OnWrite = this.callbacks.TcpOnWriteCallback

	// 网络读取缓冲创建
	buffSize := 16 * 1024
	buffRead := make([]byte, buffSize, buffSize)
	forceTimeOut := time.Second * 10

	now = time.Now()
	lastVisit = now
	// 从网络中读取数据
	for {
		// 设置超时
		now = time.Now()
		deltaTime = time.Since(lastVisit)
		if deltaTime >= forceTimeOut {
			tools.ShowDebug("connection was closed by force timeout", conn.RemoteAddr())
			conn.Close()
			break
		}

		timeout = now.Add(time.Nanosecond)
		if false {
			conn.SetReadDeadline(timeout)
			conn.SetWriteDeadline(timeout)
		}

		// 尝试从网络中读取数据
		length, err = conn.Read(buffRead)
		// tools.CaoSiShowDebug("message handler RunTcp: read", length, err)
		if err != nil {

			if CheckTimeout(err) {
				length = 0
			} else {
				// 网络错误, 强行断开连接
				conn.Close()
				tools.ShowDebug("connection was closed", err.Error())
				break
			}
		}

		if length > 0 {
			// tools.CaoSiShowDebug("tcp read content", string(buffRead[:length]))
			lastVisit = now
			// 读到数据的话,灌注到handler中
			// buffRead[length] = 0 // 尾数置零
			err = OnReceive(this, buffRead[:length], 0)
			if err != nil {
				tools.ShowError("onReceiveCallback retrun err: ", err.Error())
				break
			}
		} else if length < 0 {
			// 读取长度异常
			conn.Close()
			tools.ShowDebug("connection was closed")
			break
		}

		// 消息处理
		err = DoMessage(this)
		if err != nil {
			tools.ShowError("doMsgCallback retrun err: ", err.Error())
			break
		}

		// 回写数据
		err = OnWrite(this)
		if err != nil {
			tools.ShowError("onWriteCallback retrun err: ", err.Error())
			break
		}

		if this.needClose {
			conn.Close()
		} else {
			// time.Sleep(heartBeat)
		}
	}

	conn.Close()
	return nil
}

func (this *MessageHandler) Close() {
	this.needClose = true
}

// 工厂函数
func CreateMessageHandler(conn net.Conn) (*MessageHandler, error) {
	ret := new(MessageHandler)
	err := ret.Init(conn)

	if err != nil {
		return nil, err
	}
	return ret, nil
}
