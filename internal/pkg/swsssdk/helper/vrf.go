package helper

import (
    "errors"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/utils"
    "strings"
)

type Vrf struct {
    Key string
    Client command.Client
    Data *sonicpb.NocsysVrf_Vrf_VrfList
}

func (c *Vrf) LoadFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.NocsysVrf_Vrf_VrfList{}
    }
    if data, err := conn.GetAll(swsssdk.CONFIG_DB, []string{"VRF", c.Key}); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "fallback":
                switch strings.ToLower(v) {
                case "true":
                    c.Data.Fallback = &ywrapper.BoolValue{Value: true}
                case "false":
                    c.Data.Fallback = &ywrapper.BoolValue{Value: false}
                }
            }
        }
    }
    return nil
}

func (c *Vrf) SaveToDB() error {
    cmdstr := "config vrf add " + c.Key
    logrus.Trace(cmdstr + "|EXEC")
    if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
        return errors.New(r)
    }
    return nil
}

func (c *Vrf) RemoveFromDB() error {
    cmdstr := "config vrf del " + c.Key
    logrus.Trace(cmdstr + "|EXEC")
    if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
        return errors.New(r)
    }
    return nil
}