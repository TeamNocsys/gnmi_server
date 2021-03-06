package helper

import (
    "errors"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/utils"
)

type Interface struct {
    Key string
    Client command.Client
    Data *sonicpb.NocsysInterface_Interface_InterfaceList
}

func (c *Interface) LoadFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.NocsysInterface_Interface_InterfaceList{}
    }
    if data, err := conn.GetAll(swsssdk.CONFIG_DB, []string{"INTERFACE", c.Key}); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "vrf_name":
                c.Data.VrfName = &ywrapper.StringValue{Value: v}
            }
        }
    }

    return nil
}

func (c *Interface) SaveToDB(replace bool) error {
    e := make(map[string]interface{})
    if c.Data.VrfName != nil {
        cmdstr := "config interface vrf bind " + c.Key + " " + c.Data.VrfName.Value
        logrus.Trace(cmdstr + "|EXEC")
        if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
            return errors.New(r)
        }
        return nil
    }

    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if replace {
        if _, err := conn.SetEntry("INTERFACE", c.Key, e); err != nil {
            return err
        }
    } else {
        if _, err := conn.ModEntry("INTERFACE", c.Key, e); err != nil {
            return err
        }
    }

    return nil
}

func (c *Interface) RemoveFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.CONFIG_DB, []string{"INTERFACE", c.Key}); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "vrf_name":
                cmdstr := "config interface vrf unbind " + c.Key + " " + v
                logrus.Trace(cmdstr + "|EXEC")
                if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
                    return errors.New(r)
                }
                return nil
            }
        }
    }

    if _, err := conn.DeleteAllByPattern("INTERFACE", []string{c.Key, "*"}); err != nil {
        return err
    }
    if _, err := conn.Delete("INTERFACE", c.Key); err != nil {
        return err
    }
    return nil
}