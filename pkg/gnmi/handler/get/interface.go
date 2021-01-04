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

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, []string{config_db.INTERFACE_TABLE, spec})
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    si := &sonicpb.SonicInterface{
        Interface: &sonicpb.SonicInterface_Interface{},
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 2 {
            continue
        }
        v, err := getInterfaceList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        si.Interface.InterfaceList = append(si.Interface.InterfaceList, &sonicpb.SonicInterface_Interface_InterfaceListKey{
            PortName: keys[1],
            InterfaceList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, si)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getInterfaceList(info map[string]string) (*sonicpb.SonicInterface_Interface_InterfaceList, error) {
    r := &sonicpb.SonicInterface_Interface_InterfaceList{}

    if v, ok := info[config_db.INTERFACE_VRF_NAME]; ok {
        r.VrfName = &ywrapper.StringValue{Value: v}
    }

    return r, nil
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

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, append([]string{config_db.INTERFACE_TABLE}, spec...))
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    si := &sonicpb.SonicInterface{
        Interface: &sonicpb.SonicInterface_Interface{},
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 3 {
            continue
        }
        v, err := getInterfaceIpprefixList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        si.Interface.InterfaceIpprefixList = append(si.Interface.InterfaceIpprefixList, &sonicpb.SonicInterface_Interface_InterfaceIpprefixListKey{
            PortName: keys[1],
            IpPrefix: keys[2],
            InterfaceIpprefixList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, si)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getInterfaceIpprefixList(info map[string]string) (*sonicpb.SonicInterface_Interface_InterfaceIpprefixList, error) {
    r := &sonicpb.SonicInterface_Interface_InterfaceIpprefixList{}

    if v, ok := info[config_db.INTERFACE_IPPREFIX_SCOPE]; ok {
        if strings.ToUpper(v) == config_db.SCOPE_GLOBAL {
            r.Scope = sonicpb.SonicInterface_Interface_InterfaceIpprefixList_SCOPE_global
        } else {
            r.Scope = sonicpb.SonicInterface_Interface_InterfaceIpprefixList_SCOPE_local
        }
    }

    if v, ok := info[config_db.INTERFACE_IPPREFIX_FAMILY]; ok {
        if strings.ToUpper(v) == config_db.IP_FAMILY_IPV6 {
            r.Family = sonicpb.SonicInterfaceIpFamily_SONICINTERFACEIPFAMILY_IPv6
        } else {
            r.Family = sonicpb.SonicInterfaceIpFamily_SONICINTERFACEIPFAMILY_IPv4
        }
    }

    return r, nil
}