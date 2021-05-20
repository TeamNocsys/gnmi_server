package gnmi

import (
    "context"
    "encoding/json"
    "gnmi_server/cmd/command"
    "sync"
    "time"

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

func (entry *getMuxEntry) handle(ctx context.Context, r *gpb.GetRequest, db command.Client) (ret *gpb.GetResponse, err error) {
    logrus.Debugf("%s|GET", entry.pattern)
    begin := time.Now()
    defer func() {
        if err != nil {
            logrus.Errorf("%s|GET|%s|%s", entry.pattern, time.Now().Sub(begin), err.Error())
        } else {
            logrus.Debugf("%s|GET|%s", entry.pattern, time.Now().Sub(begin))
        }
    }()
    ret, err = entry.h(ctx, r, db)
    return
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

func (gsm *GetServeMux) DoHandle(ctx context.Context, req *gpb.GetRequest, db command.Client) (r *gpb.GetResponse, err error) {
    paths := req.GetPath()
    if len(paths) == 0 {
        return nil, status.Errorf(codes.InvalidArgument, "get request path is empty")
    } else if len(paths) > 1 {
        return nil, status.Errorf(codes.Unimplemented, "unsupported more than one path in single request")
    }

    xPath, _ := ParseXPath(req.GetPrefix(), paths[0])
    h, ok := gsm.m[xPath]
    if !ok {
        logrus.Errorf("%s|GET|unhandled", xPath)
        return nil, status.Errorf(codes.NotFound, "invalid path")
    }
    return h.handle(ctx, req, db)
}

type DeleteHandler func(context.Context, map[string]string, command.Client) error
type UpdateHandler func(context.Context, *gpb.TypedValue, command.Client) error

type setMuxEntry struct {
    dh DeleteHandler
    uh UpdateHandler
    pattern string
}

func (entry *setMuxEntry) handle(ctx context.Context, r interface{}, db command.Client) (err error) {
    begin := time.Now()
    oper := "Unknown"
    defer func() {
        if err != nil {
            logrus.Errorf("%s|%s|%s", entry.pattern, oper, time.Now().Sub(begin))
        } else {
            logrus.Debugf("%s|%s|%s", entry.pattern, oper, time.Now().Sub(begin))
        }
    }()
    switch r.(type) {
    case map[string]string:
        logrus.Debugf("%s|DEL", entry.pattern)
        oper = "DEL"
        s, _ := json.Marshal(r.(map[string]string))
        logrus.Tracef("%s|DEL|%s", entry.pattern, s)
        return entry.dh(ctx, r.(map[string]string), db)
    case *gpb.TypedValue:
        logrus.Debugf("%s|SET", entry.pattern)
        oper = "SET"
        return entry.uh(ctx, r.(*gpb.TypedValue), db)
    default:
        return status.Errorf(codes.InvalidArgument, "Unknown request parameter type")
    }
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

func (ssm *SetServeMux) AddDeleteRouter(pattern string, h DeleteHandler) *SetServeMux {
    ssm.mu.Lock()
    defer ssm.mu.Unlock()
    ssm.dm[pattern] = setMuxEntry{dh: h, pattern: pattern}
    return ssm
}

func (ssm *SetServeMux) AddReplaceRouter(pattern string, h UpdateHandler) *SetServeMux {
    ssm.mu.Lock()
    defer ssm.mu.Unlock()
    ssm.rm[pattern] = setMuxEntry{uh: h, pattern: pattern}
    return ssm
}

func (ssm *SetServeMux) AddUpdateRouter(pattern string, h UpdateHandler) *SetServeMux {
    ssm.mu.Lock()
    defer ssm.mu.Unlock()
    ssm.um[pattern] = setMuxEntry{uh: h, pattern: pattern}
    return ssm
}

func (ssm *SetServeMux) DoHandle(ctx context.Context, req *gpb.SetRequest, db command.Client) (r *gpb.SetResponse, err error) {
    if len(req.Delete) + len(req.Update) + len(req.Replace) != 1 {
        logrus.Error("Unsupported more than one path in single set request")
        return nil, status.Errorf(codes.Unimplemented, "unsupported more than one path in single request")
    }
    logrus.Debug("TRANSACTION")
    begin := time.Now()
    // TODO: 开启数据库事物
    defer func() {
        if err != nil {
            // TODO: 回滚数据库事物
            logrus.Errorf("ROLLBACK|%s|%s", time.Now().Sub(begin), err.Error())
        } else {
            // TODO: 提交数据库事物
            logrus.Debugf("COMMIT|%s", time.Now().Sub(begin))
        }
    }()
    r = &gpb.SetResponse{}
    // 删除处理
    for _, v := range req.Delete {
        xpath, kvs := ParseXPath(req.GetPrefix(), v)
        if h, ok := ssm.dm[xpath]; !ok {
            logrus.Errorf("%s|DEL|unhandled", xpath)
            return nil, status.Errorf(codes.NotFound, "invalid path")
        } else {
            if err = h.handle(ctx, kvs, db); err != nil {
                return
            } else {
                r.Response = append(r.Response, &gpb.UpdateResult{
                    Path: v,
                    Op: gpb.UpdateResult_DELETE,
                })
            }
        }
    }
    // 更新处理
    for _, v := range req.Update {
        xpath, _ := ParseXPath(req.GetPrefix(), v.Path)
        if h, ok := ssm.um[xpath]; !ok {
            logrus.Errorf("%s|SET|unhandled", xpath)
            return nil, status.Errorf(codes.NotFound, "invalid path")
        } else {
            if err = h.handle(ctx, v.Val, db); err != nil {
                return
            } else {
                r.Response = append(r.Response, &gpb.UpdateResult{
                    Path: v.Path,
                    Op: gpb.UpdateResult_UPDATE,
                })
            }
        }
    }
    // 替换处理
    for _, v := range req.Replace {
        xpath, _ := ParseXPath(req.GetPrefix(), v.Path)
        if h, ok := ssm.rm[xpath]; !ok {
            logrus.Errorf("%s|SET|unhandled", xpath)
            return nil, status.Errorf(codes.NotFound, "invalid path")
        } else {
            if err = h.handle(ctx, v.Val, db); err != nil {
                return
            } else {
                r.Response = append(r.Response, &gpb.UpdateResult{
                    Path: v.Path,
                    Op: gpb.UpdateResult_REPLACE,
                })
            }
        }
    }

    r.Prefix = req.Prefix
    r.Timestamp = time.Now().Unix()
    return
}

func ParseXPath(prefix *gpb.Path, path *gpb.Path) (string, map[string]string) {
    var xpath string
    kvs := map[string]string{}
    if prefix != nil {
        for _, elem := range prefix.GetElem() {
            xpath += "/" + elem.GetName()
            if elem.Key != nil {
                for k,v := range elem.Key {
                    kvs[k] = v
                }
            }
        }
    }

    for _, elem := range path.GetElem() {
        xpath += "/" + elem.GetName()
        if elem.Key != nil {
            for k,v := range elem.Key {
                kvs[k] = v
            }
        }
    }
    return xpath, kvs
}