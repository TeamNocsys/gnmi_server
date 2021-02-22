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

func AclTableHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysAcl{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.AclTable != nil {
            if info.AclTable.AclTableList != nil {
                for _, v := range info.AclTable.AclTableList {
                    if v.AclTableList == nil {
                        continue
                    }
                    c := helper.AclTable{
                        Key: v.TableName,
                        Client: db,
                        Data: v.AclTableList,
                    }
                    c.SaveToDB(true)
                }
            }
        }
    }

    return nil
}

func AclRuleHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysAcl{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.AclRule != nil {
            if info.AclRule.AclRuleList != nil {
                for _, v := range info.AclRule.AclRuleList {
                    if v.AclRuleList == nil {
                        continue
                    }
                    c := helper.AclRule{
                        Keys: []string{v.TableName, v.RuleName},
                        Client: db,
                        Data: v.AclRuleList,
                    }
                    c.SaveToDB(true)
                }
            }
        }
    }

    return nil
}