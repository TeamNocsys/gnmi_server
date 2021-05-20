package update

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

func VrfHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysVrf{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("UPDATE|%s", s)
        if info.Vrf != nil {
            if info.Vrf.VrfList != nil {
                for _, v := range info.Vrf.VrfList {
                    if v.VrfList == nil {
                        continue
                    }
                    c := cmd.NewVrfAdapter(v.VrfName, db)
                    if err := c.Config(v.VrfList, cmd.UPDATE); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}