package delete

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func VlanHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["vlan-name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        conn := db.Config()
        if conn == nil {
            return status.Error(codes.Internal, "")
        }

        // 如果有vrf删除vrf
        if data, err := conn.GetEntry("VLAN_INTERFACE", v); err != nil {
            return err
        } else {
            if vrf, ok := data["vrf_name"]; ok {
                c := cmd.NewIfAdapter(cmd.VLAN_INTERFACE, v, db)
                if err := c.Config(&sonicpb.NocsysVlan_VlanInterface_VlanInterfaceList{
                    VrfName: &ywrapper.StringValue{Value: vrf.(string)},
                }, cmd.DEL); err != nil {
                    return err
                }
            } else {
                // 删除IP
                if hkeys, err := conn.GetKeysWithTrace(cmd.IfType_table[int32(cmd.VLAN_INTERFACE)], []string{v, "*"}); err != nil {
                    return err
                } else {
                    for _, hkey := range hkeys {
                        keys := conn.SplitKeys(hkey)
                        c := cmd.NewIfAddrAdapter(cmd.VLAN_INTERFACE, keys[0], keys[1], db)
                        if err := c.Config(nil, cmd.DEL); err != nil {
                            return err
                        }
                    }
                }
            }
        }

        // 删除成员
        if hkeys, err := conn.GetKeysWithTrace("VLAN_MEMBER", []string{v, "*"}); err != nil {
            return err
        } else {
            for _, hkey := range hkeys {
                keys := conn.SplitKeys(hkey)
                c := cmd.NewVlanMemberAdapter(keys[0], keys[1], db)
                if err := c.Config(nil, cmd.DEL); err != nil {
                    return err
                }
            }
        }

        // 最后删除VLAN
        c := cmd.NewVlanAdapter(v, db)
        return c.Config(nil, cmd.DEL)
    }
}

func VlanMemberHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    name := "*"
    ifname := "*"
    if v, ok := kvs["vlan-name"]; ok {
        name = v
    }
    if v, ok := kvs["port-name"]; ok {
        ifname = v
    }

    conn := db.Config()
    if conn == nil {
        return status.Error(codes.Internal, "")
    }

    if hkeys, err := conn.GetKeysWithTrace("VLAN_MEMBER", []string{name, ifname}); err != nil {
        return err
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := cmd.NewVlanMemberAdapter(keys[0], keys[1], db)
            if err := c.Config(nil, cmd.DEL); err != nil {
                return err
            }
        }
        return nil
    }
}