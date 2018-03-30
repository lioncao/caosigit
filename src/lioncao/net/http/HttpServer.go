package http

import (
	"lioncao/util/tools"
	"net/http"
)

type HttpHandler struct {
	server *http.Server
}

func (this *HttpHandler) init(server *http.Server) error {
	this.server = server

	return nil
}

func (this *HttpHandler) RunHttp() error {
	server := this.server
	err := server.ListenAndServe()
	if err != nil {
		return tools.Error("RunHttp err: addr=%s,mux=%v,err=%s", server.Addr, server.Handler, err.Error())
	}
	return nil
}

func CreateHttpHandler(server *http.Server) (*HttpHandler, error) {
	if server == nil {

		return nil, tools.Error("httpHandler.Init para err, server=%v", server)
	}

	var err error
	ret := new(HttpHandler)

	err = ret.init(server)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
