package replace

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

func VlanHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.SonicVlan{}
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
                    c := helper.Vlan{
                        Key: v.VlanName,
                        Client: db,
                        Data: v.VlanList,
                    }
                    c.SaveToDB(false)
                }
            }
        }
    }

    return nil
}

func VlanMemberHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.SonicVlan{}
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
                    c := helper.VlanMember{
                        Keys: []string{v.VlanName, v.Port},
                        Client: db,
                        Data: v.VlanMemberList,
                    }
                    c.SaveToDB(true)
                }
            }
        }
    }

    return nil
}