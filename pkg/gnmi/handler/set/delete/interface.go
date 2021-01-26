package delete

import (
    "context"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func InterfaceHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["port-name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        c := &helper.Interface{
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

func InterfaceIPPrefixHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    spec := []string{}
    if v, ok := kvs["port-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["ip-prefix"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    c := &helper.InterfaceIPPrefix{
        Keys:   spec,
        Client: db,
        Data:   nil,
    }

    if err := c.RemoveFromDB(); err != nil {
        return err
    }
    return nil
}


func LoopbackInterfaceHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["loopback-interface-name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        c := &helper.LoopbackInterface{
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

func LoopbackInterfaceIPPrefixHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    spec := []string{}
    if v, ok := kvs["loopback-interface-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["ip-prefix"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    c := &helper.LoopbackInterfaceIPPrefix{
        Keys:   spec,
        Client: db,
        Data:   nil,
    }

    if err := c.RemoveFromDB(); err != nil {
        return err
    }
    return nil
}


func VlanInterfaceHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["vlan-name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        c := &helper.VlanInterface{
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

func VlanInterfaceIPPrefixHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    spec := []string{}
    if v, ok := kvs["vlan-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["ip-prefix"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    c := &helper.VlanInterfaceIPPrefix{
        Keys:   spec,
        Client: db,
        Data:   nil,
    }

    if err := c.RemoveFromDB(); err != nil {
        return err
    }
    return nil
}