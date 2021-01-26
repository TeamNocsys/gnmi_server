package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
)

type PortStatistics struct {
    Key string
    Client command.Client
    Data *sonicpb.SonicPort_Port_PortStatisticsList
}

// 参考
// https://github.com/Azure/sonic-swss/blob/master/orchagent/portsorch.cpp
func (c *PortStatistics) LoadFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if c.Data == nil {
        c.Data = &sonicpb.SonicPort_Port_PortStatisticsList{}
    }
    if data, err := conn.GetAll(swsssdk.COUNTERS_DB, []string{"COUNTERS", c.Key}); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "SAI_PORT_STAT_IF_IN_OCTETS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.InOctets = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_UCAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.InUnicastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_MULTICAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.InMulticastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_BROADCAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.InBroadcastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_DISCARDS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.InDiscards = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_ERRORS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.InErrors = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_IN_UNKNOWN_PROTOS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.InUnknownProtos = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_OCTETS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.OutOctets = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_UCAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.OutUnicastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_MULTICAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.OutMulticastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_BROADCAST_PKTS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.OutBroadcastPkts = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_DISCARDS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.OutDiscards = &ywrapper.UintValue{Value: pkts}
                }
            case "SAI_PORT_STAT_IF_OUT_ERRORS":
                if pkts, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.OutErrors = &ywrapper.UintValue{Value: pkts}
                }
            }
        }
    }

    return nil
}
