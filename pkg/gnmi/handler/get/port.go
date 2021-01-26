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

func PortHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定Port或全部Port
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["port-name"]; ok {
        spec = v
    }

    sp := &sonicpb.SonicPort{
        Port: &sonicpb.SonicPort_Port{},
    }
    if hkeys, err := conn.GetKeys("PORT", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := helper.Port{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(helper.DATA_TYPE_ALL); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sp.Port.PortList = append(sp.Port.PortList,
                &sonicpb.SonicPort_Port_PortListKey{
                    PortName: keys[0],
                    PortList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sp)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func PortStatisticsHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.State()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定端口或全部端口
    kvs := handler.FetchPathKey(r)
    statNames, err := conn.GetAll(swsssdk.COUNTERS_DB, "COUNTERS_PORT_NAME_MAP")
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    if v, ok := kvs["port-name"]; ok {
        if objId, ok := statNames[v]; ok {
            statNames = map[string]string{v: objId}
        } else {
            statNames = map[string]string{}
        }
    }
    sp := &sonicpb.SonicPort{
        Port: &sonicpb.SonicPort_Port{},
    }
    for name, objId := range statNames {
        c := helper.PortStatistics{
            Key:    objId,
            Client: db,
            Data:   nil,
        }
        if err := c.LoadFromDB(); err != nil {
            return nil, status.Errorf(codes.Internal, err.Error())
        }
        sp.Port.PortStatisticsList = append(sp.Port.PortStatisticsList,
            &sonicpb.SonicPort_Port_PortStatisticsListKey{
                PortName:           name,
                PortStatisticsList: c.Data,
            })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sp)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}