package gnmi

import (
    "context"
    "crypto/tls"
    "errors"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "google.golang.org/grpc"
    codes "google.golang.org/grpc/codes"
    "google.golang.org/grpc/credentials"
    status "google.golang.org/grpc/status"
    "math"
    "time"
)

type Auth struct {
    username string
    password string
}

type Server struct {
    Auth
    dbClient command.Client
    getServeMux *GetServeMux
    setServeMux *SetServeMux
}

func DefaultServer(username, password string, dbClient command.Client, gmux *GetServeMux, smux *SetServeMux) Server {
    return Server{
        Auth{
            username,
            password,
        },
        dbClient,
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
    return s.getServeMux.DoHandle(ctx, request, s.dbClient)
}

func (s *Server) Set(ctx context.Context, request *gpb.SetRequest) (*gpb.SetResponse, error) {
    if s.setServeMux == nil {
        return nil, status.Errorf(codes.Unimplemented, "serve multiplexer of set request is null")
    }
    return s.setServeMux.DoHandle(ctx, request, s.dbClient)
}

type gNMIAuth struct {
    username string
    password string
    secure bool
}

func (ga *gNMIAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
    return map[string]string{
        "username": ga.username,
        "password": ga.password,
    }, nil
}

func (ga *gNMIAuth) RequireTransportSecurity() bool {
    return ga.secure
}

func (s *Server) Subscribe(stream gpb.GNMI_SubscribeServer) (err error) {
    logrus.Debug("SUBSCRIBE")
    begin := time.Now()
    defer func() {
        if err != nil {
            logrus.Errorf("SUBSCRIBE|%s|%s", time.Now().Sub(begin), err.Error())
        } else {
            logrus.Debugf("SUBSCRIBE|%s", time.Now().Sub(begin))
        }
    }()
    if s.username == "" || s.password == "" {
        return errors.New("invalid user name and password")
    }

    conn, err := grpc.DialContext(stream.Context(), "127.0.0.1:8080",
        grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
        grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
        grpc.WithPerRPCCredentials(&gNMIAuth{
            username: s.username,
            password: s.password,
            secure: true,
        }),
    )
    if err != nil {
        return err
    }
    defer conn.Close()
    cli := gpb.NewGNMIClient(conn)
    proxy, err := cli.Subscribe(stream.Context())
    if err != nil {
        return err
    }
    waitc := make(chan struct{})
    go func() {
        defer proxy.CloseSend()
        var req *gpb.SubscribeRequest
        for {
            req, err = stream.Recv()
            if err != nil {
                break
            }
            subscribes := req.GetSubscribe()
            if subscribes != nil {
                subscribes.GetPrefix()
                for _, subscribe := range subscribes.GetSubscription() {
                    xpath, _ := ParseXPath(subscribes.GetPrefix(), subscribe.Path)
                    logrus.Debug(xpath + "|SUBSCRIBE")
                }
            }
            if err = proxy.Send(req); err != nil {
                break
            }
        }
        close(waitc)
    }()

    var rep *gpb.SubscribeResponse
    for {
        rep, err = proxy.Recv()
        if err != nil {
            return
        }
        if err = stream.Send(rep); err != nil {
            return
        }
    }
    <-waitc
    return nil
}
