package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
)

type PortChannelMember struct {
    Keys []string
    Client command.Client
    Data *sonicpb.AcctonPortchannel_PortchannelMember_PortchannelMemberList
}

func (c *PortChannelMember) LoadFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.AcctonPortchannel_PortchannelMember_PortchannelMemberList{}
    }
    if ok, err := conn.HasEntry("PORTCHANNEL_MEMBER", c.Keys); err != nil {
        return err
    } else if !ok {
        return swsssdk.ErrConnNotExist
    }
    return nil
}

func (c *PortChannelMember) SaveToDB(replace bool) error {
    e := make(map[string]interface{})

    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if replace {
        if _, err := conn.SetEntry("PORTCHANNEL_MEMBER", c.Keys, e); err != nil {
            return err
        }
    } else {
        if _, err := conn.ModEntry("PORTCHANNEL_MEMBER", c.Keys, e); err != nil {
            return err
        }
    }
    return nil
}

func (c *PortChannelMember) RemoveFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }
    if _, err := conn.DeleteAllByPattern(swsssdk.CONFIG_DB, append([]string{"PORTCHANNEL_MEMBER"}, c.Keys...)); err != nil {
        return err
    }
    return nil
}