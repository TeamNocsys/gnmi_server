// +build ec

package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "gnmi_server/internal/pkg/swsssdk"
)

func (adpt *LagMemberAdapter) Config(data *sonicpb.NocsysPortchannel_PortchannelMember_PortchannelMemberList, oper OperType) error {
    cmdstr := "config portchannel member"
    if oper == ADD || oper == UPDATE {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("PORTCHANNEL_MEMBER", []string{adpt.name, adpt.ifname}); err != nil {
            return err
        } else if ok {
            return nil
        }

        cmdstr += " add " + adpt.name + " " + adpt.ifname
    } else if oper == DEL {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("PORTCHANNEL_MEMBER", []string{adpt.name, adpt.ifname}); err != nil {
            return err
        } else if !ok {
            return nil
        }

        cmdstr += " del " + adpt.name + " " + adpt.ifname
    } else {
        return ErrInvalidOperType
    }

    return adpt.exec(cmdstr)
}