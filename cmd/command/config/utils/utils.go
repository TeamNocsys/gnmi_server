package utils

import (
    "fmt"
)

const (
    PORT_PATTERN = "Ethernet[0-9]+"
    PORT_CHANNEL_PATTERN = "PortChannel[0-9]{1,4}"
    VLAN_PATTERN = "Vlan[a-zA-Z0-9_-]+"
    LOOPBACK_PATTERN = "Loopback[0-9]+"
)

var (
    ErrUnknowInterface = fmt.Errorf("the format of interface name is not correct")
)