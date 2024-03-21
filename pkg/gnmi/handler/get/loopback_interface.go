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

func LoopbackInterfaceHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定Loopback Interface或全部Loopback Interface
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["loopback-interface-name"]; ok {
        spec = v
    }

    sli := &sonicpb.AcctonLoopbackInterface{
        LoopbackInterface: &sonicpb.AcctonLoopbackInterface_LoopbackInterface{},
    }
    if hkeys, err := conn.GetKeys("LOOPBACK_INTERFACE", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            if len(keys) != 1 {
                continue
            }
            c := cmd.NewIfAdapter(cmd.LOOPBACK_INTERFACE, keys[0], db)
            if data, err := c.Show(r.Type); err != nil {
                return nil, err
            } else {
                sli.LoopbackInterface.LoopbackInterfaceList = append(sli.LoopbackInterface.LoopbackInterfaceList,
                    &sonicpb.AcctonLoopbackInterface_LoopbackInterface_LoopbackInterfaceListKey{
                        LoopbackInterfaceName: keys[0],
                        LoopbackInterfaceList: data.(*sonicpb.AcctonLoopbackInterface_LoopbackInterface_LoopbackInterfaceList),
                    })
            }
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sli)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func LoopbackInterfaceIPPrefixHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定Interface或全部Interface
    kvs := handler.FetchPathKey(r)
    spec := []string{}
    if v, ok := kvs["loopback-interface-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["ip-prefix"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    sli := &sonicpb.AcctonLoopbackInterface{
        LoopbackInterface: &sonicpb.AcctonLoopbackInterface_LoopbackInterface{},
    }
    if hkeys, err := conn.GetKeys("LOOPBACK_INTERFACE", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := cmd.NewIfAddrAdapter(cmd.LOOPBACK_INTERFACE, keys[0], keys[1], db)
            if data, err := c.Show(r.Type); err != nil {
                return nil, err
            } else {
                sli.LoopbackInterface.LoopbackInterfaceIpprefixList = append(sli.LoopbackInterface.LoopbackInterfaceIpprefixList,
                    &sonicpb.AcctonLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixListKey{
                        LoopbackInterfaceName: keys[0],
                        IpPrefix: keys[1],
                        LoopbackInterfaceIpprefixList: data.(*sonicpb.AcctonLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixList),
                    })
            }
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sli)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}