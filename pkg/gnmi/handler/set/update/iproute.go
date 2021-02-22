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

func IpRouteHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysRoute{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.Route != nil {
            if info.Route.RouteList != nil {
                for _, v := range info.Route.RouteList {
                    if v.RouteList == nil {
                        continue
                    }
                    c := helper.IpRoute{
                        Keys: []string{v.VrfName, v.IpPrefix},
                        Client: db,
                        Data: v.RouteList,
                    }
                    c.SaveToDB()
                }
            }
        }
    }

    return nil
}
