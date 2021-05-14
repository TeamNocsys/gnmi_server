// +build broadcom

package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "gnmi_server/internal/pkg/swsssdk"
)

func (adpt *VrfAdapter) Config(data *sonicpb.NocsysVrf_Vrf_VrfList, oper OperType) error {
    var cmdstr string
    if oper == ADD {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("VRF", adpt.name); err != nil {
            return err
        } else if ok {
            return nil
        }

        cmdstr = "config vrf add " + adpt.name
    } else if oper == DEL {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("VRF", adpt.name); err != nil {
            return err
        } else if !ok {
            return nil
        }

        cmdstr = "config vrf del " + adpt.name
    }
    return adpt.exec(cmdstr)
}