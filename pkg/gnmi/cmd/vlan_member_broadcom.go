// +build broadcom

package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "gnmi_server/internal/pkg/swsssdk"
    "strings"
)

func (adpt *VlanMemberAdapter) Config(data *sonicpb.NocsysVlan_VlanMember_VlanMemberList, oper OperType) error {
    cmdstr := "config vlan member"
    if oper == ADD || oper == UPDATE {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("VLAN_MEMBER", []string{adpt.name, adpt.ifname}); err != nil {
            return err
        } else if ok {
            return nil
        }

        cmdstr += " add " + strings.TrimLeft(adpt.name, "Vlan") + " " +  adpt.ifname
        if data.TaggingMode == sonicpb.NocsysVlan_VlanMember_VlanMemberList_TAGGINGMODE_untagged {
            cmdstr += " -u"
        }
    } else if oper == DEL {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("VLAN_MEMBER", []string{adpt.name, adpt.ifname}); err != nil {
            return err
        } else if !ok {
            return nil
        }

        cmdstr += " del " + strings.TrimLeft(adpt.name, "Vlan") + " " +  adpt.ifname
    } else {
        return ErrInvalidOperType
    }

    return adpt.exec(cmdstr)
}