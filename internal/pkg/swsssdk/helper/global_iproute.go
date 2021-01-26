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

type GlobalIpRoute struct {
    Key string
    Client command.Client
    Data *sonicpb.SonicRoute_Route_GlobalRouteList
}

func (c *GlobalIpRoute) LoadFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.SonicRoute_Route_GlobalRouteList{}
    }
    if data, err := conn.GetAll(swsssdk.APPL_DB, []string{"ROUTE_TABLE", c.Key}); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "nexthop":
                c.Data.Nexthop = &ywrapper.StringValue{Value: v}
            case "ifname":
                c.Data.Ifname = &ywrapper.StringValue{Value: v}
            }
        }
    }

    return nil
}

func (c *GlobalIpRoute) SaveToDB() error {
    if c.Data.Nexthop == nil {
        return errors.New("nexthop is nil")
    }

    cmdstr := "ip route replace prefix " + c.Key + " nexthop " + c.Data.Nexthop.Value
    if c.Data.Ifname != nil {
        cmdstr += " dev " + c.Data.Ifname.Value
    }

    logrus.Trace(cmdstr + "|EXEC")
    if err, r := utils.Utils_execute_cmd("vtysh", "-c", cmdstr); err != nil {
        return errors.New(r)
    }
    return nil
}

func (c *GlobalIpRoute) RemoveFromDB() error {
    cmdstr := "ip route del prefix " + c.Key
    logrus.Trace(cmdstr + "|EXEC")
    if err, r := utils.Utils_execute_cmd("vtysh", "-c", cmdstr); err != nil {
        return errors.New(r)
    }
    return nil
}