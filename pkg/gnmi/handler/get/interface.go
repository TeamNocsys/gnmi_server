package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
    "gnmi_server/pkg/gnmi/handler"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func InterfaceHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定Interface或全部Interface
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["port-name"]; ok {
        spec = v
    }

    si := &sonicpb.AcctonInterface{
        Interface: &sonicpb.AcctonInterface_Interface{},
    }
    if hkeys, err := conn.GetKeys("INTERFACE", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            if len(keys) != 1 {
                continue
            }
            c := cmd.NewIfAdapter(cmd.INTERFACE, keys[0], db)
            if data, err := c.Show(r.Type); err != nil {
                return nil, err
            } else {
                si.Interface.InterfaceList = append(si.Interface.InterfaceList,
                    &sonicpb.AcctonInterface_Interface_InterfaceListKey{
                        PortName: keys[0],
                        InterfaceList: data.(*sonicpb.AcctonInterface_Interface_InterfaceList),
                    })
            }
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, si)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func InterfaceIPPrefixHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定Interface或全部Interface
    kvs := handler.FetchPathKey(r)
    spec := []string{}
    if v, ok := kvs["port-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["ip-prefix"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    si := &sonicpb.AcctonInterface{
        Interface: &sonicpb.AcctonInterface_Interface{},
    }
    if hkeys, err := conn.GetKeys("INTERFACE", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := cmd.NewIfAddrAdapter(cmd.INTERFACE, keys[0], keys[1], db)
            if data, err := c.Show(r.Type); err != nil {
                return nil, err
            } else {
                si.Interface.InterfaceIpprefixList = append(si.Interface.InterfaceIpprefixList,
                    &sonicpb.AcctonInterface_Interface_InterfaceIpprefixListKey{
                        PortName: keys[0],
                        IpPrefix: keys[1],
                        InterfaceIpprefixList: data.(*sonicpb.AcctonInterface_Interface_InterfaceIpprefixList),
                    })
            }
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, si)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}