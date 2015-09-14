package socket

import (
	"bytes"
	"encoding/binary"
	"ultralisk/3rdparty/goArrayList/goArrayList"
	"ultralisk/tools"
)

const (
	MSG_HEAD_LEN          = 4
	MSG_MIN_LEN           = 1
	CLIENT_MAX_MSG_LENGTH = 4 * 1024 // 客户端最大消息长度限制

)

var BinOrder = binary.LittleEndian

// 用于socket连接的iobuffer
type IOBuffer struct {
	// 限制参数
	inputMsgLenLimit  uint32 // 输入消息长度的最大限制(0表示无限制)
	outputMsgLenLimit uint32 // 输出消息长度的最大限制(0表示无限制)
	// 输入相关
	inputBuffer    *bytes.Buffer // 输入缓存数据
	curInputMsgLen uint32        // 当前正在处理的输入消息长度 <=0 表示没有消息正在被处理
	inputMsgs      *goArrayList.ArrayList

	// 输出相关
	outputBuffer *bytes.Buffer // 输出数据

}

////////////////////////////////////////////////////////////////////////////////
// 对接口 ulio.IOBuffer的实现
////////////////////////////////////////////////////////////////////////////////
// 初始化
func (this *IOBuffer) Init(inputMsgLenLimit, outputMsgLenLimit uint32) error {
	this.inputMsgLenLimit = inputMsgLenLimit
	this.outputMsgLenLimit = outputMsgLenLimit

	// 临时输入缓存数据初始化
	this.inputBuffer = bytes.NewBuffer([]byte{})
	this.inputMsgs = new(goArrayList.ArrayList)

	// 输出缓存数据初始化
	this.outputBuffer = bytes.NewBuffer([]byte{})
	return nil
}

// 压入输入数据
// __TODO:此函数算法有待优化
func (this *IOBuffer) PushInputData(data []byte) error {
	if data == nil {
		return nil
	}

	// 将新到的数据写入到缓冲中
	inputBuffer := this.inputBuffer
	inputBuffer.Write(data)
	inputBufferDataLen := inputBuffer.Len()

	var curInputMsgLen, curInputMsgTotalLen uint32 // 当前正在处理的消息长度(不包含消息头)
	var buf []byte                                 // 临时缓存
	var err error
	var n int
	//var inputBufBytes []byte

	curInputMsgLen = this.curInputMsgLen     // 不含包头的消息长度
	for inputBufferDataLen >= MSG_HEAD_LEN { // 只要buffer中的剩余数据长度达到一个包头的长度就需要继续解析下去
		if curInputMsgLen <= 0 { // 当前没有消息在等待数据
			// 解出包头(也就是消息长度)
			err = binary.Read(inputBuffer, BinOrder, &curInputMsgLen)
			//inputBufBytes = inputBuffer.Bytes()
			//curInputMsgLen = BinOrder.Uint32(inputBufBytes[:MSG_HEAD_LEN])
			// TODO: 非法消息长度的细致化处理, 如何通知外部?
			if err != nil {
				return err
			} else if curInputMsgLen <= MSG_MIN_LEN {
				return tools.Error("invalid msg len = %d", curInputMsgLen)
			} else {
				// 客户端消息超长检查
				if this.inputMsgLenLimit > 0 && curInputMsgLen > this.inputMsgLenLimit {
					return tools.Error("input msg length to large! len = %d", curInputMsgLen)
				}
			}
		}

		curInputMsgTotalLen = curInputMsgLen + MSG_HEAD_LEN
		if inputBufferDataLen < int(curInputMsgTotalLen) {
			break // 数据尚未完全达到,退出处理循环
		}

		// 已经有一个包的数据完全到达,准备将该包提取出来
		buf = make([]byte, curInputMsgLen)
		n, err = inputBuffer.Read(buf)
		if err != nil {
			return err
		} else if n != int(curInputMsgLen) {
			return tools.Error("read msg err n = %d , len = ", n, curInputMsgLen)
		}
		this.inputMsgs.Append(buf) // 放入完整包缓存中
		inputBufferDataLen -= n + MSG_HEAD_LEN
		curInputMsgLen = 0

	}

	this.curInputMsgLen = curInputMsgLen
	return nil
}

// 弹出输入数据
func (this *IOBuffer) PopInputData() ([]byte, error) {

	inputMsgs := this.inputMsgs
	if inputMsgs.Size() > 0 {
		obj := inputMsgs.Get(0)
		inputMsgs.Remove(0)
		if obj != nil {
			ret, ok := obj.([]byte)
			if ok {
				// tools.CaoSiShowDebug("pop input data ret", string(ret))
				return ret, nil
			} else {
				return nil, tools.Error("pop input data faild, can not format as []byte , v = %v", obj)
			}
		}
	}
	return nil, nil
}

// 压入输出数据
func (this *IOBuffer) PushOutputData(data []byte) error {
	if data == nil {
		return nil
	}
	var err error
	var length uint32
	var buf [8]byte

	length = uint32(len(data))
	BinOrder.PutUint32(buf[:], length)

	_, err = this.outputBuffer.Write(buf[0:4])
	if err != nil {
		return err
	}
	_, err = this.outputBuffer.Write(data)

	return err
}

// 弹出输出数据
func (this *IOBuffer) PopOutputData() ([]byte, error) {
	if this.outputBuffer.Len() > 0 {
		ret := this.outputBuffer.Bytes()
		this.outputBuffer.Reset()
		return ret, nil
	}
	return nil, nil
}

////////////////////////////////////////////////////////////////////////////////
// 其它代码
////////////////////////////////////////////////////////////////////////////////
