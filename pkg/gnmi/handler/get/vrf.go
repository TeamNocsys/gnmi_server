package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
    "gnmi_server/pkg/gnmi/handler"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)


func VrfHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["vrf-name"]; ok {
        spec = v
    }

    sv := &sonicpb.NocsysVrf{
        Vrf: &sonicpb.NocsysVrf_Vrf{},
    }
    if hkeys, err := conn.GetKeys("VRF", spec); err != nil {
        return nil, err
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := cmd.NewVrfAdapter(keys[0], db)
            if data, err := c.Show(r.Type); err != nil {
                return nil, err
            } else {
                sv.Vrf.VrfList = append(sv.Vrf.VrfList,
                    &sonicpb.NocsysVrf_Vrf_VrfListKey{
                        VrfName: keys[0],
                        VrfList: data,
                    })
            }
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sv)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}