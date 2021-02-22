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
    "strings"
)

func FdbHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.State()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    kvs := handler.FetchPathKey(r)
    spec := []string{}
    if v, ok := kvs["vlan-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["mac-address"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    sf := &sonicpb.NocsysFdb{
        Fdb: &sonicpb.NocsysFdb_Fdb{},
    }
    if hkeys, err := conn.GetKeys(swsssdk.STATE_DB, []string{"FDB_TABLE", strings.Join(spec, ":")}); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(swsssdk.STATE_DB, hkey)
            c := helper.Fdb{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sf.Fdb.FdbList = append(sf.Fdb.FdbList,
                &sonicpb.NocsysFdb_Fdb_FdbListKey{
                    FdbName: keys[0],
                    FdbList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sf)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}