package main

import (
    "fmt"
    "gnmi_server/pkg/gnmi"
    "gnmi_server/pkg/gnmi/handler/get"
    "gnmi_server/pkg/gnmi/handler/set"
    "net"

    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "google.golang.org/grpc"

    log "github.com/golang/glog"
)

func main() {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 5010))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer(grpc.RPCDecompressor(grpc.NewGZIPDecompressor()))

    server := gnmi.DefaultServer(get.GetServeMux(), set.SetServeMux())
    gpb.RegisterGNMIServer(grpcServer, &server)
    grpcServer.Serve(listener)
}
