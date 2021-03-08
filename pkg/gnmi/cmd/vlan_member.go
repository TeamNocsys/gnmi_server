package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strings"
)

type VlanMemberAdapter struct {
    Adapter
    name string
    ifname string
}

func NewVlanMemberAdapter(name, ifname string, cli command.Client) *VlanMemberAdapter {
    return &VlanMemberAdapter{
        Adapter: Adapter{
            client: cli,
        },
        name:    name,
        ifname:  ifname,
    }
}

func (adpt *VlanMemberAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.NocsysVlan_VlanMember_VlanMemberList, error) {
    conn := adpt.client.Config()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.CONFIG_DB, []string{"VLAN_MEMBER", adpt.name, adpt.ifname}); err != nil {
        return nil, err
    } else {
        retval := &sonicpb.NocsysVlan_VlanMember_VlanMemberList{}
        for k, v := range data {
            switch k {
            case "tagging_mode":
                switch v {
                case "tagged":
                    retval.TaggingMode = sonicpb.NocsysVlan_VlanMember_VlanMemberList_TAGGINGMODE_tagged
                case "untagged":
                    retval.TaggingMode = sonicpb.NocsysVlan_VlanMember_VlanMemberList_TAGGINGMODE_untagged
                }
            }
        }
        return retval, nil
    }
}

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