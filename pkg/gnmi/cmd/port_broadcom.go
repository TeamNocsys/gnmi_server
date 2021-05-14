// +build broadcom

package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "strconv"
)

func (adpt *PortAdapter) Config(data *sonicpb.NocsysPort_Port_PortList, oper OperType) error {
    if oper == ADD || oper == UPDATE {
        if data.Mtu != nil {
            cmdstr := "config interface mtu " + adpt.ifname + " " + strconv.FormatUint(data.Mtu.Value, 10)
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.AdminStatus != sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_UNSET {
            var cmdstr string
            if data.AdminStatus == sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_up {
                cmdstr = "config interface startup " + adpt.ifname
            } else if data.AdminStatus == sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_down {
                cmdstr = "config interface shutdown " + adpt.ifname
            }
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.Speed != nil {
            cmdstr := "config interface speed " + adpt.ifname + " " + strconv.FormatUint(data.Mtu.Value, 10)
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        return nil
    }

    return ErrInvalidOperType
}