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

func FdbHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysFdb{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.Fdb != nil {
            if info.Fdb.FdbList != nil {
                for _, v := range info.Fdb.FdbList {
                    if v.FdbList == nil {
                        continue
                    }
                    c := helper.Fdb{
                        Key: v.FdbName,
                        Client: db,
                        Data: v.FdbList,
                    }
                    c.SaveToDB(true)
                }
            }
        }
    }

    return nil
}