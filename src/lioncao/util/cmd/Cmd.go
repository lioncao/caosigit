package cmd

/*
命令行处理相关
*/
import (
	//"flag"
	"net"
	//"os"
	"3rdparty/goArrayList/goArrayList"
	bdsocket "lioncao/net/socket"
	"lioncao/util/tools"
	"time"
)

const (
	AUTH_STATUS_DEFAULT = iota
	AUTH_STATUS_AUTHING
	AUTH_STATUS_AUTH_OK
)

// 命令行的网络session数据
type CmdSession struct {
	authStatus int                    // 验证状态
	intputBuf  []byte                 // 输入缓存
	inputCmd   *goArrayList.ArrayList // 待解决的命令行
	outputResp *goArrayList.ArrayList // 输出
}

// 初始化
func (this *CmdSession) Init(inputMsgLenLimit, outputMsgLenLimit uint32) error {
	//this.inputMsgLenLimit = inputMsgLenLimit
	//this.outputMsgLenLimit = outputMsgLenLimit

	//// 临时输入缓存数据初始化
	//this.inputBuffer = bytes.NewBuffer([]byte{})
	//this.inputMsgs = new(goArrayList.ArrayList)

	//// 输出缓存数据初始化
	//this.outputBuffer = bytes.NewBuffer([]byte{})
	return nil
}

// 压入输入数据
// __TODO:此函数算法有待优化
func (this *CmdSession) PushInputData(data []byte) error {
	//if data == nil {
	//	return nil
	//}

	//// 将新到的数据写入到缓冲中
	//inputBuffer := this.inputBuffer
	//inputBuffer.Write(data)
	//inputBufferDataLen := inputBuffer.Len()

	//var curInputMsgLen, curInputMsgTotalLen uint32 // 当前正在处理的消息长度(不包含消息头)
	//var buf []byte                                 // 临时缓存
	//var err error
	//var n int
	////var inputBufBytes []byte

	//curInputMsgLen = this.curInputMsgLen     // 不含包头的消息长度
	//for inputBufferDataLen >= MSG_HEAD_LEN { // 只要buffer中的剩余数据长度达到一个包头的长度就需要继续解析下去
	//	if curInputMsgLen <= 0 { // 当前没有消息在等待数据
	//		// 解出包头(也就是消息长度)
	//		err = binary.Read(inputBuffer, binary.LittleEndian, &curInputMsgLen)
	//		//inputBufBytes = inputBuffer.Bytes()
	//		//curInputMsgLen = BinOrder.Uint32(inputBufBytes[:MSG_HEAD_LEN])
	//		// TODO: 非法消息长度的细致化处理, 如何通知外部?
	//		if err != nil {
	//			return err
	//		} else if curInputMsgLen <= MSG_MIN_LEN {
	//			return tools.Error("invalid msg len = %d", curInputMsgLen)
	//		} else {
	//			// 客户端消息超长检查
	//			if this.inputMsgLenLimit > 0 && curInputMsgLen > this.inputMsgLenLimit {
	//				return tools.Error("input msg length to large! len = %d", curInputMsgLen)
	//			}
	//		}
	//	}

	//	tools.CaoSiShowDebug("xxxxxxxxxxx  ", inputBufferDataLen, curInputMsgLen)
	//	curInputMsgTotalLen = curInputMsgLen + MSG_HEAD_LEN
	//	if inputBufferDataLen < int(curInputMsgTotalLen) {
	//		break // 数据尚未完全达到,退出处理循环
	//	}

	//	// 已经有一个包的数据完全到达,准备将该包提取出来
	//	buf = make([]byte, curInputMsgLen)
	//	n, err = inputBuffer.Read(buf)
	//	if err != nil {
	//		return err
	//	} else if n != int(curInputMsgLen) {
	//		return tools.Error("read msg err n = %d , len = ", n, curInputMsgLen)
	//	}
	//	this.inputMsgs.Append(buf) // 放入完整包缓存中
	//	tools.CaoSiShowDebug("yyyyyyyyyyyyy  ", inputBufferDataLen, curInputMsgLen, n)
	//	inputBufferDataLen -= n + MSG_HEAD_LEN
	//	curInputMsgLen = 0
	//	tools.CaoSiShowDebug("yyyyyyyyyyyyy  ", inputBufferDataLen, curInputMsgLen, n)

	//}

	//this.curInputMsgLen = curInputMsgLen
	return nil
}

// 弹出输入数据
func (this *CmdSession) PopInputData() ([]byte, error) {

	//inputMsgs := this.inputMsgs
	//if inputMsgs.Size() > 0 {
	//	obj := inputMsgs.Get(0)
	//	inputMsgs.Remove(0)
	//	if obj != nil {
	//		ret, ok := obj.([]byte)
	//		if ok {
	//			tools.CaoSiShowDebug("pop input data ret", string(ret))
	//			return ret, nil
	//		} else {
	//			return nil, tools.Error("pop input data faild, can not format as []byte , v = %v", obj)
	//		}
	//	}
	//}
	return nil, nil
}

// 压入输出数据
func (this *CmdSession) PushOutputData(data []byte) error {
	//if data == nil {
	//	return nil
	//}
	//var err error
	//var length uint32
	//var buf [8]byte

	//length = uint32(len(data))
	//BinOrder.PutUint32(buf[:], length)

	//_, err = this.outputBuffer.Write(buf[0:4])
	//if err != nil {
	//	return err
	//}
	//_, err = this.outputBuffer.Write(data)
	//return err
	return nil
}

