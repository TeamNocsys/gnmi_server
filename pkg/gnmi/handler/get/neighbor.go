package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "gnmi_server/pkg/gnmi/handler"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func NeighborHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.State()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    kvs := handler.FetchPathKey(r)
    spec := []string{}
    if v, ok := kvs["name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "Vlan*")
    }
    if v, ok := kvs["ip-prefix"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    si := &sonicpb.SonicNeighor{
        Neighor: &sonicpb.SonicNeighor_Neighor{},
    }
    if hkeys, err := conn.GetKeys(swsssdk.APPL_DB, append([]string{"NEIGH_TABLE"}, spec...)); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(swsssdk.APPL_DB, hkey)
            c := helper.Neighbor{
                Keys: keys,
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            si.Neighor.NeighorList = append(si.Neighor.NeighorList,
                &sonicpb.SonicNeighor_Neighor_NeighorListKey{
                    Name: keys[0],
                    IpPrefix: keys[1],
                    NeighorList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, si)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}