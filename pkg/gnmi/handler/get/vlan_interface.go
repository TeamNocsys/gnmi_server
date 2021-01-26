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

func VlanInterfaceHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定VLAN Interface或全部VLAN Interface
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["vlan-name"]; ok {
        spec = v
    }

    sv := &sonicpb.SonicVlan{
        VlanInterface: &sonicpb.SonicVlan_VlanInterface{},
    }
    if hkeys, err := conn.GetKeys("VLAN_INTERFACE", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            if len(keys) != 1 {
                continue
            }
            c := helper.VlanInterface{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sv.VlanInterface.VlanInterfaceList = append(sv.VlanInterface.VlanInterfaceList,
                &sonicpb.SonicVlan_VlanInterface_VlanInterfaceListKey{
                    VlanName: keys[0],
                    VlanInterfaceList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sv)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func VlanInterfaceIPPrefixHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定VLAN Interface或全部VLAN Interface
    kvs := handler.FetchPathKey(r)
    spec := []string{}
    if v, ok := kvs["vlan-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["ip-prefix"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    sv := &sonicpb.SonicVlan{
        VlanInterface: &sonicpb.SonicVlan_VlanInterface{},
    }
    if hkeys, err := conn.GetKeys("VLAN_INTERFACE", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := helper.VlanInterfaceIPPrefix{
                Keys: keys,
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sv.VlanInterface.VlanInterfaceIpprefixList = append(sv.VlanInterface.VlanInterfaceIpprefixList,
                &sonicpb.SonicVlan_VlanInterface_VlanInterfaceIpprefixListKey{
                    VlanName: keys[0],
                    IpPrefix: keys[1],
                    VlanInterfaceIpprefixList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sv)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}