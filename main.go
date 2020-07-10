package main

import (
    "fmt"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/server"
    grpc "google.golang.org/grpc"
    "net"

    log "github.com/golang/glog"
)

func main() {
    fmt.Println("123")
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 5001))
    if err != nil {
        log.Fatal("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer(grpc.RPCDecompressor(grpc.NewGZIPDecompressor()))

    server := server.DefaultServer(newGetServeMux(), newSetServeMux())
    gpb.RegisterGNMIServer(grpcServer, &server)
    grpcServer.Serve(listener)
}
