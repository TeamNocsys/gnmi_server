package cmd

import (
    "fmt"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
)

type LldpAdapter struct {
    Adapter
    ifname string
}

func NewLldpAdapter(ifname string, cli command.Client) *LldpAdapter {
    return &LldpAdapter{
        Adapter: Adapter{
            client: cli,
        },
        ifname:  ifname,
    }
}

func (adpt *LldpAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.NocsysLldp_Lldp_LldpList, error) {
    retval := &sonicpb.NocsysLldp_Lldp_LldpList{}
    if dataType == gnmi.GetRequest_ALL || dataType == gnmi.GetRequest_STATE {
        conn := adpt.client.State()
        if conn == nil {
            return nil, swsssdk.ErrConnNotExist
        }

        if data, err := conn.GetAll(swsssdk.APPL_DB, []string{"LLDP_ENTRY_TABLE", adpt.ifname}); err != nil {
            return nil, err
        } else {
            retval.State = &sonicpb.NocsysLldp_Lldp_LldpList_State{}
            for k, v := range data {
                switch k {
                case "lldp_rem_port_desc":
                    retval.State.LldpRemPortDesc = &ywrapper.StringValue{Value: v}
                case "lldp_rem_port_id":
                    retval.State.LldpRemPortId = &ywrapper.StringValue{Value: v}
                case "lldp_rem_port_id_subtype":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return nil, err
                    } else {
                        retval.State.LldpRemPortIdSubtype = sonicpb.NocsysLldp_Lldp_LldpList_State_LldpRemPortIdSubtype(i)
                    }
                case "lldp_rem_man_addr":
                    retval.State.LldpRemManAddr = &ywrapper.StringValue{Value: v}
                case "lldp_rem_time_mark":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return nil, err
                    } else {
                        retval.State.LldpRemTimeMark = &ywrapper.UintValue{Value: i}
                    }
                case "lldp_rem_chassis_id_subtype":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return nil, err
                    } else {
                        retval.State.LldpRemChassisIdSubtype = sonicpb.NocsysLldp_Lldp_LldpList_State_LldpRemChassisIdSubtype(i)
                    }
                case "lldp_rem_sys_cap_enabled":
                    var i uint64
                    if _, err := fmt.Sscanf(v, "%x 00", &i); err != nil {
                        return nil, err
                    } else {
                        retval.State.LldpRemSysCapEnabled = &ywrapper.UintValue{Value: i}
                    }
                case "lldp_rem_sys_name":
                    retval.State.LldpRemSysName = &ywrapper.StringValue{Value: v}
                case "lldp_rem_chassis_id":
                    retval.State.LldpRemChassisId = &ywrapper.StringValue{Value: v}
                case "lldp_rem_index":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return nil, err
                    } else {
                        retval.State.LldpRemIndex = &ywrapper.UintValue{Value: i}
                    }
                case "lldp_rem_sys_desc":
                    retval.State.LldpRemSysDesc = &ywrapper.StringValue{Value: v}
                case "lldp_rem_sys_cap_supported":
                    var i uint64
                    if _, err := fmt.Sscanf(v, "%x 00", &i); err != nil {
                        return nil, err
                    } else {
                        retval.State.LldpRemSysCapSupported = &ywrapper.UintValue{Value: i}
                    }
                }
            }
        }
    }
    return retval, nil
}