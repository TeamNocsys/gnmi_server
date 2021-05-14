// +build ec

package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
    "strings"
)

func (adpt *VlanAdapter) Config(data *sonicpb.NocsysVlan_Vlan_VlanList, oper OperType) error {
    if oper == ADD || oper == UPDATE {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("VLAN", adpt.name); err != nil {
            return err
        } else if !ok {
            cmdstr := "config vlan add " + strings.TrimLeft(adpt.name, "Vlan")
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.Description != nil {
            cmdstr := "config interface description " + adpt.name + " " +  data.Description.Value
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.Mtu != nil {
            cmdstr := "config interface mtu " + adpt.name + " " +  strconv.FormatUint(data.Mtu.Value, 10)
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.AdminStatus != sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_UNSET {
            var cmdstr string
            if data.AdminStatus == sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_up {
                cmdstr = "config interface startup " + adpt.name
            } else if data.AdminStatus == sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_down {
                cmdstr = "config interface shutdown " + adpt.name
            }
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.DhcpServers != nil {
            // TODO: cmdstr = "config interface ip dhcp-relay add " + vid
        }
    } else if oper == DEL {
        cmdstr := "config vlan del " + strings.TrimLeft(adpt.name, "Vlan")
        if err := adpt.exec(cmdstr); err != nil {
            return err
        }
    }
    return nil
}