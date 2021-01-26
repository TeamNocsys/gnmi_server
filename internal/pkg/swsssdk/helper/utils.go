package helper

import (
    "fmt"
    "os"
    "strings"
)

const (
    DATA_TYPE_CONFIG = 1 << 0
    DATA_TYPE_STATE = 1 << 1
    DATA_TYPE_ALL = DATA_TYPE_CONFIG | DATA_TYPE_STATE
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
