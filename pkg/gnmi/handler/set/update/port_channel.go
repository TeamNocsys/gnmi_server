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

func PortChannelHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.SonicPortchannel{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.Portchannel != nil {
            if info.Portchannel.PortchannelList != nil {
                for _, v := range info.Portchannel.PortchannelList {
                    if v.PortchannelList == nil {
                        continue
                    }
                    c := helper.PortChannel{
                        Key: v.PortchannelName,
                        Client: db,
                        Data: v.PortchannelList,
                    }
                    c.SaveToDB(false)
                }
            }
        }
    }

    return nil
}


func PortChannelMemberHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.SonicPortchannel{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.PortchannelMember != nil {
            if info.PortchannelMember.PortchannelMemberList != nil {
                for _, v := range info.PortchannelMember.PortchannelMemberList {
                    if v.PortchannelMemberList == nil {
                        continue
                    }
                    c := helper.PortChannelMember{
                        Keys: []string{v.PortchannelName, v.Port},
                        Client: db,
                        Data: v.PortchannelMemberList,
                    }
                    c.SaveToDB(false)
                }
            }
        }
    }

    return nil
}
