package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk/helper"
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

    sli := &sonicpb.SonicLoopbackInterface{
        LoopbackInterface: &sonicpb.SonicLoopbackInterface_LoopbackInterface{},
    }
    if hkeys, err := conn.GetKeys("LOOPBACK_INTERFACE", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            if len(keys) != 1 {
                continue
            }
            c := helper.LoopbackInterface{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sli.LoopbackInterface.LoopbackInterfaceList = append(sli.LoopbackInterface.LoopbackInterfaceList,
                &sonicpb.SonicLoopbackInterface_LoopbackInterface_LoopbackInterfaceListKey{
                    LoopbackInterfaceName: keys[0],
                    LoopbackInterfaceList: c.Data,
                })
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

    sli := &sonicpb.SonicLoopbackInterface{
        LoopbackInterface: &sonicpb.SonicLoopbackInterface_LoopbackInterface{},
    }
    if hkeys, err := conn.GetKeys("LOOPBACK_INTERFACE", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := helper.LoopbackInterfaceIPPrefix{
                Keys: keys,
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sli.LoopbackInterface.LoopbackInterfaceIpprefixList = append(sli.LoopbackInterface.LoopbackInterfaceIpprefixList,
                &sonicpb.SonicLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixListKey{
                    LoopbackInterfaceName: keys[0],
                    IpPrefix: keys[1],
                    LoopbackInterfaceIpprefixList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sli)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}