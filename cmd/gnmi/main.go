package main

import (
	"fmt"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"gnmi_server/pkg/gnmi"
	get_handler "gnmi_server/pkg/gnmi/handler/get"
	grpc "google.golang.org/grpc"
	"net"

	log "github.com/golang/glog"
)

func newGetServeMux() *gnmi.GetServeMux {
	mux := gnmi.NewGetServeMux()
	mux.AddRouter("/test", get_handler.Test).
		AddRouter("/components/component/fan", get_handler.Get_fan_info)
	return mux
}

func newSetServeMux() *gnmi.SetServeMux {
	return gnmi.NewSetServeMux()
}

func main() {
    fmt.Println("123")
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 5001))
    if err != nil {
        log.Fatal("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer(grpc.RPCDecompressor(grpc.NewGZIPDecompressor()))

    server := gnmi.DefaultServer(newGetServeMux(), newSetServeMux())
    gpb.RegisterGNMIServer(grpcServer, &server)
    grpcServer.Serve(listener)
}
