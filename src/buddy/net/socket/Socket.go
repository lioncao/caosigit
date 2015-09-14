package socket

import (
	"buddy/util/tools"
	"fmt"
	"net"
)

type DoTcpFunc func(conn net.Conn) error

func StartTcp(ip, port string, doTcpFunc DoTcpFunc) error {
	var err error

	if ip == "" || port == "" {
		tools.ShowError(fmt.Sprintf("startCmd addr err: ip=%s,port=%s", ip, port))
		return nil
	}

	addr := fmt.Sprintf("%s:%s", ip, port)

	l, err := net.Listen("tcp", addr) // 启动监听
	if err != nil {
		tools.ShowError("StartTcp err", err.Error())
		return nil
	}

	tools.ShowDebug("StartCmd tcp starting :", addr)

	// 连接监听
	var c net.Conn
	for {
		c, err = l.Accept()
		if err != nil {
			if c != nil {
				err = c.Close()
				if err != nil {
					tools.ShowError("tcp l.Accept err", err)
				}
			}
		} else {
			tools.ShowInfo("tcp connection ok:", c.LocalAddr(), c.RemoteAddr())
			// 有新连接建立
			go doTcpFunc(c)
		}
	}
	return nil
}
