package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strings"
)

type AclTable struct {
    Key string
    Client command.Client
    Data *sonicpb.SonicAcl_AclTable_AclTableList
}

func (c *AclTable) LoadFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.SonicAcl_AclTable_AclTableList{}
    }
    if data, err := conn.GetAll(swsssdk.APPL_DB, []string{"ACL_TABLE_TABLE", c.Key}); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "POLICY_DESC":
                c.Data.PolicyDesc = &ywrapper.StringValue{Value: v}
            case "TYPE":
                switch strings.ToUpper(v) {
                case "L2":
                    c.Data.Type = sonicpb.SonicAcl_AclTable_AclTableList_TYPE_L2
                case "L3":
                    c.Data.Type = sonicpb.SonicAcl_AclTable_AclTableList_TYPE_L3
                case "L3V6":
                    c.Data.Type = sonicpb.SonicAcl_AclTable_AclTableList_TYPE_L3V6
                case "MIRROR":
                    c.Data.Type = sonicpb.SonicAcl_AclTable_AclTableList_TYPE_MIRROR
                case "MIRRORV6":
                    c.Data.Type = sonicpb.SonicAcl_AclTable_AclTableList_TYPE_MIRRORV6
                case "MIRROR_DSCP":
                    c.Data.Type = sonicpb.SonicAcl_AclTable_AclTableList_TYPE_MIRROR_DSCP
                }
            case "stage":
                switch strings.ToUpper(v) {
                case "INGRESS":
                    c.Data.Stage = sonicpb.SonicAcl_AclTable_AclTableList_STAGE_INGRESS
                case "EGRESS":
                    c.Data.Stage = sonicpb.SonicAcl_AclTable_AclTableList_STAGE_EGRESS
                }
            case "PORTS":
                for _, servrer := range FieldToArray(v) {
                    c.Data.Ports = append(c.Data.Ports, &ywrapper.StringValue{Value: servrer})
                }
            }
        }
    }

    return nil
}

func (c *AclTable) SaveToDB(replace bool) error {
    e := make(map[string]interface{})
    if c.Data.PolicyDesc != nil {
        e["policy_desc"] = c.Data.PolicyDesc.Value
    }
    if c.Data.Type != sonicpb.SonicAcl_AclTable_AclTableList_TYPE_UNSET {
        switch c.Data.Type {
        case sonicpb.SonicAcl_AclTable_AclTableList_TYPE_L2:
            e["type"] = "L2"
        case sonicpb.SonicAcl_AclTable_AclTableList_TYPE_L3:
            e["type"] = "L3"
        case sonicpb.SonicAcl_AclTable_AclTableList_TYPE_L3V6:
            e["type"] = "L3V6"
        case sonicpb.SonicAcl_AclTable_AclTableList_TYPE_MIRROR:
            e["type"] = "MIRROR"
        case sonicpb.SonicAcl_AclTable_AclTableList_TYPE_MIRRORV6:
            e["type"] = "MIRRORV6"
        case sonicpb.SonicAcl_AclTable_AclTableList_TYPE_MIRROR_DSCP:
            e["type"] = "MIRROR_DSCP"
        }
    }
    if c.Data.Stage != sonicpb.SonicAcl_AclTable_AclTableList_STAGE_UNSET {
        switch c.Data.Stage {
        case sonicpb.SonicAcl_AclTable_AclTableList_STAGE_INGRESS:
            e["stage"] = "INGRESS"
        case sonicpb.SonicAcl_AclTable_AclTableList_STAGE_EGRESS:
            e["stage"] = "EGRESS"
        }
    }
    if c.Data.Ports != nil {
        var ports []string
        for _, port := range c.Data.Ports {
            ports = append(ports, port.Value)
        }
        e["ports"] = ports
    }

    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if replace {
        if _, err := conn.SetEntry(swsssdk.CONFIG_DB, []string{"ACL_TABLE", c.Key}, e); err != nil {
            return err
        }
    } else {
        if _, err := conn.ModEntry(swsssdk.CONFIG_DB,[]string{"ACL_TABLE", c.Key}, e); err != nil {
            return err
        }
    }

    return nil
}

func (c *AclTable) RemoveFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }
    if _, err := conn.DeleteAllByPattern(swsssdk.CONFIG_DB, []string{"ACL_RULE", c.Key, "*"}); err != nil {
        return err
    }
    if _, err := conn.SetEntry(swsssdk.CONFIG_DB,[]string{"ACL_TABLE", c.Key}, nil); err != nil {
        return err
    }
    return nil
}