package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/utils"
    "strconv"
)

type Vlan struct {
    Key string
    Client command.Client
    Data *sonicpb.NocsysVlan_Vlan_VlanList
}

func (c *Vlan) LoadFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.NocsysVlan_Vlan_VlanList{}
    }
    if data, err := conn.GetEntry("VLAN", c.Key); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "description":
                c.Data.Description = &ywrapper.StringValue{Value: v.(string)}
            case "dhcp_servers":
                for _, server := range v.([]string) {
                    c.Data.DhcpServers = append(c.Data.DhcpServers, &ywrapper.StringValue{Value: server})
                }
            case "mtu":
                if i, err := strconv.ParseUint(v.(string), 10, 64); err != nil {
                    return err
                } else {
                    c.Data.Mtu = &ywrapper.UintValue{Value: i}
                }
            case "admin_status":
                switch v {
                case "up":
                    c.Data.AdminStatus = sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_up
                case "down":
                    c.Data.AdminStatus = sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_down
                }
            }
        }
    }
    return nil
}

func (c *Vlan) SaveToDB(replace bool) error {
    e := make(map[string]interface{})
    if c.Data.Description != nil {
        e["description"] = c.Data.Description.Value
    }
    if c.Data.DhcpServers != nil {
        var servers []interface{}
        for _, server := range c.Data.DhcpServers {
            servers = append(servers, server.Value)
        }
        e["dhcp_servers"] = servers
    }
    if c.Data.Mtu != nil {
        e["mtu"] = c.Data.Mtu.Value
    }
    if c.Data.AdminStatus != sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_UNSET {
        if c.Data.AdminStatus == sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_up {
            e["admin_status"] = "up"
        } else if c.Data.AdminStatus == sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_down {
            e["admin_status"] = "down"
        }
    }

    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if replace {
        if _, err := conn.SetEntry("VLAN", c.Key, e); err != nil {
            return err
        }
    } else {
        if _, err := conn.ModEntry("VLAN", c.Key, e); err != nil {
            return err
        }
    }

    if c.Data.DhcpServers != nil {
        utils.Utils_execute_cmd("systemctl", "stop", "dhcp_relay")
        utils.Utils_execute_cmd("systemctl", "reset-failed", "dhcp_relay")
        utils.Utils_execute_cmd("systemctl", "start", "dhcp_relay")
    }

    return nil
}

func (c *Vlan) RemoveFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if _, err := conn.DeleteAllByPattern(swsssdk.CONFIG_DB, []string{"VLAN_MEMBER", c.Key, "*"}); err != nil {
        return err
    }
    if _, err := conn.SetEntry("VLAN", c.Key, nil); err != nil {
        return err
    }
    return nil
}