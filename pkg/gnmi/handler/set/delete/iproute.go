package delete

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/pkg/gnmi/cmd"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func IpRouteHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    vrf := "*"
    dest := "*"
    if v, ok := kvs["vrf-name"]; ok {
        vrf = v
    }
    if v, ok := kvs["ip-prefix"]; ok {
        dest = v
    }

    conn := db.State()
    if conn == nil {
        return status.Error(codes.Internal, "")
    }
    if v, err := conn.GetAll(swsssdk.APPL_DB, append([]string{"ROUTE_TABLE"}, vrf, dest)); err != nil {
        return err
    } else {
        data := &sonicpb.NocsysRoute_Route_RouteList{}
        for k, v := range v {
            switch k {
            case "nexthop":
                data.Nexthop = &ywrapper.StringValue{Value: v}
            case "ifname":
                data.Ifname = &ywrapper.StringValue{Value: v}
            }
        }
        c := cmd.NewVrfRouteAdapter(vrf, dest, db)
        return c.Config(data, cmd.DEL)
    }
}