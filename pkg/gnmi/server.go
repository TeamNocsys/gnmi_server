package gnmi

import (
    "context"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    codes "google.golang.org/grpc/codes"
    status "google.golang.org/grpc/status"
)

type Server struct {
    getServeMux *GetServeMux
    setServeMux *SetServeMux
}

func DefaultServer(gmux *GetServeMux, smux *SetServeMux) Server {
    return Server{
        gmux,
        smux,
    }
}

func (s *Server) AddServeMux(gmux *GetServeMux, smux *SetServeMux) {
    s.getServeMux = gmux
    s.setServeMux = smux
}

func (s *Server) Capabilities(ctx context.Context, request *gpb.CapabilityRequest) (*gpb.CapabilityResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "method Capabilities not implemented")
}

func (s *Server) Get(ctx context.Context, request *gpb.GetRequest) (*gpb.GetResponse, error) {
    if s.getServeMux == nil {
        return nil, status.Errorf(codes.Unimplemented, "serve multiplexer of get request is null")
    }
    return s.getServeMux.DoHandle(ctx, request)
}

func (s *Server) Set(ctx context.Context, request *gpb.SetRequest) (*gpb.SetResponse, error) {
    if s.setServeMux == nil {
        return nil, status.Errorf(codes.Unimplemented, "serve multiplexer of set request is null")
    }
    return s.setServeMux.DoHandle(ctx, request)
}

func (s *Server) Subscribe(server gpb.GNMI_SubscribeServer) error {
    return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}