package tools

import (
	"os"
	"os/signal"
)

// 系统信号监听者
type SigListener interface {
	OnSig(sig os.Signal)
}

// 启动监听
func StartSigListen(sl SigListener) chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	go sigHandler(ch, sl)
	return ch
}

// 信号处理
func sigHandler(ch chan os.Signal, sl SigListener) {
	for {
		sig := <-ch
		sl.OnSig(sig)
	}
}
