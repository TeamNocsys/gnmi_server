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

func VlanHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysVlan{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.Vlan != nil {
            if info.Vlan.VlanList != nil {
                for _, v := range info.Vlan.VlanList {
                    if v.VlanList == nil {
                        continue
                    }
                    c := cmd.NewVlanAdapter(v.VlanName, db)
                    if err := c.Config(v.VlanList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}

func VlanMemberHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysVlan{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.VlanMember != nil {
            if info.VlanMember.VlanMemberList != nil {
                for _, v := range info.VlanMember.VlanMemberList {
                    if v.VlanMemberList == nil {
                        continue
                    }
                    c := cmd.NewVlanMemberAdapter(v.VlanName, v.Port, db)
                    if err := c.Config(v.VlanMemberList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}