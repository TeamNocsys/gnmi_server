package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
)

type Port struct {
    Key string
    Client command.Client
    Data *sonicpb.SonicPort_Port_PortList
}

func (c *Port) LoadFromDB(flags uint) error {
    if c.Client.Config() == nil {
        return swsssdk.ErrConnNotExist
    }
    if c.Data == nil {
        c.Data = &sonicpb.SonicPort_Port_PortList{}
    }
    // 读取配置信息
    if (flags & DATA_TYPE_CONFIG) != 0 {
        if data, err := c.Client.Config().GetAll(swsssdk.CONFIG_DB, []string{"PORT", c.Key}); err != nil {
            return err
        } else {
            for k, v := range data {
                switch k {
                case "index":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return err
                    } else {
                        c.Data.Index = &ywrapper.UintValue{Value: i}
                    }
                case "lanes":
                    c.Data.Lanes = &ywrapper.StringValue{Value: v}
                case "mtu":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return err
                    } else {
                        c.Data.Mtu = &ywrapper.UintValue{Value: i}
                    }
                case "alias":
                    c.Data.Alias = &ywrapper.StringValue{Value: v}
                case "admin_status":
                    switch v {
                    case "up":
                        c.Data.AdminStatus = sonicpb.SonicTypesAdminStatus_SONICTYPESADMINSTATUS_up
                    case "down":
                        c.Data.AdminStatus = sonicpb.SonicTypesAdminStatus_SONICTYPESADMINSTATUS_down
                    }
                case "speed":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return err
                    } else {
                        c.Data.Speed = &ywrapper.UintValue{Value: i}
                    }
                }
            }
        }
    }

    // 读取状态信息
    if c.Client.State() == nil {
        return swsssdk.ErrConnNotExist
    }
    if (flags & DATA_TYPE_STATE) != 0 {
        if c.Data.State == nil {
            c.Data.State = &sonicpb.SonicPort_Port_PortList_State{}
        }
        if data, err := c.Client.State().GetAll(swsssdk.APPL_DB, []string{"PORT_TABLE", c.Key}); err != nil {
            return err
        } else {
            for k, v := range data {
                switch k {
                case "oper_status":
                    switch v {
                    case "up":
                        c.Data.State.OperStatus = sonicpb.SonicTypesOperStatus_SONICTYPESOPERSTATUS_up
                    case "down":
                        c.Data.State.OperStatus = sonicpb.SonicTypesOperStatus_SONICTYPESOPERSTATUS_down
                    }
                }
            }
        }
    }

    return nil
}

func (c *Port) SaveToDB(replace bool) error {
    e := make(map[string]interface{})
    if c.Data.Alias != nil {
        e["alias"] = c.Data.Alias.Value
    }
    if c.Data.Lanes != nil {
        e["lanes"] = c.Data.Lanes.Value
    }
    if c.Data.Speed != nil {
        e["speed"] = c.Data.Speed.Value
    }
    if c.Data.Mtu != nil {
        e["mtu"] = c.Data.Mtu.Value
    }
    if c.Data.Index != nil {
        e["index"] = c.Data.Index.Value
    }
    if c.Data.AdminStatus != sonicpb.SonicTypesAdminStatus_SONICTYPESADMINSTATUS_UNSET {
        if c.Data.AdminStatus == sonicpb.SonicTypesAdminStatus_SONICTYPESADMINSTATUS_up {
            e["admin_status"] = "up"
        } else if c.Data.AdminStatus == sonicpb.SonicTypesAdminStatus_SONICTYPESADMINSTATUS_down {
            e["admin_status"] = "down"
        }
    }

    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if replace {
        if _, err := conn.SetEntry("PORT", c.Key, e); err != nil {
            return err
        }
    } else {
        if _, err := conn.ModEntry("PORT", c.Key, e); err != nil {
            return err
        }
    }

    return nil
}

func (c *Port) RemoveFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }
    if _, err := conn.SetEntry("PORT", c.Key, nil); err != nil {
        return err
    }
    return nil
}