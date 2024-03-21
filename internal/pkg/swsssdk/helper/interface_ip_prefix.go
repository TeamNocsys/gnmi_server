package helper

import (
    "errors"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/utils"
)

type InterfaceIPPrefix struct {
    Keys []string
    Client command.Client
    Data *sonicpb.AcctonInterface_Interface_InterfaceIpprefixList
}

func (c *InterfaceIPPrefix) LoadFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.AcctonInterface_Interface_InterfaceIpprefixList{}
    }
    if data, err := conn.GetAll(swsssdk.CONFIG_DB, append([]string{"INTERFACE"}, c.Keys...)); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "scope":
                switch v {
                case "local":
                    c.Data.Scope = sonicpb.AcctonInterface_Interface_InterfaceIpprefixList_SCOPE_local
                case "global":
                    c.Data.Scope = sonicpb.AcctonInterface_Interface_InterfaceIpprefixList_SCOPE_global
                }
            case "family":
                switch v {
                case "IPv4":
                    c.Data.Family = sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_IPv4
                case "IPv6":
                    c.Data.Family = sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_IPv6
                }
            }
        }
    }

    return nil
}

func (c *InterfaceIPPrefix) SaveToDB(replace bool) error {
    if c.Data.Family != sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_UNSET {
        switch c.Data.Family {
        case sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_IPv4:
            cmdstr := "config interface ip add " + c.Keys[0] + " " + c.Keys[1]
            logrus.Trace(cmdstr + "|EXEC")
            if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
                return errors.New(r)
            }
        case sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_IPv6:
            cmdstr := "config interface ipv6 add " + c.Keys[0] + " " + c.Keys[1]
            logrus.Trace(cmdstr + "|EXEC")
            if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
                return errors.New(r)
            }
        }
    }

    return nil
}

func (c *InterfaceIPPrefix) RemoveFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }
    if _, err := conn.DeleteAllByPattern(swsssdk.CONFIG_DB, append([]string{"INTERFACE"}, c.Keys...)); err != nil {
        return err
    }
    return nil
}
