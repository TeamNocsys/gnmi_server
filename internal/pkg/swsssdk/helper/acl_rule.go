package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
    "strings"
)

type AclRule struct {
    Keys []string
    Client command.Client
    Data *sonicpb.SonicAcl_AclRule_AclRuleList
}
// 参考:
// https://github.com/Azure/sonic-swss/blob/master/orchagent/aclorch.h
func (c *AclRule) LoadFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.SonicAcl_AclRule_AclRuleList{}
    }
    if data, err := conn.GetAll(swsssdk.APPL_DB, append([]string{"ACL_RULE_TABLE"}, c.Keys...)); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "PRIORITY":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.Priority = &ywrapper.UintValue{Value: i}
                }
            case "PACKET_ACTION":
                switch strings.ToUpper(v) {
                case "DROP":
                    c.Data.PacketAction = sonicpb.SonicTypesPacketAction_SONICTYPESPACKETACTION_DROP
                case "REDIRECT":
                    c.Data.PacketAction = sonicpb.SonicTypesPacketAction_SONICTYPESPACKETACTION_REDIRECT
                case "FORWARD":
                    c.Data.PacketAction = sonicpb.SonicTypesPacketAction_SONICTYPESPACKETACTION_FORWARD
                }
            case "MIRROR_ACTION":
                c.Data.MirrorAction = &ywrapper.StringValue{Value: v}
            case "REDIRECT_ACTION":
                c.Data.RedirectAction = &ywrapper.StringValue{Value: v}
            case "MIRROR_INGRESS_ACTION":
                c.Data.MirrorIngressAction = &ywrapper.StringValue{Value: v}
            case "MIRROR_EGRESS_ACTION":
                c.Data.MirrorEgressAction = &ywrapper.StringValue{Value: v}
            case "ETHER_TYPE":
                if i, err := strconv.ParseUint(v, 0, 64); err != nil {
                    return err
                } else {
                    switch i {
                    case 0x88CC:
                        c.Data.EtherType = sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_LLDP
                    case 0x8100:
                        c.Data.EtherType = sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_VLAN
                    case 0x8915:
                        c.Data.EtherType = sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_ROCE
                    case 0x0806:
                        c.Data.EtherType = sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_ARP
                    case 0x0800:
                        c.Data.EtherType = sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_IPV4
                    case 0x86DD:
                        c.Data.EtherType = sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_IPV6
                    case 0x8847:
                        c.Data.EtherType = sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_MPLS
                    }
                }
            case "IP_TYPE":
                switch strings.ToUpper(v) {
                case "ANY":
                    c.Data.IpType = sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_ANY
                case "IP":
                    c.Data.IpType = sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_IP
                case "NON_IP":
                    c.Data.IpType = sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_NON_IP
                case "IPV4":
                    c.Data.IpType = sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_IPV4
                case "IPV6":
                    c.Data.IpType = sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_IPV6
                case "IPv4ANY":
                    c.Data.IpType = sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_IPv4ANY
                case "NON_IP4":
                    c.Data.IpType = sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_NON_IP4
                case "IPv6ANY":
                    c.Data.IpType = sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_IPv6ANY
                case "NON_IPv6":
                    c.Data.IpType = sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_NON_IPv6
                case "ARP":
                    c.Data.IpType = sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_ARP
                }
            case "IP_PROTOCOL":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    switch i {
                    case 1:
                        c.Data.IpProtocol = sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_ICMP
                    case 2:
                        c.Data.IpProtocol = sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_IGMP
                    case 6:
                        c.Data.IpProtocol = sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_TCP
                    case 17:
                        c.Data.IpProtocol = sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_UDP
                    case 46:
                        c.Data.IpProtocol = sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_RSVP
                    case 47:
                        c.Data.IpProtocol = sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_GRE
                    case 51:
                        c.Data.IpProtocol = sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_AUTH
                    case 103:
                        c.Data.IpProtocol = sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_PIM
                    case 115:
                        c.Data.IpProtocol = sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_L2TP
                    }
                }
            case "SRC_IP":
                c.Data.SrcIp = &ywrapper.StringValue{Value: v}
            case "DST_IP":
                c.Data.DstIp = &ywrapper.StringValue{Value: v}
            case "SRC_IPV6":
                c.Data.SrcIpv6 = &ywrapper.StringValue{Value: v}
            case "DST_IPV6":
                c.Data.DstIpv6 = &ywrapper.StringValue{Value: v}
            case "L4_SRC_PORT":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.L4SrcPort = &ywrapper.UintValue{Value: i}
                }
            case "L4_DST_PORT":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.L4DstPort = &ywrapper.UintValue{Value: i}
                }
            case "L4_SRC_PORT_RANGE":
                c.Data.L4SrcPortRange = &ywrapper.StringValue{Value: v}
            case "L4_DST_PORT_RANGE":
                c.Data.L4DstPortRange = &ywrapper.StringValue{Value: v}
            case "TCP_FLAGS":
                c.Data.TcpFlags = &ywrapper.StringValue{Value: v}
            case "DSCP":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.Dscp = &ywrapper.UintValue{Value: i}
                }
            }
        }
    }

    return nil
}

