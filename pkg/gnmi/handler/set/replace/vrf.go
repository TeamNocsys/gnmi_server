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

func VrfHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysVrf{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.Vrf != nil {
            if info.Vrf.VrfList != nil {
                for _, v := range info.Vrf.VrfList {
                    if v.VrfList == nil {
                        continue
                    }
                    c := helper.Vrf{
                        Key: v.VrfName,
                        Client: db,
                        Data: v.VrfList,
                    }
                    c.SaveToDB()
                }
            }
        }
    }

    return nil
}