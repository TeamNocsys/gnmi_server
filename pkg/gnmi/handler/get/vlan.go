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

func VlanHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定VLAN或全部VLAN
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["vlan-name"]; ok {
        spec = v
    }

    sv := &sonicpb.NocsysVlan{
        Vlan: &sonicpb.NocsysVlan_Vlan{},
    }
    if hkeys, err := conn.GetKeys("VLAN", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := helper.Vlan{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sv.Vlan.VlanList = append(sv.Vlan.VlanList,
                &sonicpb.NocsysVlan_Vlan_VlanListKey{
                    VlanName: keys[0],
                    VlanList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sv)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}


func VlanMemberHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定VLAN Member或全部VLAN Member
    kvs := handler.FetchPathKey(r)
    spec := []string{}
    if v, ok := kvs["vlan-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["port-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    sv := &sonicpb.NocsysVlan{
        VlanMember: &sonicpb.NocsysVlan_VlanMember{},
    }
    if hkeys, err := conn.GetKeys("VLAN_MEMBER", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := helper.VlanMember{
                Keys: keys,
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sv.VlanMember.VlanMemberList = append(sv.VlanMember.VlanMemberList,
                &sonicpb.NocsysVlan_VlanMember_VlanMemberListKey{
                    VlanName: keys[0],
                    Port: keys[1],
                    VlanMemberList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sv)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}