func (c *AclRule) SaveToDB(replace bool) error {
    e := make(map[string]interface{})
    if c.Data.Priority != nil {
        e["priority"] = c.Data.Priority.Value
    }
    if c.Data.PacketAction != sonicpb.SonicTypesPacketAction_SONICTYPESPACKETACTION_UNSET {
        switch c.Data.PacketAction {
        case sonicpb.SonicTypesPacketAction_SONICTYPESPACKETACTION_DROP:
            e["PACKET_ACTION"] = "DROP"
        case sonicpb.SonicTypesPacketAction_SONICTYPESPACKETACTION_REDIRECT:
            if c.Data.RedirectAction != nil {
                e["PACKET_ACTION"] = "REDIRECT:" + c.Data.RedirectAction.Value
            }
        case sonicpb.SonicTypesPacketAction_SONICTYPESPACKETACTION_FORWARD:
            e["PACKET_ACTION"] = "FORWARD"
        }
    }
    if c.Data.MirrorAction != nil {
        e["MIRROR_ACTION"] = c.Data.MirrorAction.Value
    }
    if c.Data.MirrorIngressAction != nil {
        e["MIRROR_INGRESS_ACTION"] = c.Data.MirrorIngressAction.Value
    }
    if c.Data.MirrorEgressAction != nil {
        e["MIRROR_EGRESS_ACTION"] = c.Data.MirrorEgressAction.Value
    }
    if c.Data.EtherType != sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_UNSET {
        switch c.Data.EtherType {
        case sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_LLDP:
            e["ETHER_TYPE"] = "0x88CC"
        case sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_VLAN:
            e["ETHER_TYPE"] = "0x8100"
        case sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_ROCE:
            e["ETHER_TYPE"] = "0x8915"
        case sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_ARP:
            e["ETHER_TYPE"] = "0x0806"
        case sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_IPV4:
            e["ETHER_TYPE"] = "0x0800"
        case sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_IPV6:
            e["ETHER_TYPE"] = "0x86DD"
        case sonicpb.SonicAcl_AclRule_AclRuleList_ETHERTYPE_MPLS:
            e["ETHER_TYPE"] = "0x8847"
        }
    }
    if c.Data.IpType != sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_UNSET {
        switch c.Data.IpType {
        case sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_ANY:
            e["ip_type"] = "ANY"
        case sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_IP:
            e["ip_type"] = "IP"
        case sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_NON_IP:
            e["ip_type"] = "NON_IP"
        case sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_IPV4:
            e["ip_type"] = "IPV4"
        case sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_IPV6:
            e["ip_type"] = "IPV6"
        case sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_IPv4ANY:
            e["ip_type"] = "IPv4ANY"
        case sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_NON_IP4:
            e["ip_type"] = "NON_IP4"
        case sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_IPv6ANY:
            e["ip_type"] = "IPv6ANY"
        case sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_NON_IPv6:
            e["ip_type"] = "NON_IPv6"
        case sonicpb.SonicTypesIpType_SONICTYPESIPTYPE_ARP:
            e["ip_type"] = "ARP"
        }
    }
    if c.Data.IpProtocol != sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_UNSET {
        switch c.Data.IpProtocol {
        case sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_ICMP:
            e["IP_PROTOCOL"] = 1
        case sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_IGMP:
            e["IP_PROTOCOL"] = 2
        case sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_TCP:
            e["IP_PROTOCOL"] = 6
        case sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_UDP:
            e["IP_PROTOCOL"] = 17
        case sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_RSVP:
            e["IP_PROTOCOL"] = 46
        case sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_GRE:
            e["IP_PROTOCOL"] = 47
        case sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_AUTH:
            e["IP_PROTOCOL"] = 51
        case sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_PIM:
            e["IP_PROTOCOL"] = 103
        case sonicpb.SonicAcl_AclRule_AclRuleList_IPPROTOCOL_L2TP:
            e["IP_PROTOCOL"] = 115
        }
    }
    if c.Data.SrcIp != nil {
        e["SRC_IP"] = c.Data.SrcIp.Value
    }
    if c.Data.DstIp != nil {
        e["DST_IP"] = c.Data.DstIp.Value
    }
    if c.Data.SrcIpv6 != nil {
        e["SRC_IPV6"] = c.Data.SrcIpv6.Value
    }
    if c.Data.DstIpv6 != nil {
        e["DST_IPV6"] = c.Data.DstIpv6.Value
    }
    if c.Data.L4SrcPort != nil {
        e["L4_SRC_PORT"] = c.Data.L4SrcPort.Value
    }
    if c.Data.L4DstPort != nil {
        e["L4_DST_PORT"] = c.Data.L4DstPort.Value
    }
    if c.Data.L4SrcPortRange != nil {
        e["L4_SRC_PORT_RANGE"] = c.Data.L4SrcPortRange.Value
    }
    if c.Data.L4DstPortRange != nil {
        e["L4_DST_PORT_RANGE"] = c.Data.L4DstPortRange.Value
    }
    if c.Data.TcpFlags != nil {
        e["TCP_FLAGS"] = c.Data.TcpFlags.Value
    }
    if c.Data.Dscp != nil {
        e["DSCP"] = c.Data.Dscp.Value
    }

    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if replace {
        if _, err := conn.SetEntry(swsssdk.CONFIG_DB, append([]string{"ACL_RULE"}, c.Keys...), e); err != nil {
            return err
        }
    } else {
        if _, err := conn.ModEntry(swsssdk.CONFIG_DB, append([]string{"ACL_RULE"}, c.Keys...), e); err != nil {
            return err
        }
    }

    return nil
}

func (c *AclRule) RemoveFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }
    if _, err := conn.DeleteAllByPattern(swsssdk.CONFIG_DB, append([]string{"ACL_RULE"}, c.Keys...)); err != nil {
        return err
    }
    return nil
}