package update

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/golang/protobuf/proto"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func PortChannelHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysPortchannel{}
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
                    c := cmd.NewLagAdapter(v.PortchannelName, db)
                    if err := c.Config(v.PortchannelList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}

func PortChannelMemberHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysPortchannel{}
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
                    c := cmd.NewLagMemberAdapter(v.PortchannelName, v.Port, db)
                    if err := c.Config(v.PortchannelMemberList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}
