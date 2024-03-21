package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
)

type PortStatisticsAdapter struct {
    Adapter
    name string
}

func NewPortStatisticsAdapter(name string, cli command.Client) *PortStatisticsAdapter {
    return &PortStatisticsAdapter{
        Adapter: Adapter{
            client: cli,
        },
        name:  name,
    }
}

func (adpt *PortStatisticsAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.AcctonPort_Port_PortStatisticsList, error) {
    conn := adpt.client.State()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.COUNTERS_DB, []string{"COUNTERS", adpt.name}); err != nil {
        return nil, err
    } else {
        retval := &sonicpb.AcctonPort_Port_PortStatisticsList{}
        for k, v := range data {
            switch k {
            case "SAI_PORT_STAT_IF_IN_OCTETS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.InOctets = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_UCAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.InUnicastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_MULTICAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.InMulticastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_BROADCAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.InBroadcastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_DISCARDS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.InDiscards = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_ERRORS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.InErrors = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_UNKNOWN_PROTOS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.InUnknownProtos = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_OCTETS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.OutOctets = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_UCAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.OutUnicastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_MULTICAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.OutMulticastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_BROADCAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.OutBroadcastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_DISCARDS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.OutDiscards = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_ERRORS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.OutErrors = &ywrapper.UintValue{Value: pkts}
                }
            }
        }
        return retval, nil
    }
}