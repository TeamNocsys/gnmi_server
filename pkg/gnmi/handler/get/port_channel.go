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

func PortChannelHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定Port Channel或全部Port Channel
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["portchannel-name"]; ok {
        spec = v
    }

    spc := &sonicpb.NocsysPortchannel{
        Portchannel: &sonicpb.NocsysPortchannel_Portchannel{},
    }
    if hkeys, err := conn.GetKeys("PORTCHANNEL", spec); err != nil {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := helper.PortChannel{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            spc.Portchannel.PortchannelList = append(spc.Portchannel.PortchannelList,
                &sonicpb.NocsysPortchannel_Portchannel_PortchannelListKey{
                    PortchannelName: keys[0],
                    PortchannelList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, spc)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func PortChannelMemberHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定Port Channel Member或全部Port Channel Member
    kvs := handler.FetchPathKey(r)
    spec := []string{}
    if v, ok := kvs["portchannel-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["port-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    // 获取成员信息
    spc := &sonicpb.NocsysPortchannel{
        PortchannelMember: &sonicpb.NocsysPortchannel_PortchannelMember{},
    }
    if hkeys, err := conn.GetKeys("PORTCHANNEL_MEMBER", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := helper.PortChannelMember{
                Keys: keys,
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            spc.PortchannelMember.PortchannelMemberList = append(spc.PortchannelMember.PortchannelMemberList,
                &sonicpb.NocsysPortchannel_PortchannelMember_PortchannelMemberListKey{
                    PortchannelName: keys[0],
                    Port: keys[1],
                    PortchannelMemberList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, spc)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}
