package update

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/golang/protobuf/proto"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk/helper"
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
        if info.Interface != nil {
            if info.Interface.InterfaceList != nil {
                for _, v := range info.Interface.InterfaceList {
                    if v.InterfaceList == nil {
                        continue
                    }
                    c := helper.Interface{
                        Key: v.PortName,
                        Client: db,
                        Data: v.InterfaceList,
                    }
                    c.SaveToDB(false)
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
        if info.Interface != nil {
            if info.Interface.InterfaceIpprefixList != nil {
                for _, v := range info.Interface.InterfaceIpprefixList {
                    if v.InterfaceIpprefixList == nil {
                        continue
                    }
                    c := helper.InterfaceIPPrefix{
                        Keys:   []string{v.PortName, v.IpPrefix},
                        Client: db,
                        Data:   v.InterfaceIpprefixList,
                    }
                    c.SaveToDB(false)
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
        if info.LoopbackInterface != nil {
            if info.LoopbackInterface.LoopbackInterfaceList != nil {
                for _, v := range info.LoopbackInterface.LoopbackInterfaceList {
                    if v.LoopbackInterfaceList == nil {
                        continue
                    }
                    c := helper.LoopbackInterface{
                        Key: v.LoopbackInterfaceName,
                        Client: db,
                        Data: v.LoopbackInterfaceList,
                    }
                    c.SaveToDB(false)
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
        if info.LoopbackInterface != nil {
            if info.LoopbackInterface.LoopbackInterfaceIpprefixList != nil {
                for _, v := range info.LoopbackInterface.LoopbackInterfaceIpprefixList {
                    if v.LoopbackInterfaceIpprefixList == nil {
                        continue
                    }
                    c := helper.LoopbackInterfaceIPPrefix{
                        Keys: []string{v.LoopbackInterfaceName, v.IpPrefix},
                        Client: db,
                        Data: v.LoopbackInterfaceIpprefixList,
                    }
                    c.SaveToDB(false)
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
        if info.VlanInterface != nil {
            if info.VlanInterface.VlanInterfaceList != nil {
                for _, v := range info.VlanInterface.VlanInterfaceList {
                    if v.VlanInterfaceList == nil {
                        continue
                    }
                    c := helper.VlanInterface{
                        Key: v.VlanName,
                        Client: db,
                        Data: v.VlanInterfaceList,
                    }
                    c.SaveToDB(false)
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
        if info.VlanInterface != nil {
            if info.VlanInterface.VlanInterfaceIpprefixList != nil {
                for _, v := range info.VlanInterface.VlanInterfaceIpprefixList {
                    if v.VlanInterfaceIpprefixList == nil {
                        continue
                    }
                    c := helper.VlanInterfaceIPPrefix{
                        Keys:   []string{v.VlanName, v.IpPrefix},
                        Client: db,
                        Data:   v.VlanInterfaceIpprefixList,
                    }
                    c.SaveToDB()
                }
            }
        }
    }

    return nil
}