// 弹出输出数据
func (this *CmdSession) PopOutputData() ([]byte, error) {
	//if this.outputBuffer.Len() > 0 {
	//	ret := this.outputBuffer.Bytes()
	//	tools.CaoSiShowDebug("pop out data ret 1   ", string(ret))
	//	this.outputBuffer.Reset()
	//	tools.CaoSiShowDebug("pop out data ret 2   ", ret)
	//	return ret, nil
	//}
	return nil, nil
}

type CmdData struct {
}

// 处理命令行的回调格式
type DoCmd func(session *CmdSession, cmdData *CmdData) error

// 命令行工具对象
type CmdMachine struct {
	title     string // 命令提示行
	password  string // 验证密码, 可为空
	DoCmdFunc DoCmd  // 具体命令的处理回调

	TcpCb bdsocket.TcpCallbacks
}

// 创建一个命令行工具
func NewCmdMachine(password string, title string, DoCmdFunc DoCmd) *CmdMachine {
	machine := new(CmdMachine)
	machine.password = password
	machine.title = title
	machine.DoCmdFunc = DoCmdFunc
	return machine
}

func NewCmdSession(machine *CmdMachine) *CmdSession {
	if machine == nil {
		return nil
	}
	session := new(CmdSession)

	if machine.password == "" {
		session.authStatus = AUTH_STATUS_AUTH_OK
	} else {
		session.authStatus = AUTH_STATUS_DEFAULT
	}
	return session
}

func NewCmdData(data []byte) *CmdData {
	if data == nil {
		return nil
	}

	ret := new(CmdData)

	return ret
}

func (this *CmdMachine) DoTcp(conn net.Conn) error {
	var err error
	//var length int
	var handler *bdsocket.MessageHandler
	var session *CmdSession

	// 创建消息处理句柄并设置参数
	handler, err = bdsocket.CreateMessageHandler(conn)
	if err != nil {
		tools.ShowError("CreateMessageHandler err,", conn, this.TcpCb.TcpDoMessageCallback)
		return nil
	}

	session = NewCmdSession(this)
	if session == nil {
		return tools.Error("create cmd session faild")
	}

	handler.SetSessionData(session)
	// 设置回调
	handler.SetCallback(nil)

	err = handler.RunTcp(time.Millisecond * 50)
	if err != nil {
		tools.ShowError("do tcp err :", err.Error())
	}
	return nil
}

type CmdTcpCallbacks struct {
	CmdMachine *CmdMachine
}

// tcp连接 MessageHandler的三个回调, 也就是bdsocket.TcpCallbacks interface{}的实现
func (this *CmdTcpCallbacks) SetSessionData(data interface{}) {
	this.CmdMachine, _ = data.(*CmdMachine)
}
func (this *CmdTcpCallbacks) GetSessionData() interface{} {
	return this.CmdMachine
}

func (this *CmdTcpCallbacks) TcpOnReceiveCallback(handler *bdsocket.MessageHandler, buf []byte, bufSize int) error {
	if handler == nil {
		return tools.Error("CmdTcpCallbacks.tcpOnReceiveCallback() get handler nil")
	}
	session, ok := handler.SessionData.(*CmdSession)
	if !ok {
		handler.Close()
		return tools.Error("CmdTcpCallbacks.tcpOnReceiveCallback() get session failed")
	}
	return session.PushInputData(buf)
}
func (this *CmdTcpCallbacks) TcpDoMessageCallback(handler *bdsocket.MessageHandler) error {
	var data []byte
	var err error
	var params *CmdData
	var session *CmdSession
	var ok bool

	session, ok = handler.SessionData.(*CmdSession)
	if !ok {
		handler.Close()
		return tools.Error("CmdTcpCallbacks.tcpDoMessageCallback() get session failed")
	}

	switch session.authStatus {
	default:
		fallthrough
	case AUTH_STATUS_DEFAULT:
		session.PushOutputData([]byte("password:"))
		session.authStatus = AUTH_STATUS_AUTHING
	case AUTH_STATUS_AUTHING:
		data, err = session.PopInputData()
		if err != nil {
			handler.Close()
			return err
		}

		if data != nil && (string(data) == this.CmdMachine.password) {
			session.authStatus = AUTH_STATUS_AUTH_OK
		} else {
			// TODO: 直接断掉可能会更好?
			session.PushOutputData([]byte("auth faild!\npassword:"))
			session.inputCmd.Clear()
		}
	case AUTH_STATUS_AUTH_OK:
		for {
			data, err = session.PopInputData()
			if err != nil {
				handler.Close()
				return err
			}

			if data == nil {
				break
			}
			params = NewCmdData(data)
			this.CmdMachine.DoCmdFunc(session, params)
		}
	}
	return nil
}
func (this *CmdTcpCallbacks) TcpOnWriteCallback(handler *bdsocket.MessageHandler) error {
	if handler == nil {

	}
	session, ok := handler.SessionData.(*CmdSession)
	if !ok {
		handler.Close()
		return tools.Error("CmdMachine.tcpDoMessageCallback() get session failed")
	}

	buf, err := session.PopOutputData()
	if err != nil {
		return err
	}

	if buf != nil {
		handler.Conn.Write(buf)
	}

	return nil
}
