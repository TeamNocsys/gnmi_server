package gnmi

import (
    "context"
    "gnmi_server/cmd/command"
    "sync"

    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type GetHandler func(context.Context, *gpb.GetRequest, command.Client) (*gpb.GetResponse, error)

type getMuxEntry struct {
    h       GetHandler
    pattern string
}

func (entry *getMuxEntry) handle(ctx context.Context, r *gpb.GetRequest, db command.Client) (*gpb.GetResponse, error) {
    return entry.h(ctx, r, db)
}

type GetServeMux struct {
    mu sync.RWMutex
    m  map[string]getMuxEntry
}

func NewGetServeMux() *GetServeMux {
    serveMux := GetServeMux{}
    serveMux.m = make(map[string]getMuxEntry)
    return &serveMux
}

func (gsm *GetServeMux) AddRouter(pattern string, h GetHandler) *GetServeMux {
    gsm.mu.Lock()
    defer gsm.mu.Unlock()
    gsm.m[pattern] = getMuxEntry{h, pattern}
    return gsm
}

func (gsm *GetServeMux) DoHandle(ctx context.Context, req *gpb.GetRequest, db command.Client) (*gpb.GetResponse, error) {
    paths := req.GetPath()
    if len(paths) == 0 {
        return nil, status.Errorf(codes.InvalidArgument, "get request path is empty")
    } else if len(paths) > 1 {
        return nil, status.Errorf(codes.Unimplemented, "unsupported more than one path in single request")
    }

    path := generalPath(req.GetPrefix(), paths[0])
    h, ok := gsm.m[path]
    if !ok {
        logrus.Error("Unhandled Get XPath: ", path)
        return nil, status.Errorf(codes.NotFound, "invalid path")
    }

    return h.handle(ctx, req, db)
}

type SetHandler func(context.Context, *gpb.SetRequest, command.Client) (*gpb.SetResponse, error)

type setMuxEntry struct {
    h       SetHandler
    pattern string
}

func (entry *setMuxEntry) handle(ctx context.Context, r *gpb.SetRequest, db command.Client) (*gpb.SetResponse, error) {
    return entry.h(ctx, r, db)
}

type SetServeMux struct {
    mu sync.RWMutex
    dm map[string]setMuxEntry // Delete request
    rm map[string]setMuxEntry // Replace request
    um map[string]setMuxEntry // Update request
}

func NewSetServeMux() *SetServeMux {
    serveMux := SetServeMux{}
    serveMux.dm = make(map[string]setMuxEntry)
    serveMux.rm = make(map[string]setMuxEntry)
    serveMux.um = make(map[string]setMuxEntry)
    return &serveMux
}

func (ssm *SetServeMux) AddDeleteRouter(pattern string, h SetHandler) *SetServeMux {
    ssm.mu.Lock()
    defer ssm.mu.Unlock()
    ssm.dm[pattern] = setMuxEntry{h, pattern}
    return ssm
}

func (ssm *SetServeMux) AddReplaceRouter(pattern string, h SetHandler) *SetServeMux {
    ssm.mu.Lock()
    defer ssm.mu.Unlock()
    ssm.rm[pattern] = setMuxEntry{h, pattern}
    return ssm
}

func (ssm *SetServeMux) AddUpdateRouter(pattern string, h SetHandler) *SetServeMux {
    ssm.mu.Lock()
    defer ssm.mu.Unlock()
    ssm.um[pattern] = setMuxEntry{h, pattern}
    return ssm
}

func (ssm *SetServeMux) DoHandle(ctx context.Context, req *gpb.SetRequest, db command.Client) (*gpb.SetResponse, error) {
    if (len(req.Delete) + len(req.Replace) + len(req.Update)) > 1 {
        return nil, status.Errorf(codes.Unimplemented, "unsupported more than one path in single request")
    }

    var h setMuxEntry
    var ok bool
    var path string
    if len(req.Delete) == 1 {
        path = generalPath(req.GetPrefix(), req.Delete[0])
        h, ok = ssm.dm[path]
    } else if len(req.Replace) == 1 {
        path = generalPath(req.GetPrefix(), req.Replace[0].GetPath())
        h, ok = ssm.rm[path]
    } else if len(req.Update) == 1 {
        path = generalPath(req.GetPrefix(), req.Update[0].GetPath())
        h, ok = ssm.um[path]
    } else {
        return nil, status.Errorf(codes.InvalidArgument, "set request path is empty")
    }

    if !ok {
        logrus.Error("Unhandled Set XPath: ", path)
        return nil, status.Errorf(codes.NotFound, "invalid path")
    }

    return h.handle(ctx, req, db)
}

func generalPath(prefix *gpb.Path, path *gpb.Path) string {
    var gPath string
    for _, pathElem := range prefix.GetElem() {
        gPath += "/" + pathElem.GetName()
    }

    for _, pathElem := range path.GetElem() {
        gPath += "/" + pathElem.GetName()
    }
    return gPath
}
