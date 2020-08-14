package get

import (
    "context"
    "errors"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "gnmi_server/pkg/gnmi/handler"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "strconv"
    "strings"
)

func PortInfoHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.State()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.APPL_DB)
    states, err := conn.GetAllByPattern(swsssdk.APPL_DB, []string{helper.PORT_STATUS_TABLE_NAME, helper.PORT_NAME_PREFIX+"*"})
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    sp := &sonicpb.SonicPort{
        Port: &sonicpb.SonicPort_Port{},
    }
    for key, value := range states {
        names := strings.SplitN(key, s, 2)
        if state, err := getPortState(value); err != nil {
            return nil, err
        } else {
            sp.Port.PortStateList = append(sp.Port.PortStateList, &sonicpb.SonicPort_Port_PortStateListKey{
                PortName:      names[len(names) - 1],
                PortStateList: state,
            })
        }
    }
    response, err := handler.CreateResponse(ctx, r, sp)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getPortState(state map[string]string) (*sonicpb.SonicPort_Port_PortStateList, error) {
    r := &sonicpb.SonicPort_Port_PortStateList{}

    if v, ok := state[helper.PORT_STATUS_ALIAS_FIELD]; ok {
        r.Alias = &ywrapper.StringValue{Value: v}
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_ALIAS_FIELD + " field")
    }

    if v, ok := state[helper.PORT_STATUS_SPEED_FIELD]; ok {
        if speed, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Speed = &ywrapper.UintValue{Value: speed}
        }
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_SPEED_FIELD + " field")
    }

    if v, ok := state[helper.PORT_STATUS_ADMIN_STATUS_FIELD]; ok {
        if strings.ToUpper(v) == helper.PORT_STATUS_VALUE_UP {
            r.AdminStatus = sonicpb.SonicPortAdminStatus_SONICPORTADMINSTATUS_up
        } else {
            r.AdminStatus = sonicpb.SonicPortAdminStatus_SONICPORTADMINSTATUS_down
        }
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_ADMIN_STATUS_FIELD + " field")
    }

    if v, ok := state[helper.PORT_STATUS_OPER_STATUS_FIELD]; ok {
        if strings.ToUpper(v) == helper.PORT_STATUS_VALUE_UP {
            r.OperStatus = sonicpb.SonicPortOperStatus_SONICPORTOPERSTATUS_up
        } else {
            r.OperStatus = sonicpb.SonicPortOperStatus_SONICPORTOPERSTATUS_down
        }
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_OPER_STATUS_FIELD + " field")
    }

    if v, ok := state[helper.PORT_STATUS_MTU_FIELD]; ok {
        if mtu, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Mtu = &ywrapper.UintValue{Value: mtu}
        }
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_MTU_FIELD + " field")
    }

    if v, ok := state[helper.PORT_STATUS_INDEX_FIELD]; ok {
        if index, err := strconv.ParseUint(v, 10, 64); err != nil {
            return nil, err
        } else {
            r.Index = &ywrapper.UintValue{Value: index}
        }
    } else {
        return nil, errors.New("missing " + helper.PORT_STATUS_INDEX_FIELD + " field")
    }

    return r, nil
}
