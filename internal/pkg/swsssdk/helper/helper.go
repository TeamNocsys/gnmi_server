package helper

import (
    "fmt"
    "os"
    "strings"
)

const (
    PORT_NAME_PREFIX     = "Ethernet"
    COUNTER_TABLE_PREFIX = "COUNTERS:"
    PORT_STATUS_VALUE_UP = "UP"

    // COUNTERS_DB
    // 表名
    COUNTERS_PORT_NAME_MAP = "COUNTERS_PORT_NAME_MAP"

    // 字段名
    COUNTERS_PORT_IN_UCAST_PKTS          = "SAI_PORT_STAT_IF_IN_UCAST_PKTS"
    COUNTERS_PORT_IN_NON_UCAST_PKTS      = "SAI_PORT_STAT_IF_IN_NON_UCAST_PKTS"
    COUNTERS_PORT_IN_ERRORS              = "SAI_PORT_STAT_IF_IN_ERRORS"
    COUNTERS_PORT_IN_DISCARDS            = "SAI_PORT_STAT_IF_IN_DISCARDS"
    COUNTERS_PORT_ETHER_RX_OVERSIZE_PKTS = "SAI_PORT_STAT_ETHER_RX_OVERSIZE_PKTS"
    COUNTERS_PORT_OUT_UCAST_PKTS         = "SAI_PORT_STAT_IF_OUT_UCAST_PKTS"
    COUNTERS_PORT_OUT_NON_UCAST_PKTS     = "SAI_PORT_STAT_IF_OUT_NON_UCAST_PKTS"
    COUNTERS_PORT_OUT_ERRORS             = "SAI_PORT_STAT_IF_OUT_ERRORS"
    COUNTERS_PORT_OUT_DISCARDS           = "SAI_PORT_STAT_IF_OUT_DISCARDS"
    COUNTERS_PORT_ETHER_TX_OVERSIZE_PKTS = "SAI_PORT_STAT_ETHER_TX_OVERSIZE_PKTS"
    COUNTERS_PORT_IN_OCTETS              = "SAI_PORT_STAT_IF_IN_OCTETS"
    COUNTERS_PORT_OUT_OCTETS             = "SAI_PORT_STAT_IF_OUT_OCTETS"

    // APPL_DB
    // 表名
    PORT_STATUS = "PORT_TABLE"

    // 字段名
    PORT_STATUS_ADMIN_STATUS = "admin_status"
    PORT_STATUS_OPER_STATUS  = "oper_status"
    PORT_STATUS_SPEED        = "speed"
    PORT_STATUS_ALIAS        = "alias"
    PORT_STATUS_MTU          = "mtu"
    PORT_STATUS_INDEX        = "index"
)

var (
    IsAliasModeIface bool
)

func VID(id int) string {
    return fmt.Sprintf("Vlan%d", id)
}

func ArrayToField(s []string) string {
    return strings.Join(s, ",")
}

func FieldToArray(s string) []string {
    return strings.Split(s, ",")
}

func init() {
    IsAliasModeIface = os.Getenv("SONIC_CLI_IFACE_MODE") == "alias"
}
