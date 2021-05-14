// +build broadcom

package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "gnmi_server/internal/pkg/swsssdk"
)

func (adpt *LagAdapter) Config(data *sonicpb.NocsysPortchannel_Portchannel_PortchannelList, oper OperType) error {
    cmdstr := "config portchannel"
    if oper == ADD {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("PORTCHANNEL", adpt.name); err != nil {
            return err
        } else if ok {
            return nil
        }

        cmdstr += " add " + adpt.name
    } else if oper == DEL {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("PORTCHANNEL", adpt.name); err != nil {
            return err
        } else if !ok {
            return nil
        }

        cmdstr += " del " + adpt.name
    } else {
        return ErrInvalidOperType
    }

    return adpt.exec(cmdstr)
}