package cmd

import (
    "fmt"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
    "github.com/sirupsen/logrus"
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

const req_LLDP_REM_PORT_ID      = 0x1
const req_LLDP_REM_PORT_ID_SUBT = 0x2
const req_LLDP_REM_CHAS_ID      = 0x4
const req_LLDP_REM_CHAS_ID_SUBT = 0x8
const req_LLDP_REM_BITS = (req_LLDP_REM_PORT_ID |
                           req_LLDP_REM_PORT_ID_SUBT |
                           req_LLDP_REM_CHAS_ID |
                           req_LLDP_REM_CHAS_ID_SUBT)

func (adpt *LldpAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.AcctonLldp_Lldp_LldpList, error) {
    retval := &sonicpb.AcctonLldp_Lldp_LldpList{}
    if dataType == gnmi.GetRequest_ALL || dataType == gnmi.GetRequest_STATE {
        conn := adpt.client.State()
        if conn == nil {
            return nil, swsssdk.ErrConnNotExist
        }

        if data, err := conn.GetAll(swsssdk.APPL_DB, []string{"LLDP_ENTRY_TABLE", adpt.ifname}); err != nil {
            return nil, err
        } else {
            retval.State = &sonicpb.AcctonLldp_Lldp_LldpList_State{}
            reqBits := 0
            for k, v := range data {
                // skip empty string in db
                if v == "" {
                    continue
                }
                switch k {
                case "lldp_rem_port_desc":
                    retval.State.LldpRemPortDesc = &ywrapper.StringValue{Value: v}
                case "lldp_rem_port_id":
                    // required
                    reqBits |= req_LLDP_REM_PORT_ID
                    retval.State.LldpRemPortId = &ywrapper.StringValue{Value: v}
                case "lldp_rem_port_id_subtype":
                    // required
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return nil, err
                    } else {
                        reqBits |= req_LLDP_REM_PORT_ID_SUBT
                        retval.State.LldpRemPortIdSubtype = sonicpb.AcctonLldp_Lldp_LldpList_State_LldpRemPortIdSubtype(i)
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
                    // required
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return nil, err
                    } else {
                        reqBits |= req_LLDP_REM_CHAS_ID_SUBT
                        retval.State.LldpRemChassisIdSubtype = sonicpb.AcctonLldp_Lldp_LldpList_State_LldpRemChassisIdSubtype(i)
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
                    // required
                    reqBits |= req_LLDP_REM_CHAS_ID
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

            // return error if required info is missing
            if reqBits != req_LLDP_REM_BITS {
                logrus.Debugf("required lldp info is missing for if/rid - %s/%v", adpt.ifname, retval.State.LldpRemIndex)
                return nil, fmt.Errorf("required info is missing")
            }
        }
    }
    return retval, nil
}
