package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
)

type PortChannel struct {
    Key string
    Client command.Client
    Data *sonicpb.AcctonPortchannel_Portchannel_PortchannelList
}

func (c *PortChannel) LoadFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.AcctonPortchannel_Portchannel_PortchannelList{}
    }
    if data, err := conn.GetAll(swsssdk.CONFIG_DB, []string{"PORTCHANNEL", c.Key}); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "mtu":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.Mtu = &ywrapper.UintValue{Value: i}
                }
            case "admin_status":
                switch v {
                case "up":
                    c.Data.AdminStatus = sonicpb.AcctonTypesAdminStatus_ACCTONTYPESADMINSTATUS_up
                case "down":
                    c.Data.AdminStatus = sonicpb.AcctonTypesAdminStatus_ACCTONTYPESADMINSTATUS_down
                }
            case "min_links":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.MinLinks = &ywrapper.UintValue{Value: i}
                }
            }
        }
    }

    return nil
}

func (c *PortChannel) SaveToDB(replace bool) error {
    e := make(map[string]interface{})
    // 需要更新状态字段
    if c.Data.AdminStatus != sonicpb.AcctonTypesAdminStatus_ACCTONTYPESADMINSTATUS_UNSET {
        if c.Data.AdminStatus == sonicpb.AcctonTypesAdminStatus_ACCTONTYPESADMINSTATUS_up {
            e["admin_status"] = "up"
        } else if c.Data.AdminStatus == sonicpb.AcctonTypesAdminStatus_ACCTONTYPESADMINSTATUS_down {
            e["admin_status"] = "down"
        }
    }
    // 需要更新MTU字段
    if c.Data.Mtu != nil {
        e["mtu"] = c.Data.Mtu.Value
    }
    // 需要更新链接数字段
    if c.Data.MinLinks != nil {
        e["min_links"] = c.Data.MinLinks.Value
    }

    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if replace {
        if _, err := conn.SetEntry("PORTCHANNEL", c.Key, e); err != nil {
            return err
        }
    } else {
        if _, err := conn.ModEntry("PORTCHANNEL", c.Key, e); err != nil {
            return err
        }
    }

    return nil
}

func (c *PortChannel) RemoveFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if _, err := conn.DeleteAllByPattern(swsssdk.CONFIG_DB, []string{"PORTCHANNEL_MEMBER", c.Key, "*"}); err != nil {
        return err
    }
    if _, err := conn.SetEntry("PORTCHANNEL", c.Key, nil); err != nil {
        return err
    }
    return nil
}