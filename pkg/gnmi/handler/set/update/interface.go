package update

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/golang/protobuf/jsonpb"
    "github.com/golang/protobuf/proto"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func InterfaceHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysInterface{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("UPDATE|%s", s)
        if info.Interface != nil {
            if info.Interface.InterfaceList != nil {
                for _, v := range info.Interface.InterfaceList {
                    if v.InterfaceList == nil {
                        continue
                    }
                    c := cmd.NewIfAdapter(cmd.INTERFACE, v.PortName, db)
                    if err := c.Config(v.InterfaceList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}

func InterfaceIPPrefixHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysInterface{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("UPDATE|%s", s)
        if info.Interface != nil {
            if info.Interface.InterfaceIpprefixList != nil {
                for _, v := range info.Interface.InterfaceIpprefixList {
                    if v.InterfaceIpprefixList == nil {
                        continue
                    }
                    c := cmd.NewIfAddrAdapter(cmd.INTERFACE, v.PortName, v.IpPrefix, db)
                    if err := c.Config(v.InterfaceIpprefixList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}

func LoopbackInterfaceHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysLoopbackInterface{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("UPDATE|%s", s)
        if info.LoopbackInterface != nil {
            if info.LoopbackInterface.LoopbackInterfaceList != nil {
                for _, v := range info.LoopbackInterface.LoopbackInterfaceList {
                    if v.LoopbackInterfaceList == nil {
                        continue
                    }
                    c := cmd.NewIfAdapter(cmd.LOOPBACK_INTERFACE, v.LoopbackInterfaceName, db)
                    if err := c.Config(v.LoopbackInterfaceList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}

func LoopbackInterfaceIPPrefixHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysLoopbackInterface{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("UPDATE|%s", s)
        if info.LoopbackInterface != nil {
            if info.LoopbackInterface.LoopbackInterfaceIpprefixList != nil {
                for _, v := range info.LoopbackInterface.LoopbackInterfaceIpprefixList {
                    if v.LoopbackInterfaceIpprefixList == nil {
                        continue
                    }
                    c := cmd.NewIfAddrAdapter(cmd.LOOPBACK_INTERFACE, v.LoopbackInterfaceName, v.IpPrefix, db)
                    if err := c.Config(v.LoopbackInterfaceIpprefixList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}

func VlanInterfaceHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysVlan{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("UPDATE|%s", s)
        if info.VlanInterface != nil {
            if info.VlanInterface.VlanInterfaceList != nil {
                for _, v := range info.VlanInterface.VlanInterfaceList {
                    if v.VlanInterfaceList == nil {
                        continue
                    }
                    c := cmd.NewIfAdapter(cmd.VLAN_INTERFACE, v.VlanName, db)
                    if err := c.Config(v.VlanInterfaceList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}

func VlanInterfaceIPPrefixHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysVlan{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("UPDATE|%s", s)
        if info.VlanInterface != nil {
            if info.VlanInterface.VlanInterfaceIpprefixList != nil {
                for _, v := range info.VlanInterface.VlanInterfaceIpprefixList {
                    if v.VlanInterfaceIpprefixList == nil {
                        continue
                    }
                    c := cmd.NewIfAddrAdapter(cmd.VLAN_INTERFACE, v.VlanName, v.IpPrefix, db)
                    if err := c.Config(v.VlanInterfaceIpprefixList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}