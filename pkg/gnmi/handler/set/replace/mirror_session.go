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

func MirrorSessionHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.AcctonMirrorSession{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("REPLACE|%s", s)
        if info.MirrorSession != nil {
            if info.MirrorSession.MirrorSessionList != nil {
                for _, v := range info.MirrorSession.MirrorSessionList {
                    if v.MirrorSessionList == nil {
                        continue
                    }
                    c := cmd.NewMirrorAdapter(v.Name, db)
                    if err := c.Config(v.MirrorSessionList, cmd.ADD); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}