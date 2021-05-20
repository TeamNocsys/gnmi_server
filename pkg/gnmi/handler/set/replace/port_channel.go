package replace

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

func PortChannelHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysPortchannel{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("SET|%s", s)
        if info.Portchannel != nil {
            if info.Portchannel.PortchannelList != nil {
                for _, v := range info.Portchannel.PortchannelList {
                    if v.PortchannelList == nil {
                        continue
                    }
                    c := cmd.NewLagAdapter(v.PortchannelName, db)
                    if err := c.Config(v.PortchannelList, cmd.ADD); err != nil {
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
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("REPLACE|%s", s)
        if info.PortchannelMember != nil {
            if info.PortchannelMember.PortchannelMemberList != nil {
                for _, v := range info.PortchannelMember.PortchannelMemberList {
                    if v.PortchannelMemberList == nil {
                        continue
                    }
                    c := cmd.NewLagMemberAdapter(v.PortchannelName, v.Port, db)
                    if err := c.Config(v.PortchannelMemberList, cmd.ADD); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}
