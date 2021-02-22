package helper

import (
    "errors"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/utils"
    "strconv"
    "strings"
)

type Ntp struct {
    Key string
    Client command.Client
    Data *sonicpb.NocsysNtp_Ntp_NtpList
}

func (c *Ntp) LoadFromDB() error {
    if err, data := utils.Utils_execute_cmd("show", "ntp"); err != nil {
        return nil
    } else {
        infos := strings.Split(data, "\n")
        for i := 2; i < len(infos) - 1; i++ {
            fields := strings.Split(infos[i], " ")
            if fields[0] == c.Key {
                if c.Data == nil {
                    c.Data = &sonicpb.NocsysNtp_Ntp_NtpList{
                        State: &sonicpb.NocsysNtp_Ntp_NtpList_State{},
                    }
                }
                if i, err := strconv.ParseUint(fields[5], 10, 64); err != nil {
                    return err
                } else {
                    c.Data.State.Poll = &ywrapper.UintValue{Value: i}
                }
                if err, now := utils.Utils_execute_cmd("date", "+\"%Y-%m-%dT%H:%M:%SZ%:z\""); err != nil {
                    return err
                } else {
                    c.Data.State.Current = &ywrapper.StringValue{Value: now}
                }
            }
        }
    }

    return nil
}

func (c *Ntp) SaveToDB(replace bool) error {
    cmdstr := "config ntp add " + c.Key
    logrus.Trace(cmdstr + "|EXEC")
    if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
        return errors.New(r)
    }
    return nil
}

func (c *Ntp) RemoveFromDB() error {
    cmdstr := "config ntp del " + c.Key
    logrus.Trace(cmdstr + "|EXEC")
    if err, r := utils.Utils_execute_cmd("bash", "-c", cmdstr); err != nil {
        return errors.New(r)
    }
    return nil
}