package helper

import (
    "fmt"
    "os"
)

const (
    VLAN_TABLE_NAME                   = "VLAN"
    VLAN_MEMBERTABLE_NAME             = "VLAN_MEMBER"
    INTERFACE_TABLE_NAME              = "INTERFACE"
    COUNTERS_PORT_NAME_MAP_TABLE_NAME = "COUNTERS_PORT_NAME_MAP"

    VLAN_SUB_INTERFACE_SEPARATOR = "."

    PORT_NAME_PREFIX = "Ethernet"
    PORT_STATUS_VALUE_UP = "UP"

    // CONFIG_DB
    // 表名

    // 字段名


    // COUNTERS_DB
    // 表名

    // 字段名

    // APPL_DB
    // 表名
    PORT_STATUS_TABLE_NAME = "PORT_TABLE"

    // 字段名
    PORT_STATUS_ADMIN_STATUS_FIELD = "admin_status"
    PORT_STATUS_OPER_STATUS_FIELD = "oper_status"
    PORT_STATUS_SPEED_FIELD = "speed"
    PORT_STATUS_ALIAS_FIELD = "alias"
    PORT_STATUS_MTU_FIELD = "mtu"
    PORT_STATUS_INDEX_FIELD = "index"
)

var (
    IsAliasModeIface bool
)

func VID(id int) string {
    return fmt.Sprintf("Vlan%d", id)
}

func init() {
    IsAliasModeIface = os.Getenv("SONIC_CLI_IFACE_MODE") == "alias"
}
