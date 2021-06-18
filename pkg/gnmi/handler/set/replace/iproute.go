package replace

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/golang/protobuf/jsonpb"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"

     handler_utils "gnmi_server/pkg/gnmi/handler/utils"
)

func IpRouteHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysRoute{}

    if err := handler_utils.UnmarshalGpbValue(value, info); err != nil {
	return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("REPLACE|%s", s)
        if info.Route != nil {
            if info.Route.RouteList != nil {
                for _, v := range info.Route.RouteList {
                    if v.RouteList == nil {
                        continue
                    }
                    c := cmd.NewVrfRouteAdapter(v.VrfName, v.IpPrefix, db)
                    if err := c.Config(v.RouteList, cmd.ADD); err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}
