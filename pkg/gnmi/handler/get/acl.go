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

func AclRuleHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定ACL Rule或全部ACL Rule
    kvs := handler.FetchPathKey(r)
    var spec []string
    if v, ok := kvs["acl-table-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["rule-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, append([]string{config_db.ACL_RULE_TABLE}, spec...))
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    sa := &sonicpb.SonicAcl{
        AclRule: &sonicpb.SonicAcl_AclRule{},
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 3 {
            continue
        }
        v, err := getAclRuleList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        sa.AclRule.AclRuleList = append(sa.AclRule.AclRuleList, &sonicpb.SonicAcl_AclRule_AclRuleListKey{
            AclTableName: keys[1],
            RuleName: keys[2],
            AclRuleList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sa)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getAclRuleList(info map[string]string) (*sonicpb.SonicAcl_AclRule_AclRuleList, error) {
    r := &sonicpb.SonicAcl_AclRule_AclRuleList{}

    if v, ok := info[config_db.ACL_RULE_PACKET_ACTION]; ok {
        switch strings.ToUpper(v) {
        case config_db.PACKETACTION_FORWARD:
            r.PacketAction = sonicpb.SonicAclPacketAction_SONICACLPACKETACTION_FORWARD
            break
        case config_db.PACKETACTION_REDIRECT:
            r.PacketAction = sonicpb.SonicAclPacketAction_SONICACLPACKETACTION_REDIRECT
            break
        default:
            r.PacketAction = sonicpb.SonicAclPacketAction_SONICACLPACKETACTION_DROP
        }
    }

    if v, ok := info[config_db.ACL_RULE_IP_TYPE]; ok {
        if index, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            switch index {
            case uint64(sonicpb.SonicAclIpType_SONICACLIPTYPE_ANY.Number()):
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_ANY
                break
            case uint64(sonicpb.SonicAclIpType_SONICACLIPTYPE_IP.Number()):
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_IP
                break
            case uint64(sonicpb.SonicAclIpType_SONICACLIPTYPE_NON_IP.Number()):
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_NON_IP
                break
            case uint64(sonicpb.SonicAclIpType_SONICACLIPTYPE_IPV4.Number()):
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_IPV4
                break
            case uint64(sonicpb.SonicAclIpType_SONICACLIPTYPE_IPV6.Number()):
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_IPV6
                break
            case uint64(sonicpb.SonicAclIpType_SONICACLIPTYPE_IPv4ANY.Number()):
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_IPv4ANY
                break
            case uint64(sonicpb.SonicAclIpType_SONICACLIPTYPE_NON_IP4.Number()):
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_NON_IP4
                break
            case uint64(sonicpb.SonicAclIpType_SONICACLIPTYPE_IPv6ANY.Number()):
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_IPv6ANY
                break
            case uint64(sonicpb.SonicAclIpType_SONICACLIPTYPE_NON_IPv6.Number()):
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_NON_IPv6
                break
            case uint64(sonicpb.SonicAclIpType_SONICACLIPTYPE_ARP.Number()):
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_ARP
                break
            default:
                r.IpType = sonicpb.SonicAclIpType_SONICACLIPTYPE_UNSET
            }
        }
    }

    if v, ok := info[config_db.ACL_RULE_PRIORITY]; ok {
        if index, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Priority = &ywrapper.UintValue{Value: index}
        }
    }

    if r.IpType == sonicpb.SonicAclIpType_SONICACLIPTYPE_ANY ||
        r.IpType == sonicpb.SonicAclIpType_SONICACLIPTYPE_IP ||
        r.IpType == sonicpb.SonicAclIpType_SONICACLIPTYPE_IPV4 ||
        r.IpType == sonicpb.SonicAclIpType_SONICACLIPTYPE_IPv4ANY ||
        r.IpType == sonicpb.SonicAclIpType_SONICACLIPTYPE_ARP {
        if v, ok := info[config_db.ACL_RULE_SRC_IP]; ok {
            r.SrcIp = &ywrapper.StringValue{Value: v}
        }
        if v, ok := info[config_db.ACL_RULE_DST_IP]; ok {
            r.SrcIp = &ywrapper.StringValue{Value: v}
        }
        if v, ok := info[config_db.ACL_RULE_L4_DST_PORT]; ok {
            if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                return nil, err
            } else {
                r.L4DstPort = &ywrapper.UintValue{Value: i}
            }
        }
        if v, ok := info[config_db.ACL_RULE_ICMP_TYPE]; ok {
            if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                return nil, err
            } else {
                r.IcmpType = &ywrapper.UintValue{Value: i}
            }
        }
        if v, ok := info[config_db.ACL_RULE_ICMP_CODE]; ok {
            if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                return nil, err
            } else {
                r.IcmpCode = &ywrapper.UintValue{Value: i}
            }
        }
    }

    if r.IpType == sonicpb.SonicAclIpType_SONICACLIPTYPE_ANY ||
        r.IpType == sonicpb.SonicAclIpType_SONICACLIPTYPE_IP ||
        r.IpType == sonicpb.SonicAclIpType_SONICACLIPTYPE_IPV6 ||
        r.IpType == sonicpb.SonicAclIpType_SONICACLIPTYPE_IPv6ANY {
        if v, ok := info[config_db.ACL_RULE_SRC_IP6]; ok {
            r.SrcIp6 = &ywrapper.StringValue{Value: v}
        }
        if v, ok := info[config_db.ACL_RULE_DST_IP6]; ok {
            r.SrcIp6 = &ywrapper.StringValue{Value: v}
        }
        if v, ok := info[config_db.ACL_RULE_ICMPV6_TYPE]; ok {
            if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                return nil, err
            } else {
                r.Icmpv6Type = &ywrapper.UintValue{Value: i}
            }
        }
        if v, ok := info[config_db.ACL_RULE_ICMPV6_CODE]; ok {
            if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                return nil, err
            } else {
                r.Icmpv6Code = &ywrapper.UintValue{Value: i}
            }
        }
    }

    if v, ok := info[config_db.ACL_RULE_IN_PORTS]; ok {
        for _, s := range helper.FieldToArray(v) {
            if i, err := strconv.ParseUint(s, 10, 64); err != nil {
                return nil, err
            } else {
                r.InPorts = append(r.InPorts, &ywrapper.UintValue{Value: i})
            }
        }
    }

    if v, ok := info[config_db.ACL_RULE_OUT_PORTS]; ok {
        for _, s := range helper.FieldToArray(v) {
            if i, err := strconv.ParseUint(s, 10, 64); err != nil {
                return nil, err
            } else {
                r.OutPorts = append(r.OutPorts, &ywrapper.UintValue{Value: i})
            }
        }
    }

    if v, ok := info[config_db.ACL_RULE_L4_SRC_PORT]; ok {
        if i, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.L4SrcPort = &ywrapper.UintValue{Value: i}
        }
    }

    if v, ok := info[config_db.ACL_RULE_L4_SRC_PORT_RANGE]; ok {
        r.L4SrcPortRange = &ywrapper.StringValue{Value: v}
    }

    if v, ok := info[config_db.ACL_RULE_L4_DST_PORT]; ok {
        if i, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.L4DstPort = &ywrapper.UintValue{Value: i}
        }
    }

    /*if v, ok := info[config_db.ACL_RULE_L4_DST_PORT_RANGE]; ok {
        r.L4DstPortRange = &ywrapper.StringValue{Value: v}
    }*/

    if v, ok := info[config_db.ACL_RULE_ETHER_TYPE]; ok {
        r.EtherType = &ywrapper.StringValue{Value: v}
    }

    if v, ok := info[config_db.ACL_RULE_IP_PROTOCOL]; ok {
        if i, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.IpProtocol = &ywrapper.UintValue{Value: i}
        }
    }

    if v, ok := info[config_db.ACL_RULE_TCP_FLAGS]; ok {
        r.TcpFlags = &ywrapper.StringValue{Value: v}
    }

    if v, ok := info[config_db.ACL_RULE_DSCP]; ok {
        if i, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Dscp = &ywrapper.UintValue{Value: i}
        }
    }

    if v, ok := info[config_db.ACL_RULE_TC]; ok {
        if i, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Tc = &ywrapper.UintValue{Value: i}
        }
    }

    if v, ok := info[config_db.ACL_RULE_INNER_ETHER_TYPE]; ok {
        r.InnerEtherType = &ywrapper.StringValue{Value: v}
    }

    if v, ok := info[config_db.ACL_RULE_INNER_IP_PROTOCOL]; ok {
        if i, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.InnerIpProtocol = &ywrapper.UintValue{Value: i}
        }
    }

    if v, ok := info[config_db.ACL_RULE_INNER_L4_SRC_PORT]; ok {
        if i, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.InnerL4SrcPort = &ywrapper.UintValue{Value: i}
        }
    }

    if v, ok := info[config_db.ACL_RULE_INNER_L4_DST_PORT]; ok {
        if i, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.InnerL4DstPort = &ywrapper.UintValue{Value: i}
        }
    }

    return r, nil
}

func AclTableHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定ACL Table或全部ACL Table
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["acl-table-name"]; ok {
        spec = v
    }

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, []string{config_db.ACL_TABLE, spec})
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    sa := &sonicpb.SonicAcl{
        AclTable: &sonicpb.SonicAcl_AclTable{},
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 2 {
            continue
        }
        v, err := getAclTableList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        sa.AclTable.AclTableList = append(sa.AclTable.AclTableList, &sonicpb.SonicAcl_AclTable_AclTableListKey{
            AclTableName: keys[1],
            AclTableList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sa)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getAclTableList(info map[string]string) (*sonicpb.SonicAcl_AclTable_AclTableList, error) {
    r := &sonicpb.SonicAcl_AclTable_AclTableList{}

    if v, ok := info[config_db.ACL_TABLE_POLICY_DESC]; ok {
        r.PolicyDesc = &ywrapper.StringValue{Value: v}
    }

    if v, ok := info[config_db.ACL_TABLE_TYPE]; ok {
        if i, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            switch i {
            case uint64(sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_L2.Number()):
                r.Type = sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_L2
                break
            case uint64(sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_L3.Number()):
                r.Type = sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_L3
                break
            case uint64(sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_L3V6.Number()):
                r.Type = sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_L3V6
                break
            case uint64(sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_MIRROR.Number()):
                r.Type = sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_MIRROR
                break
            case uint64(sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_MIRRORV6.Number()):
                r.Type = sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_MIRRORV6
                break
            case uint64(sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_MIRROR_DSCP.Number()):
                r.Type = sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_MIRROR_DSCP
                break
            case uint64(sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_CTRLPLANE.Number()):
                r.Type = sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_CTRLPLANE
                break
            default:
                r.Type = sonicpb.SonicAclAclTableType_SONICACLACLTABLETYPE_UNSET
            }
        }
    }

    if v, ok := info[config_db.ACL_TABLE_STAGE]; ok {
        if i, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            switch i {
            case uint64(sonicpb.SonicAcl_AclTable_AclTableList_STAGE_INGRESS.Number()):
                r.Stage = sonicpb.SonicAcl_AclTable_AclTableList_STAGE_INGRESS
                break
            case uint64(sonicpb.SonicAcl_AclTable_AclTableList_STAGE_EGRESS.Number()):
                r.Stage = sonicpb.SonicAcl_AclTable_AclTableList_STAGE_EGRESS
                break
            default:
                r.Stage = sonicpb.SonicAcl_AclTable_AclTableList_STAGE_UNSET
            }
        }

    }

    if v, ok := info[config_db.ACL_TABLE_PORTS]; ok {
        for _, s := range helper.FieldToArray(v) {
            r.Ports = append(r.Ports, &ywrapper.StringValue{Value: s})
        }
    }

    return r, nil
}