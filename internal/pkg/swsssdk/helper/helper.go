package helper

import (
    "fmt"
    "os"
)

const (
    VLAN_TABLE_NAME                   = "VLAN"
    VLAN_MEMBERTABLE_NAME             = "VLAN_MEMBER"
    INTERFACE_TABLE_NAME              = "INTERFACE"
    PORT_TABLE_NAME                   = "PORT"
    COUNTERS_PORT_NAME_MAP_TABLE_NAME = "COUNTERS_PORT_NAME_MAP"

    VLAN_SUB_INTERFACE_SEPARATOR = "."
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
