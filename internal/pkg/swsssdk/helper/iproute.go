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

type IpRoute struct {
    Keys []string
    Client command.Client
    Data *sonicpb.AcctonRoute_Route_RouteList
}

func (c *IpRoute) LoadFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.AcctonRoute_Route_RouteList{}
    }
    if data, err := conn.GetAll(swsssdk.APPL_DB, append([]string{"ROUTE_TABLE"}, c.Keys...)); err != nil {
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

func (c *IpRoute) SaveToDB() error {
    if c.Data.Nexthop == nil {
        return errors.New("nexthop is nil")
    }

    cmdstr := "config route add prefix vrf " + c.Keys[0] + " " + c.Keys[1] + " nexthop vrf " + c.Keys[0]
    if c.Data.Ifname != nil {
        cmdstr += " dev " + c.Data.Ifname.Value
    } else {
        cmdstr += " " + c.Data.Nexthop.Value
    }

    logrus.Trace(cmdstr + "|EXEC")
    if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
        return errors.New(r)
    }
    return nil
}

func (c *IpRoute) RemoveFromDB() error {
    cmdstr := "config route del prefix vrf " + c.Keys[0] + " " + c.Keys[1] + " nexthop vrf " + c.Keys[0] + " " + c.Keys[2]
    logrus.Trace(cmdstr + "|EXEC")
    if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
        return errors.New(r)
    }
    return nil
}
