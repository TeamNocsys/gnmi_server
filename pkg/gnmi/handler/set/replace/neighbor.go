package replace

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/golang/protobuf/jsonpb"
    "github.com/golang/protobuf/proto"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func NeighborHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.AcctonNeighor{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("REPLACE|%s", s)
        if info.Neighor != nil {
            if info.Neighor.NeighorList != nil {
                for _, v := range info.Neighor.NeighorList {
                    if v.NeighorList == nil {
                        continue
                    }
                    c := helper.Neighbor{
                        Keys: []string{v.Name, v.IpPrefix},
                        Client: db,
                        Data: v.NeighorList,
                    }
                    c.SaveToDB(true)
                }
            }
        }
    }

    return nil
}