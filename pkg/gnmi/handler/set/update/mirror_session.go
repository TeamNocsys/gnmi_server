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

func MirrorSessionHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.SonicMirrorSession{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.MirrorSession != nil {
            if info.MirrorSession.MirrorSessionList != nil {
                for _, v := range info.MirrorSession.MirrorSessionList {
                    if v.MirrorSessionList == nil {
                        continue
                    }
                    c := helper.MirrorSession{
                        Key: v.Name,
                        Client: db,
                        Data: v.MirrorSessionList,
                    }
                    c.SaveToDB()
                }
            }
        }
    }

    return nil
}