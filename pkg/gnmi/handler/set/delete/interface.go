package delete

import (
    "context"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func InterfaceHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["port-name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        return ifAddrAutoRemove(cmd.INTERFACE, v, "*", db)
    }
}

func InterfaceIPPrefixHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    ifname := "*"
    ipaddr := "*"
    if v, ok := kvs["port-name"]; ok {
        ifname = v
    }
    if v, ok := kvs["ip-prefix"]; ok {
        ipaddr = v
    }
    return ifAddrAutoRemove(cmd.INTERFACE, ifname, ipaddr, db)
}

func LoopbackInterfaceHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["loopback-interface-name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        return ifAddrAutoRemove(cmd.LOOPBACK_INTERFACE, v, "*", db)
    }
}

func LoopbackInterfaceIPPrefixHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    ifname := "*"
    ipaddr := "*"
    if v, ok := kvs["loopback-interface-name"]; ok {
        ifname = v
    }
    if v, ok := kvs["ip-prefix"]; ok {
        ipaddr = v
    }
    return ifAddrAutoRemove(cmd.LOOPBACK_INTERFACE, ifname, ipaddr, db)
}


func VlanInterfaceHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["vlan-name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        return ifAddrAutoRemove(cmd.VLAN_INTERFACE, v, "*", db)
    }
}

func VlanInterfaceIPPrefixHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    ifname := "*"
    ipaddr := "*"
    if v, ok := kvs["vlan-name"]; ok {
        ifname = v
    }
    if v, ok := kvs["ip-prefix"]; ok {
        ipaddr = v
    }
    return ifAddrAutoRemove(cmd.VLAN_INTERFACE, ifname, ipaddr, db)
}

func ifAddrAutoRemove(ifType cmd.IfType, ifname string, ipaddr string, db command.Client) error {
    conn := db.Config()
    if conn == nil {
        return status.Error(codes.Internal, "")
    }

    if hkeys, err := conn.GetKeys(cmd.IfType_table[int32(ifType)], []string{ifname, ipaddr}); err != nil {
        return err
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := cmd.NewIfAddrAdapter(ifType, keys[0], keys[1], db)
            if err := c.Config(nil, cmd.DEL); err != nil {
                return err
            }
        }
        return nil
    }
}