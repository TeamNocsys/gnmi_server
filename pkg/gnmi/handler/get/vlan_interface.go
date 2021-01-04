package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "gnmi_server/pkg/gnmi/handler"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "strings"
)

func VlanInterfaceHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定VLAN Interface或全部VLAN Interface
    kvs := handler.FetchPathKey(r)
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    spec := "*"
    if v, ok := kvs["vlan-name"]; ok {
        spec = v
    }

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, []string{config_db.VLAN_INTERFACE_TABLE, spec})
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    sv := &sonicpb.SonicVlan{
        VlanInterface: &sonicpb.SonicVlan_VlanInterface{},
    }
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 2 {
            continue
        }
        v, err := getVlanInterfaceList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        sv.VlanInterface.VlanInterfaceList = append(sv.VlanInterface.VlanInterfaceList, &sonicpb.SonicVlan_VlanInterface_VlanInterfaceListKey{
            VlanName: keys[1],
            VlanInterfaceList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sv)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getVlanInterfaceList(info map[string]string) (*sonicpb.SonicVlan_VlanInterface_VlanInterfaceList, error) {
    r := &sonicpb.SonicVlan_VlanInterface_VlanInterfaceList{}

    if v, ok := info[config_db.VLAN_INTERFACE_VRF_NAME]; ok {
        r.VrfName = &ywrapper.StringValue{Value: v}
    }

    return r, nil
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

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, append([]string{config_db.VLAN_MEMBER_TABLE}, spec...))
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    sv := &sonicpb.SonicVlan{
        VlanInterface: &sonicpb.SonicVlan_VlanInterface{},
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 3 {
            continue
        }
        v, err := getVlanInterfaceIpprefixList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        sv.VlanInterface.VlanInterfaceIpprefixList = append(sv.VlanInterface.VlanInterfaceIpprefixList, &sonicpb.SonicVlan_VlanInterface_VlanInterfaceIpprefixListKey{
            VlanName: keys[1],
            IpPrefix: keys[2],
            VlanInterfaceIpprefixList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sv)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getVlanInterfaceIpprefixList(info map[string]string) (*sonicpb.SonicVlan_VlanInterface_VlanInterfaceIpprefixList, error) {
    r := &sonicpb.SonicVlan_VlanInterface_VlanInterfaceIpprefixList{}

    if v, ok := info[config_db.VLAN_INTERFACE_IPPREFIX_SCOPE]; ok {
        if strings.ToUpper(v) == config_db.SCOPE_GLOBAL {
            r.Scope = sonicpb.SonicVlan_VlanInterface_VlanInterfaceIpprefixList_SCOPE_global
        } else {
            r.Scope = sonicpb.SonicVlan_VlanInterface_VlanInterfaceIpprefixList_SCOPE_local
        }
    }

    if v, ok := info[config_db.VLAN_INTERFACE_IPPREFIX_FAMILY]; ok {
        if strings.ToUpper(v) == config_db.IP_FAMILY_IPV6 {
            r.Family = sonicpb.SonicVlanIpFamily_SONICVLANIPFAMILY_IPv6
        } else {
            r.Family = sonicpb.SonicVlanIpFamily_SONICVLANIPFAMILY_IPv4
        }
    }

    return r, nil
}