package helper

import (
    "errors"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/utils"
)

type VlanInterfaceIPPrefix struct {
    Keys []string
    Client command.Client
    Data *sonicpb.SonicVlan_VlanInterface_VlanInterfaceIpprefixList
}

func (c *VlanInterfaceIPPrefix) LoadFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.SonicVlan_VlanInterface_VlanInterfaceIpprefixList{}
    }
    if data, err := conn.GetAll(swsssdk.CONFIG_DB, append([]string{"VLAN_INTERFACE"}, c.Keys...)); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "scope":
                switch v {
                case "local":
                    c.Data.Scope = sonicpb.SonicVlan_VlanInterface_VlanInterfaceIpprefixList_SCOPE_local
                case "global":
                    c.Data.Scope = sonicpb.SonicVlan_VlanInterface_VlanInterfaceIpprefixList_SCOPE_global
                }
            case "family":
                switch v {
                case "IPv4":
                    c.Data.Family = sonicpb.SonicTypesIpFamily_SONICTYPESIPFAMILY_IPv4
                case "IPv6":
                    c.Data.Family = sonicpb.SonicTypesIpFamily_SONICTYPESIPFAMILY_IPv6
                }
            }
        }
    }

    return nil
}

func (c *VlanInterfaceIPPrefix) SaveToDB() error {
    if c.Data.Family != sonicpb.SonicTypesIpFamily_SONICTYPESIPFAMILY_UNSET {
        switch c.Data.Family {
        case sonicpb.SonicTypesIpFamily_SONICTYPESIPFAMILY_IPv4:
            cmdstr := "config interface ip add " + c.Keys[0] + " " + c.Keys[1]
            logrus.Trace(cmdstr + "|EXEC")
            if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
                return errors.New(r)
            }
        case sonicpb.SonicTypesIpFamily_SONICTYPESIPFAMILY_IPv6:
            cmdstr := "config interface ipv6 add " + c.Keys[0] + " " + c.Keys[1]
            logrus.Trace(cmdstr + "|EXEC")
            if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
                return errors.New(r)
            }
        }
    }

    return nil
}

func (c *VlanInterfaceIPPrefix) RemoveFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }
    if _, err := conn.DeleteAllByPattern(swsssdk.CONFIG_DB, append([]string{"VLAN_INTERFACE"}, c.Keys...)); err != nil {
        return err
    }
    return nil
}
