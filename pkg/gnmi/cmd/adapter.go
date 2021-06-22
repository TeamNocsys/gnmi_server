package cmd

import (
    "context"
    "errors"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/utils"
    "strings"
)

type OperType int32

const (
    ADD     OperType = 0
    UPDATE  OperType = 1
    DEL     OperType = 2
)

type Adapter struct {
    client command.Client
}

func (adpt *Adapter) exec(cmdstr string) error {
    logrus.Trace(cmdstr + "|EXEC")
    params := strings.Split(cmdstr, " ")
    cmd := params[0]
    args := params[1:]
    if err, r := utils.Utils_execute_cmd(cmd, args...); err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return err
        }
        return errors.New(r)
    }
    return nil
}
