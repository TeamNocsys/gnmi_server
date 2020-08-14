package get

import (
    "context"
    "errors"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "gnmi_server/pkg/gnmi/handler"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "strconv"
    "strings"
)

func PortStateHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.State()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定端口或全部端口
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["port-name"]; ok {
        spec = v
    }

    // 获取统计表存储的端口名称
    statNames, err := conn.GetAll(swsssdk.COUNTERS_DB, helper.COUNTERS_PORT_NAME_MAP)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.APPL_DB)
    states, err := conn.GetAllByPattern(swsssdk.APPL_DB, []string{helper.PORT_STATUS, helper.PORT_NAME_PREFIX + spec})
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    sp := &sonicpb.SonicPort{
        Port: &sonicpb.SonicPort_Port{},
    }
    for hash, value := range states {
        keys := strings.SplitN(hash, s, 2)
        name := keys[len(keys)-1]
        statName, ok := statNames[name]
        if !ok {
            logrus.Warningf("Missing ")
            continue
        }
        counters, err := conn.GetAll(swsssdk.COUNTERS_DB, helper.COUNTER_TABLE_PREFIX+statName)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        state, err := getPortState(value, counters)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        sp.Port.PortStateList = append(sp.Port.PortStateList, &sonicpb.SonicPort_Port_PortStateListKey{
            PortName:      name,
            PortStateList: state,
        })
    }

    response, err := handler.CreateResponse(ctx, r, sp)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getPortState(state map[string]string, counters map[string]string) (*sonicpb.SonicPort_Port_PortStateList, error) {
    r := &sonicpb.SonicPort_Port_PortStateList{}

    if v, ok := state[helper.PORT_STATUS_ALIAS]; ok {
        r.Alias = &ywrapper.StringValue{Value: v}
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_ALIAS + " field")
    }

    if v, ok := state[helper.PORT_STATUS_SPEED]; ok {
        if speed, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Speed = &ywrapper.UintValue{Value: speed}
        }
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_SPEED + " field")
    }

    if v, ok := state[helper.PORT_STATUS_ADMIN_STATUS]; ok {
        if strings.ToUpper(v) == helper.PORT_STATUS_VALUE_UP {
            r.AdminStatus = sonicpb.SonicPortAdminStatus_SONICPORTADMINSTATUS_up
        } else {
            r.AdminStatus = sonicpb.SonicPortAdminStatus_SONICPORTADMINSTATUS_down
        }
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_ADMIN_STATUS + " field")
    }

    if v, ok := state[helper.PORT_STATUS_OPER_STATUS]; ok {
        if strings.ToUpper(v) == helper.PORT_STATUS_VALUE_UP {
            r.OperStatus = sonicpb.SonicPortOperStatus_SONICPORTOPERSTATUS_up
        } else {
            r.OperStatus = sonicpb.SonicPortOperStatus_SONICPORTOPERSTATUS_down
        }
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_OPER_STATUS + " field")
    }

    if v, ok := state[helper.PORT_STATUS_MTU]; ok {
        if mtu, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Mtu = &ywrapper.UintValue{Value: mtu}
        }
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_MTU + " field")
    }

    if v, ok := state[helper.PORT_STATUS_INDEX]; ok {
        if index, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Index = &ywrapper.UintValue{Value: index}
        }
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_INDEX + " field")
    }

    var err error
    r.Counters, err = getPortStateCounters(counters)
    if err != nil {
        return nil, err
    }

    return r, nil
}

func getPortStateCounters(counters map[string]string) (*sonicpb.SonicPort_Port_PortStateList_Counters, error) {
    r := &sonicpb.SonicPort_Port_PortStateList_Counters{}

    if v, ok := counters[helper.COUNTERS_PORT_IN_UCAST_PKTS]; ok {
        if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.InUnicastPkts = &ywrapper.UintValue{Value: pkts}
        }
    } else {
        return nil, errors.New("missing " + helper.COUNTERS_PORT_IN_UCAST_PKTS + " field")
    }

    if v, ok := counters[helper.COUNTERS_PORT_IN_ERRORS]; ok {
        if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.InErrors = &ywrapper.UintValue{Value: pkts}
        }
    } else {
        return nil, errors.New("missing " + helper.COUNTERS_PORT_IN_ERRORS + " field")
    }

    if v, ok := counters[helper.COUNTERS_PORT_IN_DISCARDS]; ok {
        if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.InDiscards = &ywrapper.UintValue{Value: pkts}
        }
    } else {
        return nil, errors.New("missing " + helper.COUNTERS_PORT_IN_DISCARDS + " field")
    }

    if v, ok := counters[helper.COUNTERS_PORT_OUT_UCAST_PKTS]; ok {
        if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.OutUnicastPkts = &ywrapper.UintValue{Value: pkts}
        }
    } else {
        return nil, errors.New("missing " + helper.COUNTERS_PORT_OUT_UCAST_PKTS + " field")
    }

    if v, ok := counters[helper.COUNTERS_PORT_OUT_ERRORS]; ok {
        if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.OutErrors = &ywrapper.UintValue{Value: pkts}
        }
    } else {
        return nil, errors.New("missing " + helper.COUNTERS_PORT_OUT_ERRORS + " field")
    }

    if v, ok := counters[helper.COUNTERS_PORT_OUT_DISCARDS]; ok {
        if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.OutDiscards = &ywrapper.UintValue{Value: pkts}
        }
    } else {
        return nil, errors.New("missing " + helper.COUNTERS_PORT_OUT_DISCARDS + " field")
    }

    if v, ok := counters[helper.COUNTERS_PORT_IN_OCTETS]; ok {
        if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.InOctets = &ywrapper.UintValue{Value: pkts}
        }
    } else {
        return nil, errors.New("missing " + helper.COUNTERS_PORT_IN_OCTETS + " field")
    }

    if v, ok := counters[helper.COUNTERS_PORT_OUT_OCTETS]; ok {
        if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.OutOctets = &ywrapper.UintValue{Value: pkts}
        }
    } else {
        return nil, errors.New("missing " + helper.COUNTERS_PORT_OUT_OCTETS + " field")
    }

    return r, nil
}