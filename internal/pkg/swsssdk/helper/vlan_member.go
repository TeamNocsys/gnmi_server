package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
)

type VlanMember struct {
    Keys []string
    Client command.Client
    Data *sonicpb.SonicVlan_VlanMember_VlanMemberList
}

func (c *VlanMember) LoadFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.SonicVlan_VlanMember_VlanMemberList{}
    }
    if data, err := conn.GetAll(swsssdk.CONFIG_DB, append([]string{"VLAN_MEMBER"}, c.Keys...)); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "tagging_mode":
                switch v {
                case "tagged":
                    c.Data.TaggingMode = sonicpb.SonicVlan_VlanMember_VlanMemberList_TAGGINGMODE_tagged
                case "untagged":
                    c.Data.TaggingMode = sonicpb.SonicVlan_VlanMember_VlanMemberList_TAGGINGMODE_untagged
                }
            }
        }
    }
    return nil
}

func (c *VlanMember) SaveToDB(replace bool) error {
    e := make(map[string]interface{})
    if c.Data.TaggingMode != sonicpb.SonicVlan_VlanMember_VlanMemberList_TAGGINGMODE_UNSET {
        if c.Data.TaggingMode == sonicpb.SonicVlan_VlanMember_VlanMemberList_TAGGINGMODE_tagged {
            e["tagging_mode"] = "tagged"
        } else if c.Data.TaggingMode == sonicpb.SonicVlan_VlanMember_VlanMemberList_TAGGINGMODE_untagged {
            e["tagging_mode"] = "untagged"
        }
    }

    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if replace {
        if _, err := conn.SetEntry("VLAN_MEMBER", c.Keys, e); err != nil {
            return err
        }
    } else {
        if _, err := conn.ModEntry("VLAN_MEMBER", c.Keys, e); err != nil {
            return err
        }
    }
    return nil
}

func (c *VlanMember) RemoveFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }
    if _, err := conn.DeleteAllByPattern(swsssdk.CONFIG_DB, append([]string{"VLAN_MEMBER"}, c.Keys...)); err != nil {
        return err
    }
    return nil
}