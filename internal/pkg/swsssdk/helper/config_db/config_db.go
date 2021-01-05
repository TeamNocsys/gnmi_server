package config_db

// CONFIG_DB
const (
    ADMIN_STATUS_UP                     = "UP"
    SCOPE_GLOBAL                        = "GLOBAL"
    IP_FAMILY_IPV6                      = "IPV6"
    PACKETACTION_DROP                   = "DROP"
    PACKETACTION_FORWARD                = "FORWARD"
    PACKETACTION_REDIRECT               = "REDIRECT"

    // 表名
    VLAN_TABLE                          = "VLAN"
    // 字段名
    VLAN_TABLE_VLANID                   = "vlanid"
    VLAN_TABLE_DESCRIPTION              = "description"
    VLAN_TABLE_DHCP_SERVERS             = "dhcp-servers"
    VLAN_TABLE_MTU                      = "mtu"
    VLAN_TABLE_ADMIN_STATUS             = "admin_status"
    VLAN_TABLE_MEMBERS                  = "members"

    // 表名
    VLAN_MEMBER_TABLE                   = "VLAN_MEMBER"
    // 字段名
    VLAN_MEMBER_TAGGING_MODE            = "tagging_mode"

    // 表名
    VLAN_INTERFACE_TABLE                = "VLAN_INTERFACE"
    // 字段名
    VLAN_INTERFACE_VRF_NAME             = "vrf_name"
    VLAN_INTERFACE_IPPREFIX_SCOPE       = "scope"
    VLAN_INTERFACE_IPPREFIX_FAMILY      = "family"

    // 表名
    INTERFACE_TABLE                     = "INTERFACE"
    // 字段名
    INTERFACE_VRF_NAME                  = "vrf_name"
    INTERFACE_IPPREFIX_SCOPE            = "scope"
    INTERFACE_IPPREFIX_FAMILY           = "family"

    // 表名
    PORTCHANNEL_TABLE                   = "PORTCHANNEL"
    // 字段名
    PORTCHANNEL_MEMBERS                 = "members"
    PORTCHANNEL_MIN_LINKS               = "min_links"
    PORTCHANNEL_DESCRIPTION             = "description"
    PORTCHANNEL_MTU                     = "mtu"
    PORTCHANNEL_ADMIN_STATUS            = "admin_status"

    // 表名
    PORT_TABLE                          = "PORT"
    // 字段名
    PORT_ALIAS                          = "alias"
    PORT_LANES                          = "lanes"
    PORT_DESCRIPTION                    = "description"
    PORT_SPEED                          = "speed"
    PORT_MTU                            = "mtu"
    PORT_INDEX                          = "index"
    PORT_ADMIN_STATUS                   = "admin_status"
    PORT_FEC                            = "fec"

    // 表名
    ACL_RULE_TABLE                      = "ACL_RULE"
    // 字段名
    /*ACL_RULE_PACKET_ACTION              = "packet_action"
    ACL_RULE_IP_TYPE                    = "ip_type"
    ACL_RULE_PRIORITY                   = "priority"
    ACL_RULE_SRC_IP                     = "src_ip"
    ACL_RULE_DST_IP                     = "dst_ip"
    ACL_RULE_SRC_IP6                    = "src_ip6"
    ACL_RULE_DST_IP6                    = "dst_ip6"
    ACL_RULE_IN_PORTS                   = "in_ports"
    ACL_RULE_OUT_PORTS                  = "out_ports"
    ACL_RULE_L4_SRC_PORT                = "l4_src_port"
    ACL_RULE_L4_SRC_PORT_RANGE          = "l4_src_port_range"
    ACL_RULE_L4_DST_PORT                = "l4_dst_port"
    ACL_RULE_L4_DST_PORT_RANGE          = "l4_dst_port_range"
    ACL_RULE_ETHER_TYPE                 = "ether_type"
    ACL_RULE_IP_PROTOCOL                = "ip_protocol"
    ACL_RULE_TCP_FLAGS                  = "tcp_flags"
    ACL_RULE_DSCP                       = "dscp"
    ACL_RULE_TC                         = "tc"
    ACL_RULE_ICMP_TYPE                  = "icmp_type"
    ACL_RULE_ICMP_CODE                  = "icmp_code"
    ACL_RULE_ICMPV6_TYPE                = "icmpv6_type"
    ACL_RULE_ICMPV6_CODE                = "icmpv6_code"
    ACL_RULE_INNER_ETHER_TYPE           = "inner_ether_type"
    ACL_RULE_INNER_IP_PROTOCOL          = "inner_ip_protocol"
    ACL_RULE_INNER_L4_SRC_PORT          = "inner_l4_src_port"
    ACL_RULE_INNER_L4_DST_PORT          = "inner_l4_dst_port"*/
    ACL_RULE_PACKET_ACTION              = "PACKET_ACTION"
    ACL_RULE_IP_TYPE                    = "IP_TYPE"
    ACL_RULE_PRIORITY                   = "PRIORITY"
    ACL_RULE_SRC_IP                     = "SRC_IP"
    ACL_RULE_DST_IP                     = "DST_IP"
    ACL_RULE_SRC_IP6                    = "SRC_IP6"
    ACL_RULE_DST_IP6                    = "DST_IP6"
    ACL_RULE_IN_PORTS                   = "IN_PORTS"
    ACL_RULE_OUT_PORTS                  = "OUT_PORTS"
    ACL_RULE_L4_SRC_PORT                = "L4_SRC_PORT"
    ACL_RULE_L4_SRC_PORT_RANGE          = "L4_SRC_PORT_RANGE"
    ACL_RULE_L4_DST_PORT                = "L4_DST_PORT"
    ACL_RULE_L4_DST_PORT_RANGE          = "L4_DST_PORT_RANGE"
    ACL_RULE_ETHER_TYPE                 = "ETHER_TYPE"
    ACL_RULE_IP_PROTOCOL                = "IP_PROTOCOL"
    ACL_RULE_TCP_FLAGS                  = "TCP_FLAGS"
    ACL_RULE_DSCP                       = "DSCP"
    ACL_RULE_TC                         = "TC"
    ACL_RULE_ICMP_TYPE                  = "ICMP_TYPE"
    ACL_RULE_ICMP_CODE                  = "ICMP_CODE"
    ACL_RULE_ICMPV6_TYPE                = "ICMPV6_TYPE"
    ACL_RULE_ICMPV6_CODE                = "ICMPV6_CODE"
    ACL_RULE_INNER_ETHER_TYPE           = "INNER_ETHER_TYPE"
    ACL_RULE_INNER_IP_PROTOCOL          = "INNER_IP_PROTOCOL"
    ACL_RULE_INNER_L4_SRC_PORT          = "INNER_L4_SRC_PORT"
    ACL_RULE_INNER_L4_DST_PORT          = "INNER_L4_DST_PORT"

    // 表名
    ACL_TABLE                           = "ACL_TABLE"
    // 字段名
    ACL_TABLE_POLICY_DESC               = "policy_desc"
    ACL_TABLE_TYPE                      = "type"
    ACL_TABLE_STAGE                     = "stage"
    ACL_TABLE_PORTS                     = "ports"

    // 表名
    LOOPBACK_INTERFACE_TABLE            = "LOOPBACK_INTERFACE"
    // 字段名
    LOOPBACK_INTERFACE_VRF_NAME         = "vrf_name"
    LOOPBACK_INTERFACE_IPPREFIX_SCOPE   = "scope"
    LOOPBACK_INTERFACE_IPPREFIX_FAMILY  = "family"
)