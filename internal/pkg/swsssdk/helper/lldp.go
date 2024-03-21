package helper

import (
    "fmt"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
)

type Lldp struct {
    Key string
    Client command.Client
    Data *sonicpb.AcctonLldp_Lldp_LldpList
}

// 参考:
// https://github.com/lldpd/lldpd/blob/master/src/client/display.c
// https://github.com/Azure/sonic-dbsyncd/blob/master/src/lldp_syncd/daemon.py
// https://github.com/Azure/sonic-swss/blob/master/doc/swss-schema.md
func (c *Lldp) LoadFromDB(flags uint) error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.AcctonLldp_Lldp_LldpList{}
    }
    // 加载配置数据
    // 是否加载状态数据
    if (flags & DATA_TYPE_STATE) != 0 {
        if c.Data.State == nil {
            c.Data.State = &sonicpb.AcctonLldp_Lldp_LldpList_State{}
        }

        if data, err := conn.GetAll(swsssdk.APPL_DB, []string{"LLDP_ENTRY_TABLE", c.Key}); err != nil {
            return err
        } else {
            for k, v := range data {
                switch k {
                case "lldp_rem_port_desc":
                    c.Data.State.LldpRemPortDesc = &ywrapper.StringValue{Value: v}
                case "lldp_rem_port_id":
                    c.Data.State.LldpRemPortId = &ywrapper.StringValue{Value: v}
                case "lldp_rem_port_id_subtype":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return err
                    } else {
                        c.Data.State.LldpRemPortIdSubtype = sonicpb.AcctonLldp_Lldp_LldpList_State_LldpRemPortIdSubtype(i)
                    }
                case "lldp_rem_man_addr":
                    c.Data.State.LldpRemManAddr = &ywrapper.StringValue{Value: v}
                case "lldp_rem_time_mark":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return err
                    } else {
                        c.Data.State.LldpRemTimeMark = &ywrapper.UintValue{Value: i}
                    }
                case "lldp_rem_chassis_id_subtype":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return err
                    } else {
                        c.Data.State.LldpRemChassisIdSubtype = sonicpb.AcctonLldp_Lldp_LldpList_State_LldpRemChassisIdSubtype(i)
                    }
                case "lldp_rem_sys_cap_enabled":
                    var i uint64
                    if _, err := fmt.Sscanf(v, "%x 00", &i); err != nil {
                        return err
                    } else {
                        c.Data.State.LldpRemSysCapEnabled = &ywrapper.UintValue{Value: i}
                    }
                case "lldp_rem_sys_name":
                    c.Data.State.LldpRemSysName = &ywrapper.StringValue{Value: v}
                case "lldp_rem_chassis_id":
                    c.Data.State.LldpRemChassisId = &ywrapper.StringValue{Value: v}
                case "lldp_rem_index":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return err
                    } else {
                        c.Data.State.LldpRemIndex = &ywrapper.UintValue{Value: i}
                    }
                case "lldp_rem_sys_desc":
                    c.Data.State.LldpRemSysDesc = &ywrapper.StringValue{Value: v}
                case "lldp_rem_sys_cap_supported":
                    var i uint64
                    if _, err := fmt.Sscanf(v, "%x 00", &i); err != nil {
                        return err
                    } else {
                        c.Data.State.LldpRemSysCapSupported = &ywrapper.UintValue{Value: i}
                    }
                }
            }
        }
    }

    return nil
}