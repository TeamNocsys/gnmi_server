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

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, []string{config_db.VLAN_TABLE, spec})
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    sv := &sonicpb.SonicVlan{
        Vlan: &sonicpb.SonicVlan_Vlan{},
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 2 {
            continue
        }
        name := keys[len(keys)-1]
        v, err := getVlanList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        sv.Vlan.VlanList = append(sv.Vlan.VlanList, &sonicpb.SonicVlan_Vlan_VlanListKey{
            VlanName: name,
            VlanList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sv)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getVlanList(info map[string]string) (*sonicpb.SonicVlan_Vlan_VlanList, error) {
    r := &sonicpb.SonicVlan_Vlan_VlanList{}

    if v, ok := info[config_db.VLAN_DESCRIPTION]; ok {
        r.Description = &ywrapper.StringValue{Value: v}
    } else  {
        r.Description = &ywrapper.StringValue{Value: ""}
    }

    if s, ok := info[config_db.VLAN_DHCP_SERVERS]; ok {
        for _, v := range helper.FieldToArray(s) {
            r.DhcpServers = append(r.DhcpServers, &ywrapper.StringValue{Value: v})
        }
    } else  {
        r.DhcpServers = []*ywrapper.StringValue{}
    }

    if v, ok := info[config_db.VLAN_MTU]; ok {
        if index, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Mtu = &ywrapper.UintValue{Value: index}
        }
    } else  {
        r.Mtu = &ywrapper.UintValue{Value: 0}
    }

    if v, ok := info[config_db.VLAN_ADMIN_STATUS]; ok {
        if strings.ToUpper(v) == config_db.ADMIN_STATUS_UP {
            r.AdminStatus = sonicpb.SonicVlanAdminStatus_SONICVLANADMINSTATUS_up
        } else {
            r.AdminStatus = sonicpb.SonicVlanAdminStatus_SONICVLANADMINSTATUS_down
        }
    } else  {
        r.AdminStatus = sonicpb.SonicVlanAdminStatus_SONICVLANADMINSTATUS_down
    }

    if s, ok := info[config_db.VLAN_MEMBERS]; ok {
        for _, v := range helper.FieldToArray(s) {
            r.Members = append(r.Members, &ywrapper.StringValue{Value: v})
        }
    } else  {
        r.Members = []*ywrapper.StringValue{}
    }

    return r, nil
}
