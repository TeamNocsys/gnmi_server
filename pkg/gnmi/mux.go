package gnmi

import (
    "context"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "sync"
)

type GetHandler func(context.Context, *gpb.GetRequest) (*gpb.GetResponse, error)

type getMuxEntry struct {
    h       GetHandler
    pattern string
}

func (entry *getMuxEntry) handle(ctx context.Context, r *gpb.GetRequest) (*gpb.GetResponse, error) {
    return entry.h(ctx, r)
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

func (gsm *GetServeMux) DoHandle(ctx context.Context, r *gpb.GetRequest) (*gpb.GetResponse, error) {
    paths := r.GetPath()
    if len(paths) == 0 {
        return nil, status.Errorf(codes.InvalidArgument, "get request path is empty")
    } else if len(paths) > 1 {
        return nil, status.Errorf(codes.Unimplemented, "unsupported more than one path in single request")
    }

    path := generalPath(paths[0])
    h, ok := gsm.m[path]
    if !ok {
        return nil, status.Errorf(codes.NotFound, "invalid path")
    }

    return h.handle(ctx, r)
}

type SetHandler func(context.Context, *gpb.SetRequest) (*gpb.SetResponse, error)

type setMuxEntry struct {
    h       SetHandler
    pattern string
}

func (entry *setMuxEntry) handle(ctx context.Context, r *gpb.SetRequest) (*gpb.SetResponse, error) {
    return entry.h(ctx, r)
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

func (ssm *SetServeMux) DoHandle(ctx context.Context, r *gpb.SetRequest) (*gpb.SetResponse, error) {
    if (len(r.Delete) + len(r.Replace) + len(r.Update)) > 1 {
        return nil, status.Errorf(codes.Unimplemented, "unsupported more than one path in single request")
    }

    var h setMuxEntry
    var ok bool
    if len(r.Delete) == 1 {
        path := generalPath(r.Delete[0])
        h, ok = ssm.dm[path]
    } else if len(r.Replace) == 1 {
        path := generalPath(r.Replace[0].Path)
        h, ok = ssm.rm[path]
    } else if len(r.Update) == 1 {
        path := generalPath(r.Update[0].Path)
        h, ok = ssm.um[path]
    } else {
        return nil, status.Errorf(codes.InvalidArgument, "set request path is empty")
    }

    if !ok {
        return nil, status.Errorf(codes.NotFound, "invalid path")
    }

    return h.handle(ctx, r)
}

func generalPath(path *gpb.Path) string {
    var gPath string
    elements := path.GetElem()
    for _, pathElem := range elements {
        gPath += "/" + pathElem.GetName()
    }
    return gPath
}
