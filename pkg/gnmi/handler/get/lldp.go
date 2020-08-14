package get

import (
    "context"
    "encoding/json"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
    "strings"

    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

const (
    LLDP_ENTRY_TABLE = "LLDP_ENTRY_TABLE*"
)

func LLDPHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    otherDB := db.State()

    dbConfig := swsssdk.Config()
    delimiter := dbConfig.GetDBSeparator(swsssdk.APPL_DB)

    lldpEntry, err := otherDB.GetAllByPattern(swsssdk.APPL_DB, LLDP_ENTRY_TABLE)
    if err != nil {
        logrus.Errorf("Get key %s from %s failed: %s", LLDP_ENTRY_TABLE,
            swsssdk.APPL_DB, err.Error())
        return nil, status.Errorf(codes.Internal, "Get table from database failed")
    }

    lldp := &sonicpb.SonicLldp_Lldp{}

    for key, value := range lldpEntry {
        words := strings.Split(key, delimiter)
        if len(words) != 2 {
            logrus.Errorf("Error key-%s when process lldp entry", key)
            continue
        }

        state := &sonicpb.SonicLldp_Lldp_DeviceList_State{}
        state.LocalPort = &ywrapper.StringValue{Value:words[1]}
        var name string

        for field, v := range value {
            switch field {
            case "lldp_rem_time_mark":
                age, err := parseUint(v, key, field)
                if err == nil {
                    state.Age = &ywrapper.UintValue{Value: age}
                }
            case "lldp_rem_chassis_id":
                state.ChassisId = &ywrapper.StringValue{Value: v}
            case "lldp_rem_chassis_id_subtype":
                subType, err := parseInt(v, key, field)
                if err == nil {
                    if subType < 0 || subType > 7 {
                        logrus.Errorf("unknown chassis id type: %s in %s %s",
                            v, swsssdk.APPL_DB, LLDP_ENTRY_TABLE)
                    } else {
                        state.ChassisIdType = sonicpb.SonicLldpChassisIdType(subType)
                    }
                }
            case "lldp_rem_index":
                index, err := parseUint(v, key, field)
                if err == nil {
                    state.Index = &ywrapper.UintValue{Value: index}
                }
            case "lldp_rem_port_desc":
                state.PortDescription = &ywrapper.StringValue{Value: v}
            case "lldp_rem_port_id":
                state.PortId = &ywrapper.StringValue{Value: v}
            case "lldp_rem_port_id_subtype":
                portId, err := parseInt(v, key, field)
                if err == nil {
                    state.PortIdType = sonicpb.SonicLldpPortIdType(portId)
                }
            case "lldp_rem_sys_desc":
                state.SystemDescription = &ywrapper.StringValue{Value: v}
            case "lldp_rem_sys_name":
                name = v
                state.SystemName = &ywrapper.StringValue{Value: v}
            }
        }

        deviceList := &sonicpb.SonicLldp_Lldp_DeviceList{
            Config: nil,
            State:  state,
        }
        deviceListKey := &sonicpb.SonicLldp_Lldp_DeviceListKey{
            DeviceName: name,
            DeviceList: deviceList,
        }
        lldp.DeviceList = append(lldp.DeviceList, deviceListKey)
    }

    bytes, err := json.Marshal(lldp)
    if err != nil {
        logrus.Errorf("marshal struct failed: %s", err.Error())
        return nil, status.Errorf(codes.Internal, "marshal json failed")
    }

    response, err := createResponse(ctx, r, bytes)
    if err != nil {
        logrus.Errorf("create response failed: %s", err.Error())
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    return response, nil
}

func parseUint(str string, key string, field string) (uint64, error) {
    value, err := strconv.ParseUint(str, 10, 64)
    if err != nil {
        logrus.Errorf("parse %s value-%s from %s %s to uint64 failed: %s",
            field, str, swsssdk.APPL_DB, key, err.Error())
        return 0, err
    }

    return value, nil
}

func parseInt(str string, key string, field string) (int64, error) {
    value, err := strconv.ParseInt(str, 10, 64)
    if err != nil {
        logrus.Errorf("parse %s value-%s from %s %s to int64 failed: %s",
            field, str, swsssdk.APPL_DB, key, err.Error())
        return 0, err
    }

    return value, nil
}
