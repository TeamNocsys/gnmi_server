package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk/helper"
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
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := helper.Vrf{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sv.Vrf.VrfList = append(sv.Vrf.VrfList,
                &sonicpb.NocsysVrf_Vrf_VrfListKey{
                    VrfName: keys[0],
                    VrfList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sv)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}