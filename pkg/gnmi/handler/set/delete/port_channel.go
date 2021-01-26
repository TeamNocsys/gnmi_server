package delete

import (
    "context"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func PortChannelHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["portchannel-name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        c := &helper.PortChannel{
            Key:    v,
            Client: db,
            Data:   nil,
        }
        if err := c.RemoveFromDB(); err != nil {
            return err
        }
    }

    return nil
}

func PortChannelMemberHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    spec := []string{}
    if v, ok := kvs["portchannel-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["port-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    c := &helper.PortChannelMember{
        Keys:   spec,
        Client: db,
        Data:   nil,
    }

    if err := c.RemoveFromDB(); err != nil {
        return err
    }
    return nil
}