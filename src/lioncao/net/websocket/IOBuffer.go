package websocket

import (
	"lioncao/util/tools"
)

// 用于websocket连接的IOBuffer
type IOBuffer struct {
	lockInput  tools.FastLock
	lockOutput tools.FastLock

	// 输入相关
	inputMsgs []string // 输入缓存数据

	// 输出相关
	outputMsgs []string // 输入缓存数据

}

////////////////////////////////////////////////////////////////////////////////
// 对接口 ulio.IOBuffer的实现
////////////////////////////////////////////////////////////////////////////////
// 初始化
func (this *IOBuffer) Init() error {
	this.inputMsgs = make([]string, 0)
	this.outputMsgs = make([]string, 0)
	return nil
}

// 压入输入数据
func (this *IOBuffer) PushInputMsg(msg string) error {
	this.lockInput.Lock()
	defer this.lockInput.Unlock()

	if msg == "" {
		return nil
	}
	this.inputMsgs = append(this.inputMsgs, msg)
	return nil
}

// 弹出输入数据
func (this *IOBuffer) PopInputMsgs() ([]string, error) {
	this.lockInput.Lock()
	defer this.lockInput.Unlock()

	cnt := len(this.inputMsgs)
	// tools.TODO("msg len:", cnt, this)
	if cnt > 0 {
		msgs := this.inputMsgs[:cnt]
		this.inputMsgs = make([]string, 0)
		// tools.TODO("msg len:", len(this.inputMsgs), this)
		return msgs, nil
	}
	return nil, nil
}

// 压入输出数据
func (this *IOBuffer) PushOutputMsg(msg string) error {
	this.lockOutput.Lock()
	defer this.lockOutput.Unlock()

	if msg == "" {
		return nil
	}
	this.outputMsgs = append(this.outputMsgs, msg)
	return nil
}

// 弹出输出数据
func (this *IOBuffer) PopOutputMsg() ([]string, error) {
	this.lockOutput.Lock()
	defer this.lockOutput.Unlock()

	cnt := len(this.outputMsgs)
	if cnt > 0 {
		msgs := this.outputMsgs[:cnt]
		this.outputMsgs = make([]string, 0)
		return msgs, nil
	}
	return nil, nil
}

////////////////////////////////////////////////////////////////////////////////
// 其它代码
////////////////////////////////////////////////////////////////////////////////
