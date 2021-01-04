package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "gnmi_server/pkg/gnmi/handler"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "strconv"
    "strings"
)

func PortChannelHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定Port Channel或全部Port Channel
    kvs := handler.FetchPathKey(r)
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    spec := "*"
    if v, ok := kvs["portchannel-name"]; ok {
        spec = v
    }

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, []string{config_db.PORTCHANNEL_TABLE, spec})
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    spc := &sonicpb.SonicPortchannel{
        Portchannel: &sonicpb.SonicPortchannel_Portchannel{},
    }
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 2 {
            continue
        }
        v, err := getPortchannelList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        spc.Portchannel.PortchannelList = append(spc.Portchannel.PortchannelList, &sonicpb.SonicPortchannel_Portchannel_PortchannelListKey{
            PortchannelName: keys[1],
            PortchannelList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, spc)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getPortchannelList(info map[string]string) (*sonicpb.SonicPortchannel_Portchannel_PortchannelList, error) {
    r := &sonicpb.SonicPortchannel_Portchannel_PortchannelList{}

    if s, ok := info[config_db.PORTCHANNEL_MEMBERS]; ok {
        for _, v := range helper.FieldToArray(s) {
            r.Members = append(r.Members, &ywrapper.StringValue{Value: v})
        }
    } else  {
        r.Members = []*ywrapper.StringValue{}
    }

    if v, ok := info[config_db.PORTCHANNEL_MIN_LINKS]; ok {
        if index, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.MinLinks = &ywrapper.UintValue{Value: index}
        }
    }

    if v, ok := info[config_db.PORTCHANNEL_DESCRIPTION]; ok {
        r.Description = &ywrapper.StringValue{Value: v}
    } else  {
        r.Description = &ywrapper.StringValue{Value: ""}
    }

    if v, ok := info[config_db.PORTCHANNEL_MTU]; ok {
        if index, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Mtu = &ywrapper.UintValue{Value: index}
        }
    } else  {
        r.Mtu = &ywrapper.UintValue{Value: 0}
    }

    if v, ok := info[config_db.PORTCHANNEL_ADMIN_STATUS]; ok {
        if strings.ToUpper(v) == config_db.ADMIN_STATUS_UP {
            r.AdminStatus = sonicpb.SonicPortchannelAdminStatus_SONICPORTCHANNELADMINSTATUS_up
        } else {
            r.AdminStatus = sonicpb.SonicPortchannelAdminStatus_SONICPORTCHANNELADMINSTATUS_down
        }
    } else  {
        r.AdminStatus = sonicpb.SonicPortchannelAdminStatus_SONICPORTCHANNELADMINSTATUS_down
    }

    return r, nil
}
