package main

import (
	"gnmi_server/server"
	get_handler "gnmi_server/handler/get"
)

func newGetServeMux() *server.GetServeMux {
	mux := server.NewGetServeMux()
    mux.AddRouter("/test", get_handler.Test).
		AddRouter("/components/component/fan", get_handler.Get_fan_info)
	return mux
}

func newSetServeMux() *server.SetServeMux {
	return server.NewSetServeMux()
}
