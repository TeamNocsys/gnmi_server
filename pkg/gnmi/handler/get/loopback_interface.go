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

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, []string{config_db.LOOPBACK_INTERFACE_TABLE, spec})
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    sli := &sonicpb.SonicLoopbackInterface{
        LoopbackInterface: &sonicpb.SonicLoopbackInterface_LoopbackInterface{},
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 2 {
            continue
        }
        v, err := getLoopbackInterfaceList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        sli.LoopbackInterface.LoopbackInterfaceList = append(sli.LoopbackInterface.LoopbackInterfaceList, &sonicpb.SonicLoopbackInterface_LoopbackInterface_LoopbackInterfaceListKey{
            LoopbackInterfaceName: keys[1],
            LoopbackInterfaceList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sli)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getLoopbackInterfaceList(info map[string]string) (*sonicpb.SonicLoopbackInterface_LoopbackInterface_LoopbackInterfaceList, error) {
    r := &sonicpb.SonicLoopbackInterface_LoopbackInterface_LoopbackInterfaceList{}

    if v, ok := info[config_db.LOOPBACK_INTERFACE_VRF_NAME]; ok {
        r.VrfName = &ywrapper.StringValue{Value: v}
    }

    return r, nil
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

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, append([]string{config_db.LOOPBACK_INTERFACE_TABLE}, spec...))
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    sli := &sonicpb.SonicLoopbackInterface{
        LoopbackInterface: &sonicpb.SonicLoopbackInterface_LoopbackInterface{},
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 3 {
            continue
        }
        v, err := getLoopbackInterfaceIpprefixList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        sli.LoopbackInterface.LoopbackInterfaceIpprefixList = append(sli.LoopbackInterface.LoopbackInterfaceIpprefixList, &sonicpb.SonicLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixListKey{
            LoopbackInterfaceName: keys[1],
            IpPrefix: keys[2],
            LoopbackInterfaceIpprefixList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sli)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getLoopbackInterfaceIpprefixList(info map[string]string) (*sonicpb.SonicLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixList, error) {
    r := &sonicpb.SonicLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixList{}

    if v, ok := info[config_db.LOOPBACK_INTERFACE_IPPREFIX_SCOPE]; ok {
        if strings.ToUpper(v) == config_db.SCOPE_GLOBAL {
            r.Scope = sonicpb.SonicLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixList_SCOPE_global
        } else {
            r.Scope = sonicpb.SonicLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixList_SCOPE_local
        }
    }

    if v, ok := info[config_db.LOOPBACK_INTERFACE_IPPREFIX_FAMILY]; ok {
        if strings.ToUpper(v) == config_db.IP_FAMILY_IPV6 {
            r.Family = sonicpb.SonicLoopbackInterfaceIpFamily_SONICLOOPBACKINTERFACEIPFAMILY_IPv6
        } else {
            r.Family = sonicpb.SonicLoopbackInterfaceIpFamily_SONICLOOPBACKINTERFACEIPFAMILY_IPv4
        }
    }

    return r, nil
}