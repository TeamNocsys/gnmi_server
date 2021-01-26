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

func PortHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.SonicPort{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.Port != nil {
            if info.Port.PortList != nil {
                for _, v := range info.Port.PortList {
                    if v.PortList == nil {
                        continue
                    }
                    c := helper.Port{
                        Key: v.PortName,
                        Client: db,
                        Data: v.PortList,
                    }
                    c.SaveToDB(false)
                }
            }
        }
    }

    return nil
}