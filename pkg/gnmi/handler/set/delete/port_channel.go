package delete

import (
    "context"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func PortChannelHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["portchannel-name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        c := cmd.NewLagAdapter(v, db)
        return c.Config(nil, cmd.DEL)
    }
}

func PortChannelMemberHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    name := "*"
    ifname := "*"
    if v, ok := kvs["portchannel-name"]; ok {
        name = v
    }
    if v, ok := kvs["port-name"]; ok {
        ifname = v
    }

    conn := db.Config()
    if conn == nil {
        return status.Error(codes.Internal, "")
    }

    if hkeys, err := conn.GetKeys("PORTCHANNEL_MEMBER", []string{name, ifname}); err != nil {
        return err
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := cmd.NewLagMemberAdapter(keys[0], keys[1], db)
            if err := c.Config(nil, cmd.DEL); err != nil {
                return err
            }
        }
        return nil
    }
}