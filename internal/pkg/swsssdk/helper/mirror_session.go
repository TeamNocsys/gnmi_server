package helper

import (
    "errors"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/utils"
    "strconv"
    "strings"
)

type MirrorSession struct {
    Key string
    Client command.Client
    Data *sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList
}

// 参考
// https://github.com/Azure/sonic-swss/blob/master/orchagent/mirrororch.cpp
func (c *MirrorSession) LoadFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList{}
    }
    if data, err := conn.GetAll(swsssdk.APPL_DB, []string{"MIRROR_SESSION_TABLE", c.Key}); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "status":
                switch v {
                case "active":
                    c.Data.Status = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_STATUS_active
                case "inactive":
                    c.Data.Status = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_STATUS_inactive
                }
            case "src_ip":
                c.Data.SrcIp = &ywrapper.StringValue{Value: v}
            case "dst_ip":
                c.Data.DstIp = &ywrapper.StringValue{Value: v}
            case "gre_type":
                // mellanox => 0x8949 other=> 0x88be
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.GreType = &ywrapper.UintValue{Value: i}
                }
            case "dscp":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.Dscp = &ywrapper.UintValue{Value: i}
                }
            case "ttl":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.Ttl = &ywrapper.UintValue{Value: i}
                }
            case "queue":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return err
                } else {
                    c.Data.Queue = &ywrapper.UintValue{Value: i}
                }
            case "dst_port":
                c.Data.DstPort = &ywrapper.StringValue{Value: v}
            case "src_port":
                c.Data.SrcPort = &ywrapper.StringValue{Value: v}
            case "direction":
                switch strings.ToUpper(v) {
                case "RX":
                    c.Data.Direction = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_RX
                case "TX":
                    c.Data.Direction = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_TX
                case "BOTH":
                    c.Data.Direction = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_BOTH
                }
            case "type":
                switch strings.ToUpper(v) {
                case "SPAN":
                    c.Data.Type = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_TYPE_SPAN
                case "ERSPAN":
                    c.Data.Type = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_TYPE_ERSPAN
                }
            }
        }
    }
    return nil
}

func (c *MirrorSession) SaveToDB() error {
    switch c.Data.Type {
    case sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_TYPE_SPAN:
        if c.Data.DstPort == nil ||
            c.Data.SrcPort == nil ||
            c.Data.Direction == sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_UNSET {
            return nil
        }
        cmdstr := "config mirror_session span add " + c.Key + " " + c.Data.DstPort.Value + " " + c.Data.SrcPort.Value
        switch c.Data.Direction {
        case sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_RX:
            cmdstr += " rx"
        case sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_TX:
            cmdstr += " tx"
        case sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_BOTH:
            cmdstr += " both"
        }
        logrus.Trace(cmdstr + "|EXEC")
        if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
            return errors.New(r)
        }
    case sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_TYPE_ERSPAN:
        logrus.Trace("config mirror_session erspan add |EXEC")
    }
    return nil
}

func (c *MirrorSession) RemoveFromDB() error {
    cmdstr := "config mirror_session remove " + c.Key
    logrus.Trace(cmdstr + "|EXEC")
    if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
        return errors.New(r)
    }
    return nil
}
