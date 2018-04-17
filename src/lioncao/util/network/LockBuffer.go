package network

import (
	"lioncao/util/tools"
)

/*const (
	MSG_HEAD_LEN          = 4
	MSG_MIN_LEN           = 1
	CLIENT_MAX_MSG_LENGTH = 4 * 1024 // 客户端最大消息长度限制

)

var BinOrder = binary.LittleEndian*/

// 用于socket连接的iobuffer
type LockBuffer struct {
	IOBuffer
	inputLock tools.FastLock
	outLock   tools.FastLock
}

////////////////////////////////////////////////////////////////////////////////
// 对接口 ulio.IOBuffer的实现
////////////////////////////////////////////////////////////////////////////////
// 压入输入数据
// __TODO:此函数算法有待优化
func (this *LockBuffer) PushInputDataLock(data []byte) error {
	this.inputLock.Lock()
	defer this.inputLock.Unlock()
	return this.PushInputData(data)
}

// 弹出输入数据
func (this *LockBuffer) PopInputDataLock() ([]byte, error) {
	this.inputLock.Lock()
	defer this.inputLock.Unlock()
	return this.PopInputData()
}

// 压入输出数据
func (this *LockBuffer) PushOutputDataLock(data []byte) error {
	this.outLock.Lock()
	defer this.outLock.Unlock()
	return this.PushOutputData(data)
}

// 弹出输出数据
func (this *LockBuffer) PopOutputDataLock() ([]byte, error) {
	this.outLock.Lock()
	defer this.outLock.Unlock()
	return this.PopOutputData()
}

////////////////////////////////////////////////////////////////////////////////
// 其它代码
////////////////////////////////////////////////////////////////////////////////
