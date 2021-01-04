package config_db

// CONFIG_DB
const (
    ADMIN_STATUS_UP                     = "UP"
    SCOPE_GLOBAL                        = "GLOBAL"
    IP_FAMILY_IPV6                      = "IPV6"

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
)