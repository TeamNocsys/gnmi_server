package main

import (
	"fmt"
	"gnmi_server/pkg/gnmi"
	get_handler "gnmi_server/pkg/gnmi/handler/get"
	"net"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	grpc "google.golang.org/grpc"

	log "github.com/golang/glog"
)

func newGetServeMux() *gnmi.GetServeMux {
	mux := gnmi.NewGetServeMux()
	mux.AddRouter("/test", get_handler.Test).
		AddRouter("/sonic-platform/platform", get_handler.ComponentInfoHandler).
		AddRouter("/sonic-platform/platform/component-list/fan", get_handler.FanInfoHandler).
		AddRouter("/sonic-platform/platform/component-list/power-supply", get_handler.PowerSupplyInfoHandler).
		AddRouter("/sonic-platform/platform/component-list/temperature", get_handler.TemperatureInfoHandler)
	return mux
}

func newSetServeMux() *gnmi.SetServeMux {
	return gnmi.NewSetServeMux()
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 5010))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.RPCDecompressor(grpc.NewGZIPDecompressor()))

	server := gnmi.DefaultServer(newGetServeMux(), newSetServeMux())
	gpb.RegisterGNMIServer(grpcServer, &server)
	grpcServer.Serve(listener)
}
