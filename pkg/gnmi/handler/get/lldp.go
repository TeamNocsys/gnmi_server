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

func LLDPHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.State()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定LLDP或全部LLDP
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["port-name"]; ok {
        spec = v
    }

    sl := &sonicpb.NocsysLldp{
        Lldp: &sonicpb.NocsysLldp_Lldp{},
    }
    if hkeys, err := conn.GetKeys(swsssdk.APPL_DB, []string{"LLDP_ENTRY_TABLE", spec}); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(swsssdk.APPL_DB, hkey)
            c := helper.Lldp{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(helper.DATA_TYPE_ALL); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sl.Lldp.LldpList = append(sl.Lldp.LldpList,
                &sonicpb.NocsysLldp_Lldp_LldpListKey{
                    PortName: keys[0],
                    LldpList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sl)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}
